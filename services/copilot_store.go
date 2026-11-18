package services

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"sale-app/models"

	"gorm.io/gorm/clause"
)

const (
	copilotMaxStoredMessages = 400 // full history kept per thread
	copilotMaxMemImagesB64   = 12  // newest images keeping base64 in memory (UI preview)
	copilotMemSessions       = 30  // sessions kept in memory / listed on the Copilot screen
	copilotAIContextMessages = 60  // history tail sent to OpenAI for "suggest"
)

// CopilotSessionRecord persists a Copilot conversation so it survives deploys / restarts.
// Data holds the JSON of CopilotSession (without image_base64 to keep rows small).
type CopilotSessionRecord struct {
	ThreadID     string `gorm:"primaryKey;size:191"`
	CustomerName string `gorm:"size:255"`
	Data         string `gorm:"type:longtext"`
	CreatedAt    time.Time
	UpdatedAt    time.Time `gorm:"index"`
}

func (CopilotSessionRecord) TableName() string { return "copilot_sessions" }

func (s *CopilotService) initStore() {
	if s.DB == nil {
		return
	}
	if err := s.DB.AutoMigrate(&CopilotSessionRecord{}); err != nil {
		log.Printf("[copilot] AutoMigrate copilot_sessions error: %v", err)
		return
	}
	s.purgeLegacyExtensionSessions()

	var recs []CopilotSessionRecord
	if err := s.DB.Order("updated_at DESC").Limit(copilotMemSessions).Find(&recs).Error; err != nil {
		log.Printf("[copilot] load sessions error: %v", err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, rec := range recs {
		sess := decodeSessionRecord(rec)
		if sess == nil {
			continue
		}
		s.refreshMessagesFromRecords(sess)
		s.sessions[sess.ThreadID] = sess
		s.sessionList = append(s.sessionList, sess.ThreadID)
	}
	log.Printf("[copilot] loaded %d sessions from DB", len(s.sessionList))
}

func decodeSessionRecord(rec CopilotSessionRecord) *CopilotSession {
	var sess CopilotSession
	if err := json.Unmarshal([]byte(rec.Data), &sess); err != nil {
		return nil
	}
	if sess.ThreadID == "" {
		sess.ThreadID = rec.ThreadID
	}
	return &sess
}

// loadSession returns a session from memory, falling back to the DB (older threads not in memory).
func (s *CopilotService) loadSession(threadID string) (*CopilotSession, bool) {
	if sess, ok := s.GetSession(threadID); ok {
		return sess, true
	}
	if s.DB == nil {
		return nil, false
	}
	var rec CopilotSessionRecord
	if err := s.DB.Where("thread_id = ?", threadID).Limit(1).Find(&rec).Error; err != nil || rec.ThreadID == "" {
		return nil, false
	}
	sess := decodeSessionRecord(rec)
	return sess, sess != nil
}

func (s *CopilotService) persistSession(sess *CopilotSession) {
	if s.DB == nil || sess == nil {
		return
	}
	cp := *sess
	cp.Messages = make([]CopilotMessage, len(sess.Messages))
	for i, m := range sess.Messages {
		m.ImageBase64 = ""
		cp.Messages[i] = m
	}
	b, err := json.Marshal(&cp)
	if err != nil {
		log.Printf("[copilot] marshal session %s error: %v", sess.ThreadID, err)
		return
	}
	now := time.Now()
	rec := CopilotSessionRecord{
		ThreadID:     sess.ThreadID,
		CustomerName: sess.CustomerName,
		Data:         string(b),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	// Upsert that never rewrites created_at: DB.Save() updated every column including a zero
	// created_at, which TiDB rejects, so only the first snapshot of a thread was ever stored.
	err = s.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "thread_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"customer_name", "data", "updated_at"}),
	}).Create(&rec).Error
	if err != nil {
		log.Printf("[copilot] save session %s error: %v", sess.ThreadID, err)
	}
}

// messageKey identifies a message across snapshots. Facebook CDN image URLs keep the same
// path but rotate their signed query string, so images are compared by path only.
func messageKey(m CopilotMessage) string {
	img := ""
	if m.Type == "image" && m.ImageURL != "" {
		if u, err := url.Parse(m.ImageURL); err == nil {
			// Filename is unique per photo; CDN host shards differ between page loads.
			img = u.Path[strings.LastIndex(u.Path, "/")+1:]
		} else {
			img = m.ImageURL
		}
	}
	return m.Sender + "|" + m.Type + "|" + strings.TrimSpace(m.Text) + "|" + img
}

// mergeCopilotMessages appends a new extension snapshot onto the stored history.
//   - Tail of history overlaps head of snapshot  -> history[:i] + snapshot
//   - Snapshot is fully inside history (older view) -> history unchanged
//   - No overlap (gap)                             -> history + snapshot
func mergeCopilotMessages(history, snap []CopilotMessage) []CopilotMessage {
	if len(history) == 0 {
		return capMessages(append([]CopilotMessage(nil), snap...))
	}
	if len(snap) == 0 {
		return history
	}
	hk := make([]string, len(history))
	for i, m := range history {
		hk[i] = messageKey(m)
	}
	sk := make([]string, len(snap))
	for i, m := range snap {
		sk[i] = messageKey(m)
	}

	start := len(history) - len(snap)
	if start < 0 {
		start = 0
	}
	for i := start; i < len(history); i++ {
		match := true
		for j := 0; i+j < len(history); j++ {
			if hk[i+j] != sk[j] {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		out := make([]CopilotMessage, 0, i+len(snap))
		out = append(out, history[:i]...)
		for j, m := range snap {
			// Keep OCR text already computed for this message earlier.
			if m.OCRText == "" && i+j < len(history) {
				m.OCRText = history[i+j].OCRText
			}
			out = append(out, m)
		}
		return capMessages(out)
	}

	if containsSequence(hk, sk) {
		return history
	}
	return capMessages(append(append([]CopilotMessage(nil), history...), snap...))
}

func containsSequence(hay, needle []string) bool {
	for i := 0; i+len(needle) <= len(hay); i++ {
		ok := true
		for j := range needle {
			if hay[i+j] != needle[j] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func capMessages(msgs []CopilotMessage) []CopilotMessage {
	if len(msgs) > copilotMaxStoredMessages {
		msgs = msgs[len(msgs)-copilotMaxStoredMessages:]
	}
	// Only the newest images keep base64 in memory; older ones render from image_url.
	kept := 0
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].ImageBase64 == "" {
			continue
		}
		if kept < copilotMaxMemImagesB64 {
			kept++
			continue
		}
		msgs[i].ImageBase64 = ""
	}
	return msgs
}

func tailMessages(msgs []CopilotMessage, n int) []CopilotMessage {
	if len(msgs) > n {
		return msgs[len(msgs)-n:]
	}
	return msgs
}

// PurgeThreads removes sessions stored under invalid thread ids (e.g. the Page id that older
// extension builds used when selected_item_id was missing, which split one customer in two).
func (s *CopilotService) PurgeThreads(ids ...string) {
	var valid []string
	for _, id := range ids {
		if id = strings.TrimSpace(id); id != "" {
			valid = append(valid, id)
		}
	}
	if len(valid) == 0 {
		return
	}
	s.mu.Lock()
	for _, id := range valid {
		delete(s.sessions, id)
	}
	kept := s.sessionList[:0]
	for _, id := range s.sessionList {
		if _, ok := s.sessions[id]; ok {
			kept = append(kept, id)
		}
	}
	s.sessionList = kept
	s.mu.Unlock()
	if s.DB != nil {
		if err := s.DB.Where("thread_id IN ?", valid).Delete(&CopilotSessionRecord{}).Error; err != nil {
			log.Printf("[copilot] purge threads error: %v", err)
		}
	}
}

// purgeLegacyExtensionSessions removes sessions created by the Chrome extension, which keyed
// threads by the Business Suite id instead of the Messenger PSID, so the same customer showed
// up twice next to the webhook/Graph session. The webhook is now the source of truth: any stored
// session whose thread_id is not a known PSID is dropped.
func (s *CopilotService) purgeLegacyExtensionSessions() {
	if s.DB == nil || !s.DB.Migrator().HasTable(&models.CopilotConversation{}) {
		return
	}
	var convCount int64
	if err := s.DB.Model(&models.CopilotConversation{}).Count(&convCount).Error; err != nil || convCount == 0 {
		return // never wipe everything when the PSID table is empty or unreadable
	}
	res := s.DB.Where("thread_id NOT IN (?)", s.DB.Model(&models.CopilotConversation{}).Select("ps_id")).
		Delete(&CopilotSessionRecord{})
	if res.Error != nil {
		log.Printf("[copilot] purge legacy sessions error: %v", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("[copilot] purged %d legacy extension sessions", res.RowsAffected)
	}
}

// refreshMessagesFromRecords reloads a webhook/Graph session's messages and LastMessageAt from
// copilot_message_records (the source of truth), so a stale stored snapshot never hides newer
// messages after a restart. AI fields are kept as stored (no OpenAI calls at startup).
func (s *CopilotService) refreshMessagesFromRecords(sess *CopilotSession) {
	if s.DB == nil || sess == nil || sess.ThreadID == "" {
		return
	}
	var conv models.CopilotConversation
	if err := s.DB.Where(&models.CopilotConversation{PSID: sess.ThreadID}).Limit(1).Find(&conv).Error; err != nil || conv.ID == 0 {
		return
	}
	// Order was deleted: drop the stale marker so earlier messages show again.
	if conv.LastOrderID != nil {
		var n int64
		s.DB.Model(&models.Order{}).Where("id = ?", *conv.LastOrderID).Count(&n)
		if n == 0 {
			s.DB.Model(&conv).Updates(map[string]interface{}{"last_order_id": nil, "last_order_at": nil})
			conv.LastOrderID = nil
			conv.LastOrderAt = nil
		}
	}
	q := s.DB.Where("conversation_id = ?", conv.ID)
	if conv.LastOrderAt != nil {
		q = q.Where("msg_timestamp > ?", *conv.LastOrderAt)
	}
	var recs []models.CopilotMessageRecord
	if err := q.Order("msg_timestamp DESC, id DESC").Limit(30).Find(&recs).Error; err != nil {
		return
	}
	msgs := make([]CopilotMessage, 0, len(recs))
	for i := len(recs) - 1; i >= 0; i-- {
		r := recs[i]
		if r.Type == "text" && IsFBSystemNotice(r.Text) {
			continue
		}
		msgs = append(msgs, CopilotMessage{Sender: r.Sender, Type: r.Type, Text: r.Text, ImageURL: r.ImageURL, OCRText: r.OCRText})
	}
	sess.Messages = msgs
	if len(recs) > 0 {
		sess.LastMessageAt = recs[0].MsgTimestamp
	}
	if conv.CustomerName != "" {
		sess.CustomerName = conv.CustomerName
	}
	if conv.AvatarURL != "" {
		sess.AvatarURL = conv.AvatarURL
	}
	if conv.LastOrderID == nil {
		sess.CreatedOrder = nil
	} else if sess.CreatedOrder == nil || sess.CreatedOrder.ID != *conv.LastOrderID {
		var ord models.Order
		if err := s.DB.Preload("Products").Limit(1).Find(&ord, *conv.LastOrderID).Error; err == nil && ord.ID != 0 {
			s.FormatOrderMedia(&ord)
			sess.CreatedOrder = &ord
		} else {
			sess.CreatedOrder = nil
		}
	}
}

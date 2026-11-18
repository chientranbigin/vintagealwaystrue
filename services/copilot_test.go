package services

import (
	"context"
	"testing"
)

func TestCopilot_ExtractProductCodesFromText(t *testing.T) {
	s := NewCopilotService(nil, nil, "")
	codes := s.ExtractProductCodesFromText("Dạ em hỏi áo BZ 105 với quần ts-42 và GL09 còn không ạ? BZ105 nữa")
	if len(codes) != 3 {
		t.Fatalf("expected 3 unique codes, got %d: %v", len(codes), codes)
	}
	expected := map[string]bool{"BZ105": true, "TS42": true, "GL09": true}
	for _, c := range codes {
		if !expected[c] {
			t.Fatalf("unexpected code %s", c)
		}
	}
}

func TestCopilot_ExtractHeightWeight(t *testing.T) {
	h, w := extractHeightWeight("Mình cao 1m72 nặng 68kg mặc vừa không shop?")
	if h != 172 || w != 68 {
		t.Fatalf("expected 172cm 68kg, got %d %d", h, w)
	}
}

func TestCopilot_AnalyzeAndBroadcast(t *testing.T) {
	s := NewCopilotService(nil, nil, "")
	s.OpenAIKey = "" // force deterministic fallback in unit test

	ch := s.Subscribe()
	defer s.Unsubscribe(ch)

	sess := s.AnalyzeIngest(context.Background(), CopilotIngestRequest{
		ThreadID:     "t-1",
		CustomerName: "Minh Anh",
		Action:       "suggest",
		Messages: []CopilotMessage{
			{Sender: "customer", Type: "text", Text: "Shop ơi áo này giảm giá chút được không?"},
		},
	})

	if !sess.NeedsHuman {
		t.Fatalf("expected NeedsHuman=true for discount request")
	}

	s.SaveAndBroadcast(sess)
	select {
	case got := <-ch:
		if got.ThreadID != "t-1" {
			t.Fatalf("expected thread t-1, got %s", got.ThreadID)
		}
	default:
		t.Fatalf("expected broadcast message on SSE channel")
	}

	list := s.ListSessions()
	if len(list) != 1 || list[0].CustomerName != "Minh Anh" {
		t.Fatalf("unexpected sessions list: %+v", list)
	}
}

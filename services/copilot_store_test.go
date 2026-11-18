package services

import "testing"

func txt(sender, t string) CopilotMessage { return CopilotMessage{Sender: sender, Type: "text", Text: t} }
func img(sender, u string) CopilotMessage { return CopilotMessage{Sender: sender, Type: "image", ImageURL: u} }

func texts(ms []CopilotMessage) []string {
	out := make([]string, len(ms))
	for i, m := range ms {
		out[i] = messageKey(m)
	}
	return out
}

func eq(t *testing.T, got, want []CopilotMessage) {
	t.Helper()
	g, w := texts(got), texts(want)
	if len(g) != len(w) {
		t.Fatalf("len %d != %d\ngot  %v\nwant %v", len(g), len(w), g, w)
	}
	for i := range g {
		if g[i] != w[i] {
			t.Fatalf("idx %d: %q != %q", i, g[i], w[i])
		}
	}
}

func TestMergeAppendsNewMessages(t *testing.T) {
	h := []CopilotMessage{txt("customer", "a"), txt("shop", "b"), txt("customer", "c")}
	snap := []CopilotMessage{txt("shop", "b"), txt("customer", "c"), txt("shop", "d")}
	eq(t, mergeCopilotMessages(h, snap), []CopilotMessage{h[0], h[1], h[2], snap[2]})
}

func TestMergeSameSnapshotIsIdempotent(t *testing.T) {
	h := []CopilotMessage{txt("customer", "a"), txt("shop", "b")}
	eq(t, mergeCopilotMessages(h, h), h)
}

func TestMergeOlderViewKeepsHistory(t *testing.T) {
	h := []CopilotMessage{txt("customer", "a"), txt("shop", "b"), txt("customer", "c"), txt("shop", "d")}
	snap := []CopilotMessage{txt("shop", "b"), txt("customer", "c")}
	eq(t, mergeCopilotMessages(h, snap), h)
}

func TestMergeRepeatedTextUsesLongestOverlap(t *testing.T) {
	h := []CopilotMessage{txt("shop", "Dạ okay"), txt("customer", "x"), txt("shop", "Dạ okay")}
	snap := []CopilotMessage{txt("customer", "x"), txt("shop", "Dạ okay"), txt("customer", "y")}
	eq(t, mergeCopilotMessages(h, snap), []CopilotMessage{h[0], h[1], h[2], snap[2]})
}

func TestMergeImageIgnoresSignedQueryAndHost(t *testing.T) {
	h := []CopilotMessage{img("shop", "https://scontent.fsgn2-3.fna.fbcdn.net/v/t1/123_n.jpg?oh=AAA&oe=1")}
	snap := []CopilotMessage{img("shop", "https://scontent.fsgn2-8.fna.fbcdn.net/v/t1/123_n.jpg?oh=BBB&oe=2"), txt("customer", "ok")}
	got := mergeCopilotMessages(h, snap)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestMergeGapAppends(t *testing.T) {
	h := []CopilotMessage{txt("customer", "a")}
	snap := []CopilotMessage{txt("customer", "z")}
	eq(t, mergeCopilotMessages(h, snap), []CopilotMessage{h[0], snap[0]})
}

func TestFBSystemNoticeAndLikeSticker(t *testing.T) {
	if !IsFBSystemNotice("Tran Anh Tu đã trả lời một quảng cáo.") {
		t.Fatal("ad banner should be a system notice")
	}
	if !IsFBSystemNotice("Lê Minh Quangg đã trả lời về một bài viết. Xem bài viết") {
		t.Fatal("post-reply banner should be a system notice")
	}
	if !IsFBSystemNotice("Lê Minh Quangg đã trả lời về một bài viết. Xem bài viết(https://www.facebook.com/x/posts/pfbid0q)") {
		t.Fatal("post-reply banner with link should be a system notice")
	}
	if IsFBSystemNotice("Shop đã trả lời anh rồi ạ, anh xem giúp em") {
		t.Fatal("normal text must not be filtered")
	}
	if !IsLikeStickerAttachment("sticker_369239263222822", "") || !IsLikeStickerAttachment("", "sticker-369239383222810") {
		t.Fatal("like sticker not detected")
	}
	if IsLikeStickerAttachment("1703845684703879", "image-1703845684703879") {
		t.Fatal("photo detected as like")
	}
}

package services

import (
	"strings"
	"testing"
)

func txt(sender, t string) CopilotMessage {
	return CopilotMessage{Sender: sender, Type: "text", Text: t}
}
func img(sender, u string) CopilotMessage {
	return CopilotMessage{Sender: sender, Type: "image", ImageURL: u}
}

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
	if !IsFBSystemNotice("Đã thêm nhãn tự động: Đã đánh dấu trạng thái đơn đặt hàng là Đã đặt hàng.") {
		t.Fatal("Business Suite auto-label line should be a system notice")
	}
	if IsFBSystemNotice("Dạ em đã thêm nhãn cho đơn của anh rồi nha") {
		t.Fatal("normal shop text mentioning a label must not be filtered")
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

func TestOpenAICostUSD(t *testing.T) {
	// 1M input + 1M output on gpt-4o-mini = 0.15 + 0.60
	if c := openAICostUSD("gpt-4o-mini", 1_000_000, 1_000_000); c < 0.749 || c > 0.751 {
		t.Fatalf("gpt-4o-mini cost = %v", c)
	}
	if c := openAICostUSD("gpt-4o-mini-2024-07-18", 1_000_000, 0); c < 0.149 || c > 0.151 {
		t.Fatalf("dated model should match prefix pricing, got %v", c)
	}
}

func TestSoldCodesPresentedAsAvailable(t *testing.T) {
	sold := []string{"TS4799", "TS4779"}
	bad := soldCodesPresentedAsAvailable("Dạ, hiện tại em còn chiếc quần mã TS4764 giá 350k. Chiếc mã TS4799 cũng có giá 350k, eo 42.", sold)
	if len(bad) != 1 || bad[0] != "TS4799" {
		t.Fatalf("expected TS4799 flagged, got %v", bad)
	}
	if bad := soldCodesPresentedAsAvailable("Dạ TS4764 350k còn nha 😊 còn TS4799, TS4779 em bán mất r ạ", sold); len(bad) != 0 {
		t.Fatalf("sold mention with 'bán' must pass, got %v", bad)
	}
}

func TestPendingCustomerCodesOnlyAfterLastShopReply(t *testing.T) {
	s := &CopilotService{}
	msgs := []CopilotMessage{
		{Sender: "customer", Type: "image", OCRText: "TS4763 E80"},
		{Sender: "shop", Type: "text", Text: "Dạ TS4763 em bán r ạ"},
		{Sender: "customer", Type: "image", OCRText: "TS4796 E82 - DA33"},
		{Sender: "customer", Type: "text", Text: "bz1536 còn k shop"},
	}
	got := s.pendingCustomerCodes(msgs)
	if len(got) != 2 || got[0] != "TS4796" || got[1] != "BZ1536" {
		t.Fatalf("pending = %v", got)
	}
}

func TestPriceK(t *testing.T) {
	for in, want := range map[int64]string{350000: "350k", 1250000: "1tr250", 2000000: "2tr"} {
		if got := priceK(in); got != want {
			t.Fatalf("priceK(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildProductContextGroups(t *testing.T) {
	ctx, sold := buildProductContext([]CopilotProductMatch{
		{Code: "TS4764", Type: "TROUSERS", Price: 350000, Status: "AVAILABLE"},
		{Code: "TS4799", Type: "TROUSERS", Price: 350000, Status: "SOLD"},
	}, []string{"TS4764", "TS4799"})
	if len(sold) != 1 || sold[0] != "TS4799" {
		t.Fatalf("sold = %v", sold)
	}
	if !strings.Contains(ctx, "CÒN HÀNG") || !strings.Contains(ctx, "- TS4764 | TROUSERS | 350k") || !strings.Contains(ctx, "ĐÃ BÁN (không còn, không được nói là còn hàng): TS4799") {
		t.Fatalf("context:\n%s", ctx)
	}
}

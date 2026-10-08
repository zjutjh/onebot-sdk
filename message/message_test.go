package message

import (
	"reflect"
	"testing"

	json "github.com/zjutjh/onebot-sdk/internal/jsonx"
)

// TestConstructorsRoundTrip 确认每个构造器产出的段能编码、再解析回同类型,且再编码字节一致。
// node 段仍按 UnknownSegmentData 保留,只校验 type 不校验字节。
func TestConstructorsRoundTrip(t *testing.T) {
	cases := []struct {
		name     string
		seg      Segment
		wantType string
		wantData any
	}{
		{"text", Text("hi"), "text", TextData{}},
		{"at", At(10001), "at", AtData{}},
		{"at_all", AtAll(), "at", AtData{}},
		{"image", Image("a.png"), "image", ImageData{}},
		{"reply", Reply(42), "reply", ReplyData{}},
		{"record", Record("a.silk"), "record", RecordData{}},
		{"video", Video("a.mp4"), "video", VideoData{}},
		{"face", Face("14"), "face", FaceData{}},
		{"json", JSON(`{"a":1}`), "json", JSONData{}},
		{"mface", MFace(3, "e1"), "mface", MFaceData{}},
		{"card", Card("json", `{"k":"v"}`), "card", CardData{}},
		{"forward", Forward("res"), "forward", ForwardData{}},
		{"file", File("f1"), "file", FileData{}},
		{"node_custom", NodeCustom("n", 10001, Text("x")), "node", nil},
		{"node_id", NodeID("123"), "node", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			first, err := json.Marshal(tc.seg)
			if err != nil {
				t.Fatalf("编码失败: %v", err)
			}
			var got Segment
			if err := json.Unmarshal(first, &got); err != nil {
				t.Fatalf("解析失败: %v", err)
			}
			if got.Type != tc.wantType {
				t.Fatalf("type = %q, want %q", got.Type, tc.wantType)
			}
			if tc.wantData == nil {
				return
			}
			if reflect.TypeOf(got.Data) != reflect.TypeOf(tc.wantData) {
				t.Fatalf("data 类型 = %T, want %T", got.Data, tc.wantData)
			}
			second, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("重新编码失败: %v", err)
			}
			if string(first) != string(second) {
				t.Fatalf("往返字节不一致: %s -> %s", first, second)
			}
		})
	}
}

// TestCardUnmarshal 确认 card 段被识别为 CardData 且内层 data 原样保留。
func TestCardUnmarshal(t *testing.T) {
	var seg Segment
	if err := json.Unmarshal([]byte(`{"type":"card","data":{"type":"json","data":{"x":1}}}`), &seg); err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	card, ok := seg.Data.(CardData)
	if !ok {
		t.Fatalf("data 类型 = %T, want CardData", seg.Data)
	}
	if card.Type != "json" {
		t.Fatalf("card.type = %q, want json", card.Type)
	}
	if string(card.Data) != `{"x":1}` {
		t.Fatalf("card.data = %s, want {\"x\":1}", card.Data)
	}
}

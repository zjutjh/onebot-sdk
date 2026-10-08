package onebot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/message"
)

// loadFixture 读取 testdata 下的双后端事件样本。
// 样本字段形态按两后端源码核对后构造,作为字段级 golden 回归基线。
func loadFixture(t *testing.T, parts ...string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(append([]string{"testdata"}, parts...)...))
	if err != nil {
		t.Fatalf("读取 fixture 失败: %v", err)
	}
	return data
}

// 断言事件是预期类型并返回,类型不符即失败。
func wantEvent[E event.Event](t *testing.T, data []byte) E {
	t.Helper()
	ev, ok := event.Parse(data).(E)
	if !ok {
		t.Fatalf("解析为 %T,期望 %T", event.Parse(data), ev)
	}
	return ev
}

func TestNapcatGroupMessageGolden(t *testing.T) {
	msg := wantEvent[*event.GroupMessage](t, loadFixture(t, "napcat", "group_message.json"))
	if msg.MessageID.Int64() != 1234567890 || msg.GroupID.Int64() != 30001 || msg.UserID.Int64() != 20001 {
		t.Fatalf("群消息 ID 字段错误: %+v", msg)
	}
	if msg.MessageSeq != 41 || msg.Sender.Card != "爱丽丝" || msg.Sender.Role != "admin" {
		t.Fatalf("群消息扩展字段错误: %+v", msg.Sender)
	}
	if len(msg.Message) != 7 {
		t.Fatalf("段数量 = %d，期望 7", len(msg.Message))
	}
	if d, ok := msg.Message[0].Data.(message.ReplyData); !ok || d.ID.Int64() != 1234567889 {
		t.Fatalf("reply 段解析错误: %+v", msg.Message[0].Data)
	}
	if d, ok := msg.Message[1].Data.(message.AtData); !ok || d.QQ != "all" {
		t.Fatalf("at(all) 段解析错误: %+v", msg.Message[1].Data)
	}
	if d, ok := msg.Message[3].Data.(message.ImageData); !ok || d.URL != "https://example.com/abc.jpg" || d.Summary != "一张图片" || d.SubType != "0" {
		t.Fatalf("image 段解析错误: %+v", msg.Message[3].Data)
	}
	if d, ok := msg.Message[4].Data.(message.MFaceData); !ok || d.EmojiPackageID != 1 || d.EmojiID != "2" {
		t.Fatalf("mface 段解析错误: %+v", msg.Message[4].Data)
	}
	if d, ok := msg.Message[5].Data.(message.FaceData); !ok || d.ID != "178" {
		t.Fatalf("face 段解析错误: %+v", msg.Message[5].Data)
	}
	if d, ok := msg.Message[6].Data.(message.UnknownSegmentData); !ok || d.Fields["foo"] != "bar" {
		t.Fatalf("未知段应保留原始字段: %+v", msg.Message[6].Data)
	}
}

func TestNapcatPrivateMessageGolden(t *testing.T) {
	msg := wantEvent[*event.PrivateMessage](t, loadFixture(t, "napcat", "private_message.json"))
	// NapCat 私聊 target_id 恒为会话对端(取自 peerUin),收到的私聊里对端就是发送者,
	// 因此与 user_id 相等;SnowLuma 只在自发消息上报,由管线补成同一语义。
	if msg.TargetID.Int64() != 20002 || msg.UserID.Int64() != 20002 || msg.Sender.Nickname != "Bob" {
		t.Fatalf("私聊消息字段错误: %+v", msg)
	}
	if msg.Message.Text() != "在吗" {
		t.Fatalf("私聊文本 = %q", msg.Message.Text())
	}
}

func TestNapcatForwardInlineGolden(t *testing.T) {
	msg := wantEvent[*event.GroupMessage](t, loadFixture(t, "napcat", "group_forward_inline.json"))
	fd, ok := msg.Message[0].Data.(message.ForwardData)
	if !ok || fd.ID != "napres-1" {
		t.Fatalf("forward 段解析错误: %+v", msg.Message[0].Data)
	}
	if len(fd.Content) != 2 || fd.FetchErr != nil {
		t.Fatalf("内嵌 forward 内容解析错误: %d 个节点, err=%v", len(fd.Content), fd.FetchErr)
	}
	if fd.Content[0].UserID.Int64() != 20001 || fd.Content[0].Nickname != "Alice" || fd.Content[0].Message.Text() != "内嵌第一条" {
		t.Fatalf("forward 节点 0 错误: %+v", fd.Content[0])
	}
	// 节点级 content 字符串形态也应按 OneBot 约定解码为文本段
	if fd.Content[1].Message.Text() != "字符串也接受" {
		t.Fatalf("forward 节点 1 消息体错误: %+v", fd.Content[1].Message)
	}
}

func TestNapcatPokeGolden(t *testing.T) {
	poke := wantEvent[*event.Poke](t, loadFixture(t, "napcat", "poke_group.json"))
	if poke.GroupID.Int64() != 30001 || poke.UserID.Int64() != 20001 || poke.TargetID.Int64() != 10000 || poke.Action != "戳一戳" {
		t.Fatalf("群戳一戳字段错误: %+v", poke)
	}
	if poke.SenderID != 0 || poke.ActionImgURL != "" {
		t.Fatalf("群戳一戳不应携带好友侧字段: %+v", poke)
	}
}

func TestFriendRequestGolden(t *testing.T) {
	for _, backend := range []string{"napcat", "snowluma"} {
		req := wantEvent[*event.FriendRequest](t, loadFixture(t, backend, "friend_request.json"))
		if req.UserID == 0 || req.Flag == "" || req.Comment == "" {
			t.Fatalf("%s 好友请求字段错误: %+v", backend, req)
		}
	}
}

func TestSnowlumaGroupMessageGolden(t *testing.T) {
	msg := wantEvent[*event.GroupMessage](t, loadFixture(t, "snowluma", "group_message.json"))
	// SnowLuma message_id 是有符号 int32,负值是常态,SDK 只按不透明 ID 对待
	if msg.MessageID.Int64() != -1234567890 {
		t.Fatalf("负数 message_id 应原样保留: %d", msg.MessageID.Int64())
	}
	if msg.MessageSeq != 42 {
		t.Fatalf("message_seq = %d", msg.MessageSeq)
	}
	if d, ok := msg.Message[0].Data.(message.AtData); !ok || d.QQ != "10000" {
		t.Fatalf("at 段解析错误: %+v", msg.Message[0].Data)
	}
}

func TestSnowlumaPrivateMessageGolden(t *testing.T) {
	msg := wantEvent[*event.PrivateMessage](t, loadFixture(t, "snowluma", "private_message.json"))
	// target_id 缺席为零值,由归一化管线补齐(管线测试单独覆盖)
	if msg.TargetID != 0 || msg.UserID.Int64() != 20005 {
		t.Fatalf("私聊消息字段错误: target_id=%d user_id=%d", msg.TargetID.Int64(), msg.UserID.Int64())
	}
}

func TestSnowlumaForwardIDOnlyGolden(t *testing.T) {
	msg := wantEvent[*event.GroupMessage](t, loadFixture(t, "snowluma", "group_forward_id_only.json"))
	fd, ok := msg.Message[0].Data.(message.ForwardData)
	if !ok || fd.ID != "resXYZ" || len(fd.Content) != 0 || fd.FetchErr != nil {
		t.Fatalf("id-only forward 段解析错误: %+v", msg.Message[0].Data)
	}
}

func TestSnowlumaPokeFriendGolden(t *testing.T) {
	poke := wantEvent[*event.Poke](t, loadFixture(t, "snowluma", "poke_friend.json"))
	if poke.GroupID != 0 || poke.UserID.Int64() != 20005 || poke.SenderID.Int64() != 20005 || poke.ActionImgURL != "https://example.com/poke.png" {
		t.Fatalf("好友戳一戳字段错误: %+v", poke)
	}
}

func TestSnowlumaGroupUploadGolden(t *testing.T) {
	upload := wantEvent[*event.GroupUpload](t, loadFixture(t, "snowluma", "group_upload.json"))
	if upload.GroupID.Int64() != 30002 || upload.File.ID != "file-1" || upload.File.Name != "合并转发.json" || upload.File.Size != 1234 || upload.File.BusID != 102 {
		t.Fatalf("群文件上传字段错误: %+v", upload.File)
	}
}

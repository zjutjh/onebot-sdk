package onebot_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/message"
)

// wsURL 把 httptest HTTP 地址换成 ws 协议。
func wsURL(server *httptest.Server) string {
	return "ws" + strings.TrimPrefix(server.URL, "http")
}

// TestDialectDetectionBranches 覆盖三种 app_name 检测结果。
func TestDialectDetectionBranches(t *testing.T) {
	cases := []struct {
		appName string
		want    string
	}{
		{"NapCat.Onebot", "napcat"},
		{"SnowLuma", "snowluma"},
		{"MysteryBackend", "generic"},
	}
	for _, tc := range cases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := websocket.Accept(w, r, nil)
			if err != nil {
				return
			}
			defer conn.CloseNow()
			if err := respondDetect(context.Background(), conn, tc.appName); err != nil {
				return
			}
			frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
			if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
				return
			}
			// 持续读取避免读端半关闭干扰
			for {
				if _, _, err := readAction(context.Background(), conn); err != nil {
					return
				}
			}
		}))
		client, err := onebot.DialWebSocket(context.Background(), wsURL(server))
		if err != nil {
			server.Close()
			t.Fatal(err)
		}
		if got := client.Dialect(context.Background()).Name; got != tc.want {
			t.Fatalf("%s 检测为 %q,期望 %q", tc.appName, got, tc.want)
		}
		select {
		case _, ok := <-client.Events():
			if !ok {
				t.Fatalf("%s 检测后事件通道意外关闭", tc.appName)
			}
		case <-time.After(time.Second):
			t.Fatalf("%s 检测后未收到事件", tc.appName)
		}
		_ = client.Close()
		server.Close()
	}
}

// TestBackendOverrideSkipsDetection 显式指定后端时不发起 get_version_info。
func TestBackendOverrideSkipsDetection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		// 不读 action:若客户端仍发起检测,检测会超时降级为 generic
		frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
		if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
			return
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}))
	defer server.Close()

	client, err := onebot.DialWebSocket(
		context.Background(), wsURL(server),
		onebot.WithBackend(onebot.BackendNapCat),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := client.Dialect(context.Background()).Name; got != "napcat" {
		t.Fatalf("显式后端被忽略,方言 = %q", got)
	}
	count := 0
	for range client.Events() {
		count++
	}
	if count != 1 {
		t.Fatalf("收到 %d 个事件,期望 1", count)
	}
}

// snowlumaForwardServer 模拟 SnowLuma:检测后下发仅含 id 的 forward 私聊消息,
// 对 get_forward_msg 按 fail 参数决定返回节点数组或错误。
func snowlumaForwardServer(t *testing.T, fail bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		if err := respondDetect(context.Background(), conn, "SnowLuma"); err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"message","message_type":"private","self_id":10,` +
			`"message_id":-42,"user_id":20005,` +
			`"message":[{"type":"forward","data":{"id":"resXYZ"}}],"raw_message":"","sender":{"user_id":20005,"nickname":"Eve"}}`)
		if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
			return
		}
		action, echo, err := readAction(context.Background(), conn)
		if err != nil {
			return
		}
		if action != "get_forward_msg" {
			t.Errorf("补拉 action = %s,期望 get_forward_msg", action)
			return
		}
		if fail {
			resp := `{"status":"failed","retcode":1404,"message":"not found","data":null,"echo":"` + echo + `"}`
			_ = conn.Write(context.Background(), websocket.MessageText, []byte(resp))
		} else {
			data := `{"messages":[{"user_id":20005,"sender":{"nickname":"Eve"},"time":1,` +
				`"message":[{"type":"text","data":{"text":"补拉内容"}}]}]}`
			if err := writeActionResp(context.Background(), conn, echo, data); err != nil {
				return
			}
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}))
}

func TestForwardEnrichSuccess(t *testing.T) {
	server := snowlumaForwardServer(t, false)
	defer server.Close()

	client, err := onebot.DialWebSocket(context.Background(), wsURL(server))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case ev, ok := <-client.Events():
		if !ok {
			t.Fatal("事件通道提前关闭")
		}
		msg, ok := ev.(*event.PrivateMessage)
		if !ok {
			t.Fatalf("事件类型 = %T", ev)
		}
		// 同一事件同时验证 target_id 补齐与 forward 补拉
		if msg.TargetID.Int64() != 20005 {
			t.Fatalf("私聊 target_id 未补齐: %d", msg.TargetID.Int64())
		}
		fd, ok := msg.Message[0].Data.(message.ForwardData)
		if !ok || fd.ID != "resXYZ" || len(fd.Content) != 1 || fd.FetchErr != nil {
			t.Fatalf("forward 补拉结果错误: %+v", msg.Message[0].Data)
		}
		if fd.Content[0].Message.Text() != "补拉内容" {
			t.Fatalf("补拉节点内容错误: %+v", fd.Content[0])
		}
	case <-time.After(time.Second):
		t.Fatal("未收到事件")
	}
}

func TestForwardFetchFailureDeliversEvent(t *testing.T) {
	server := snowlumaForwardServer(t, true)
	defer server.Close()

	client, err := onebot.DialWebSocket(context.Background(), wsURL(server))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case ev, ok := <-client.Events():
		if !ok {
			t.Fatal("事件通道提前关闭")
		}
		msg := ev.(*event.PrivateMessage)
		fd, ok := msg.Message[0].Data.(message.ForwardData)
		if !ok || fd.ID != "resXYZ" || len(fd.Content) != 0 || fd.FetchErr == nil {
			t.Fatalf("补拉失败应保留 id 并暴露错误: %+v", msg.Message[0].Data)
		}
	case <-time.After(time.Second):
		t.Fatal("补拉失败不应吞掉事件")
	}
}

// TestInlineForwardSkipsFetch NapCat 内嵌内容不走补拉。
func TestInlineForwardSkipsFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		if err := respondDetect(context.Background(), conn, "NapCat.Onebot"); err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"message","message_type":"group","self_id":10,` +
			`"message_id":777,"group_id":30001,"user_id":20001,` +
			`"message":[{"type":"forward","data":{"id":"napres-1","content":[` +
			`{"user_id":20001,"sender":{"nickname":"Alice"},"time":1,"message":[{"type":"text","data":{"text":"内嵌"}}]}]}}],` +
			`"raw_message":"","sender":{"user_id":20001,"nickname":"Alice"}}`)
		if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
			return
		}
		// 收到任何 action(除检测外)即说明错误地发起了补拉
		action, _, err := readAction(context.Background(), conn)
		if err == nil && action != "get_version_info" {
			t.Errorf("内嵌 forward 不应补拉,却收到 action %s", action)
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}))
	defer server.Close()

	client, err := onebot.DialWebSocket(context.Background(), wsURL(server))
	if err != nil {
		t.Fatal(err)
	}
	ev := <-client.Events()
	msg := ev.(*event.GroupMessage)
	fd, ok := msg.Message[0].Data.(message.ForwardData)
	if !ok || len(fd.Content) != 1 || fd.Content[0].Message.Text() != "内嵌" {
		t.Fatalf("内嵌 forward 解析错误: %+v", msg.Message[0].Data)
	}
}

// TestRawEventsBypassTeesFrames 旁路与归一化流并存。
func TestRawEventsBypassTeesFrames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		if err := respondDetect(context.Background(), conn, "SnowLuma"); err != nil {
			return
		}
		// 收到 ping action 后才下发事件,保证客户端已先启用旁路
		if _, _, err := readAction(context.Background(), conn); err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"notice","notice_type":"group_increase","self_id":10,` +
			`"group_id":30001,"user_id":20007}`)
		if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
			return
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}))
	defer server.Close()

	client, err := onebot.DialWebSocket(context.Background(), wsURL(server))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// 先启用旁路,再触发事件下发
	rawCh := client.RawEvents()
	var callDone = make(chan error, 1)
	go func() {
		callDone <- client.Call(context.Background(), "ping", nil, nil)
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	var rawCount, evCount int
	go func() {
		defer wg.Done()
		for range rawCh {
			rawCount++
		}
	}()
	go func() {
		defer wg.Done()
		for range client.Events() {
			evCount++
		}
	}()
	wg.Wait()
	if rawCount != 1 || evCount != 1 {
		t.Fatalf("旁路 %d / 归一化 %d 个事件,期望各 1", rawCount, evCount)
	}
}

// TestHTTPDialectAndFriendRemark HTTP 懒检测方言 + remark 按方言取舍。
func TestHTTPDialectAndFriendRemark(t *testing.T) {
	cases := []struct {
		appName     string
		wantDialect string
		wantRemark  bool
	}{
		{"NapCat.Onebot", "napcat", true},
		{"SnowLuma", "snowluma", false},
	}
	for _, tc := range cases {
		var gotParams map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "get_version_info") {
				_, _ = io.WriteString(w, `{"status":"ok","retcode":0,"data":{"app_name":`+quoteJSON(tc.appName)+`}}`)
				return
			}
			body, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(body, &gotParams)
			_, _ = io.WriteString(w, `{"status":"ok","retcode":0,"data":null}`)
		}))
		client := onebot.NewHTTPClient(server.URL)
		if err := client.SetFriendRequest(context.Background(), "flag-1", true, "备注名"); err != nil {
			t.Fatal(err)
		}
		if got := client.Dialect(context.Background()).Name; got != tc.wantDialect {
			t.Fatalf("%s 懒检测为 %q,期望 %q", tc.appName, got, tc.wantDialect)
		}
		_, hasRemark := gotParams["remark"]
		if hasRemark != tc.wantRemark {
			t.Fatalf("%s remark = %v,期望携带 %v", tc.appName, gotParams, tc.wantRemark)
		}
		server.Close()
	}
}

func quoteJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// TestMessageReactionsDualPaths 表情回应双后端路径:NapCat 走 get_msg 的
// emoji_likes_list(likes_cnt 为字符串),SnowLuma 走独立 action get_msg_emoji_likes。
func TestMessageReactionsDualPaths(t *testing.T) {
	cases := []struct {
		appName    string
		wantAction string
		respData   string
		wantCount  int64
	}{
		{"SnowLuma", "/get_msg_emoji_likes", `[{"emoji_id":"76","count":3}]`, 3},
		{"NapCat.Onebot", "/get_msg", `{"emoji_likes_list":[{"emoji_id":"76","likes_cnt":"2"}]}`, 2},
	}
	for _, tc := range cases {
		var lastPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lastPath = r.URL.Path
			if strings.HasSuffix(r.URL.Path, "get_version_info") {
				_, _ = io.WriteString(w, `{"status":"ok","retcode":0,"data":{"app_name":`+quoteJSON(tc.appName)+`}}`)
				return
			}
			_, _ = io.WriteString(w, `{"status":"ok","retcode":0,"data":`+tc.respData+`}`)
		}))
		client := onebot.NewHTTPClient(server.URL)
		reactions, err := client.GetMessageReactions(context.Background(), -123)
		if err != nil {
			t.Fatal(err)
		}
		if lastPath != tc.wantAction {
			t.Fatalf("%s 调用了 %s,期望 %s", tc.appName, lastPath, tc.wantAction)
		}
		if len(reactions) != 1 || reactions[0].EmojiID != "76" || reactions[0].Count != tc.wantCount {
			t.Fatalf("%s 表情回应归一错误: %+v", tc.appName, reactions)
		}
		server.Close()
	}

	// 通用方言明确报错,不猜测
	client := onebot.NewHTTPClient("http://127.0.0.1:1", onebot.WithBackend(onebot.BackendGeneric))
	if _, err := client.GetMessageReactions(context.Background(), 1); err == nil {
		t.Fatal("通用方言应拒绝读取表情回应")
	}
}

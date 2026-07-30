package napcat_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	napcat "github.com/zjutjh/napcat-sdk"
	"github.com/zjutjh/napcat-sdk/event"
)

func TestHTTPCallAndAPIError(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"failed","retcode":42,"message":"denied","data":null}`))
	}))
	defer server.Close()

	client := napcat.NewHTTPClient(server.URL, napcat.WithHTTPClient(&http.Client{}))
	err := client.Call(context.Background(), "send_msg", nil, nil)
	var apiErr *napcat.APIError
	if !errors.As(err, &apiErr) || !strings.Contains(err.Error(), "send_msg") || !strings.Contains(err.Error(), "retcode 42") {
		t.Fatalf("API 错误缺少 action 或 retcode: %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("HTTP action 执行了 %d 次，期望 1", calls.Load())
	}
	if client.Err() != nil {
		t.Fatalf("HTTP client Err() = %v", client.Err())
	}
}

func TestEventIDAndRawFallback(t *testing.T) {
	valid := map[string]int64{
		`0`:                      0,
		`-1`:                     -1,
		`"9223372036854775807"`:  9223372036854775807,
		`"-9223372036854775808"`: -9223372036854775808,
	}
	for input, want := range valid {
		var id event.ID
		if err := id.UnmarshalJSON([]byte(input)); err != nil {
			t.Fatalf("ID %s 解析失败: %v", input, err)
		}
		if got := id.Int64(); got != want {
			t.Fatalf("ID %s = %d，期望 %d", input, got, want)
		}
	}
	for _, input := range []string{`null`, `true`, `[]`, `{}`, `1.0`, `1e3`, `""`, `"1.0"`, `"1e3"`, `"abc"`, `9223372036854775808`, `"9223372036854775808"`} {
		var id event.ID
		if err := id.UnmarshalJSON([]byte(input)); err == nil {
			t.Fatalf("ID %s 应解析失败", input)
		}
	}

	raw := []byte(`{"time":1,"post_type":"message","message_type":"group","self_id":"10","message_id":"20","group_id":"30","user_id":"40","message":[],"raw_message":"","sender":{"user_id":"40","nickname":"n"}}`)
	message, ok := event.Parse(raw).(*event.GroupMessage)
	if !ok {
		t.Fatalf("字符串 ID 应解析为群消息，得到 %T", event.Parse(raw))
	}
	if message.SelfID() != 10 || message.MessageID.Int64() != 20 || message.GroupID.Int64() != 30 || message.UserID.Int64() != 40 || message.Sender.UserID.Int64() != 40 {
		t.Fatalf("群消息 ID 解析错误: %+v", message)
	}

	malformed := []byte(`{"time":2,"post_type":"message","message_type":"group","self_id":"10","message_id":"bad"}`)
	failure, ok := event.Parse(malformed).(event.ParseFailure)
	if !ok || failure.ParseError() == nil {
		t.Fatalf("畸形群消息应保留为 ParseFailure，得到 %T", event.Parse(malformed))
	}
	if failure.Time() != 2 || failure.PostType() != "message" || failure.SelfID() != 10 || string(failure.Raw()) != string(malformed) {
		t.Fatalf("RawEvent 未完整保留 envelope 或原始 JSON")
	}

	badSelfID := []byte(`{"self_id":null,"post_type":"message","time":3,"message_type":"group"}`)
	failure, ok = event.Parse(badSelfID).(event.ParseFailure)
	if !ok || failure.Time() != 3 || failure.PostType() != "message" {
		t.Fatalf("RawEvent 应独立保留可解析的 envelope 字段，得到 %T", event.Parse(badSelfID))
	}
}

func TestEventBackpressureTerminatesPendingCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		if _, _, err := conn.Read(context.Background()); err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
		for range 3 {
			if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
				return
			}
		}
		_, _, _ = conn.Read(context.Background())
	}))
	defer server.Close()

	client, err := napcat.DialWebSocket(
		context.Background(),
		"ws"+strings.TrimPrefix(server.URL, "http"),
		napcat.WithEventBuffer(1),
		napcat.WithEventDeliveryTimeout(50*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	callDone := make(chan error, 1)
	go func() {
		callDone <- client.Call(context.Background(), "blocked", nil, nil)
	}()
	select {
	case err := <-callDone:
		if !errors.Is(err, napcat.ErrEventBackpressure) {
			t.Fatalf("pending 调用错误 = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("事件背压未唤醒 pending 调用")
	}
	if !errors.Is(client.Err(), napcat.ErrEventBackpressure) {
		t.Fatalf("Client.Err() = %v", client.Err())
	}
	if err := client.Close(); err != nil {
		t.Fatalf("重复关闭失败: %v", err)
	}
	if !errors.Is(client.Err(), napcat.ErrEventBackpressure) {
		t.Fatalf("Close 覆盖了首次终止原因: %v", client.Err())
	}

	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-client.Events():
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("事件 channel 未由读循环关闭")
		}
	}
}

func TestNormalCloseDrainsBufferedEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
		for range 2 {
			if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
				conn.CloseNow()
				return
			}
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}))
	defer server.Close()

	client, err := napcat.DialWebSocket(
		context.Background(),
		"ws"+strings.TrimPrefix(server.URL, "http"),
		napcat.WithEventBuffer(2),
	)
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for range client.Events() {
		count++
	}
	if count != 2 {
		t.Fatalf("正常关闭后收到 %d 个缓冲事件，期望 2", count)
	}
	if !errors.Is(client.Err(), napcat.ErrClosed) {
		t.Fatalf("正常关闭原因 = %v", client.Err())
	}
}

func TestDefaultEventDropDoesNotBlockActionResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
		for range 3 {
			if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
				return
			}
		}
		if _, _, err := conn.Read(context.Background()); err != nil {
			return
		}
		response := []byte(`{"status":"ok","retcode":0,"data":null,"echo":"1"}`)
		if err := conn.Write(context.Background(), websocket.MessageText, response); err != nil {
			return
		}
		_, _, _ = conn.Read(context.Background())
	}))
	defer server.Close()

	client, err := napcat.DialWebSocket(
		context.Background(),
		"ws"+strings.TrimPrefix(server.URL, "http"),
		napcat.WithEventBuffer(1),
		napcat.WithRequestTimeout(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if err := client.Call(context.Background(), "action", nil, nil); err != nil {
		t.Fatalf("默认事件丢弃模式阻塞了 action 响应: %v", err)
	}
	if err := client.Err(); err != nil {
		t.Fatalf("默认事件丢弃不应终止连接: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("主动关闭失败: %v", err)
	}
	if !errors.Is(client.Err(), napcat.ErrClosed) {
		t.Fatalf("主动关闭原因 = %v", client.Err())
	}
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-client.Events():
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("主动关闭后事件 channel 未关闭")
		}
	}
}

func TestActionResponseWinsConcurrentClose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		if _, _, err := conn.Read(context.Background()); err != nil {
			return
		}
		response := []byte(`{"status":"ok","retcode":0,"data":null,"echo":"1"}`)
		if err := conn.Write(context.Background(), websocket.MessageText, response); err != nil {
			return
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}))
	defer server.Close()

	for range 10 {
		client, err := napcat.DialWebSocket(context.Background(), "ws"+strings.TrimPrefix(server.URL, "http"))
		if err != nil {
			t.Fatal(err)
		}
		if err := client.Call(context.Background(), "action", nil, nil); err != nil {
			t.Fatalf("已收到的 action 响应被关闭原因覆盖: %v", err)
		}
		_ = client.Close()
	}
}

package onebot_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/zjutjh/onebot-sdk"
	"github.com/zjutjh/onebot-sdk/event"
)

// readAction 读取一帧 action 请求,返回 action 与 echo。
func readAction(ctx context.Context, conn *websocket.Conn) (action, echo string, err error) {
	_, data, err := conn.Read(ctx)
	if err != nil {
		return "", "", err
	}
	var req struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
	}
	if err := json.Unmarshal(data, &req); err != nil || req.Echo == "" {
		return "", "", fmt.Errorf("非 action 帧: %s", data)
	}
	return req.Action, req.Echo, nil
}

// writeActionResp 按 echo 写回成功响应,data 为 data 字段的 JSON 文本。
func writeActionResp(ctx context.Context, conn *websocket.Conn, echo, data string) error {
	resp := fmt.Sprintf(`{"status":"ok","retcode":0,"data":%s,"echo":%q}`, data, echo)
	return conn.Write(ctx, websocket.MessageText, []byte(resp))
}

// respondDetect 响应方言检测请求,app_name 决定客户端识别出的方言。
func respondDetect(ctx context.Context, conn *websocket.Conn, appName string) error {
	action, echo, err := readAction(ctx, conn)
	if err != nil {
		return err
	}
	if action != "get_version_info" {
		return fmt.Errorf("首帧应为方言检测,得到 %s", action)
	}
	return writeActionResp(ctx, conn, echo, fmt.Sprintf(`{"app_name":%q}`, appName))
}

func TestHTTPCallAndAPIError(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"failed","retcode":42,"message":"denied","data":null}`))
	}))
	defer server.Close()

	client := onebot.NewHTTPClient(server.URL, onebot.WithHTTPClient(&http.Client{}))
	err := client.Call(context.Background(), "send_msg", nil, nil)
	var apiErr *onebot.APIError
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
		if err := respondDetect(context.Background(), conn, "SnowLuma"); err != nil {
			return
		}
		// 先等一次 action 请求再投事件:此时管线已在消费,事件能进入用户通道并被有界队列挡住;
		// 否则洪峰在管线接管前就被传输队列丢弃,填不满用户通道,验证不到背压
		if _, _, err := readAction(context.Background(), conn); err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
		for range 3 {
			if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
				return
			}
		}
		// 持续读取,吞掉 blocked 请求但不响应
		for {
			if _, _, err := readAction(context.Background(), conn); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client, err := onebot.DialWebSocket(
		context.Background(),
		wsURL(server),
		onebot.WithEventBuffer(1),
		onebot.WithEventDeliveryTimeout(50*time.Millisecond),
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
		if !errors.Is(err, onebot.ErrEventBackpressure) {
			t.Fatalf("pending 调用错误 = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("事件背压未唤醒 pending 调用")
	}
	if !errors.Is(client.Err(), onebot.ErrEventBackpressure) {
		t.Fatalf("Client.Err() = %v", client.Err())
	}
	if err := client.Close(); err != nil {
		t.Fatalf("重复关闭失败: %v", err)
	}
	if !errors.Is(client.Err(), onebot.ErrEventBackpressure) {
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
			t.Fatal("事件 channel 未由管线关闭")
		}
	}
}

func TestNormalCloseDrainsBufferedEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		if err := respondDetect(context.Background(), conn, "SnowLuma"); err != nil {
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

	client, err := onebot.DialWebSocket(
		context.Background(),
		wsURL(server),
		onebot.WithEventBuffer(2),
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
	if !errors.Is(client.Err(), onebot.ErrClosed) {
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
		if err := respondDetect(context.Background(), conn, "SnowLuma"); err != nil {
			return
		}
		frame := []byte(`{"time":1,"post_type":"notice","self_id":"1"}`)
		for range 3 {
			if err := conn.Write(context.Background(), websocket.MessageText, frame); err != nil {
				return
			}
		}
		for {
			action, echo, err := readAction(context.Background(), conn)
			if err != nil {
				return
			}
			if err := writeActionResp(context.Background(), conn, echo, actionData(action)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client, err := onebot.DialWebSocket(
		context.Background(),
		wsURL(server),
		onebot.WithEventBuffer(1),
		onebot.WithRequestTimeout(time.Second),
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
	if !errors.Is(client.Err(), onebot.ErrClosed) {
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

// TestHTTPDialectDetectionRetriesAfterFailure 覆盖 HTTP 懒检测失败:
// 临时失败不得被永久缓存,后端恢复后必须重新识别出真实方言。
func TestHTTPDialectDetectionRetriesAfterFailure(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if calls.Add(1) == 1 {
			_, _ = w.Write([]byte(`{"status":"failed","retcode":1,"message":"busy","data":null}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"ok","retcode":0,"data":{"app_name":"NapCat.Onebot"}}`))
	}))
	defer server.Close()

	client := onebot.NewHTTPClient(server.URL, onebot.WithHTTPClient(&http.Client{}))
	if name := client.Dialect(context.Background()).Name; name != "generic" {
		t.Fatalf("首次检测失败应降级为 generic,得到 %q", name)
	}
	if name := client.Dialect(context.Background()).Name; name != "napcat" {
		t.Fatalf("检测失败不得被缓存,重试后应识别为 napcat,得到 %q", name)
	}
	if name := client.Dialect(context.Background()).Name; name != "napcat" || calls.Load() != 2 {
		t.Fatalf("识别结果应缓存,第三次调用得到 %q、请求数 %d", name, calls.Load())
	}
}

// actionData 返回 action 对应的响应 data,用于按 action 名分发模拟响应。
func actionData(action string) string {
	switch action {
	case "get_version_info":
		return `{"app_name":"SnowLuma"}`
	default:
		return `null`
	}
}

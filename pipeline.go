package onebot

import (
	"context"
	"time"

	"github.com/zjutjh/onebot-sdk/event"
	"github.com/zjutjh/onebot-sdk/internal/errorsx"
	"github.com/zjutjh/onebot-sdk/message"
)

// forwardFetchTimeout 是归一化管线补拉合并转发内容的独立超时,
// 不计入事件交付超时,避免补拉往返误杀连接。
const forwardFetchTimeout = 10 * time.Second

// runPipeline 是事件归一化管线:单 goroutine 串行处理原始事件,
// 依次完成强类型解析、方言归一化和投递,保证投递顺序与到达顺序一致。
// 管线消费传输层主原始队列(TakeRawEvents),用户旁路 TeeRawEvents 与其并存,
// 管线退出时关闭用户事件通道。
func (c *Client) runPipeline() {
	defer close(c.events)
	// 管线自带 WS 连接,构造时已同步解析方言,读取只取缓存值
	forwardFetch := c.currentDialect(context.Background()).ForwardFetch
	for raw := range c.ws.TakeRawEvents() {
		ev := event.Parse(raw)
		if failure, ok := ev.(event.ParseFailure); ok {
			if c.opts.logger != nil {
				c.opts.logger.Warn("事件强类型解析失败，已保留原始事件", "error", failure.ParseError())
			}
		}
		c.normalize(ev, forwardFetch)
		if !c.deliverEvent(ev) {
			return
		}
	}
}

// deliverEvent 把事件投递给用户,背压语义在这一跳生效:
// 交付超时为正时,用户持续不消费会在超时后以 ErrEventBackpressure 终止连接;
// 未配置时满队列直接丢弃并记日志,不影响后续事件。
// 连接终止时立即停止投递,避免 Close 后管线还停在交付等待里。
func (c *Client) deliverEvent(ev event.Event) bool {
	if c.opts.eventDeliveryTimeout <= 0 {
		select {
		case c.events <- ev:
			return true
		default:
		}
		if c.opts.logger != nil {
			c.opts.logger.Warn("事件缓冲区已满，丢弃事件", "post_type", ev.PostType())
		}
		return true
	}

	timer := time.NewTimer(c.opts.eventDeliveryTimeout)
	defer timer.Stop()
	select {
	case c.events <- ev:
		return true
	case <-c.ws.Done():
		return false
	case <-timer.C:
		if c.opts.logger != nil {
			c.opts.logger.Warn("事件交付超时，终止连接", "timeout", c.opts.eventDeliveryTimeout)
		}
		_ = c.ws.Terminate(errorsx.ErrEventBackpressure)
		return false
	}
}

// normalize 按方言补齐事件字段缺口,forwardFetch 表示是否补拉只有 id 的合并转发段。
func (c *Client) normalize(ev event.Event, forwardFetch bool) {
	switch e := ev.(type) {
	case *event.PrivateMessage:
		// SnowLuma 收到的私聊不带 target_id,按对端补齐
		if e.TargetID == 0 {
			e.TargetID = e.UserID
		}
		c.checkMessageID(e.MessageID)
		if forwardFetch {
			c.enrichForward(e.Message)
		}
	case *event.GroupMessage:
		c.checkMessageID(e.MessageID)
		if forwardFetch {
			c.enrichForward(e.Message)
		}
	}
}

func (c *Client) checkMessageID(id event.ID) {
	if id == 0 && c.opts.logger != nil {
		c.opts.logger.Warn("事件 message_id 缺失，该消息无法用于回复、撤回或去重")
	}
}

// enrichForward 补拉事件内只有 id 的合并转发段。
// 补拉失败时保留原 id 并把错误挂在 FetchErr 上,不伪造内容。
func (c *Client) enrichForward(chain message.Chain) {
	for i, seg := range chain {
		fd, ok := seg.Data.(message.ForwardData)
		if !ok {
			continue
		}
		if len(fd.Content) > 0 || fd.FetchErr != nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), forwardFetchTimeout)
		nodes, err := c.GetForwardContent(ctx, fd.ID)
		cancel()
		if err != nil {
			chain[i].Data = message.ForwardData{ID: fd.ID, FetchErr: err}
		} else {
			chain[i].Data = message.ForwardData{ID: fd.ID, Content: nodes}
		}
	}
}

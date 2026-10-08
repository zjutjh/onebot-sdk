package onebot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zjutjh/onebot-sdk/message"
)

// GetForwardContent 按转发 ID(resId 或数值消息 ID 字符串)拉取合并转发消息内容。
// 归一化管线对 SnowLuma 事件自动调用;其他场景可直接复用。
func (c *Client) GetForwardContent(ctx context.Context, id string) ([]message.ForwardNode, error) {
	var res struct {
		Messages []message.ForwardNode `json:"messages"`
	}
	params := struct {
		ID string `json:"id"`
	}{ID: id}
	if err := c.Call(ctx, "get_forward_msg", params, &res); err != nil {
		return nil, fmt.Errorf("拉取合并转发内容失败: %w", err)
	}
	return res.Messages, nil
}

// EmojiReaction 是归一化后的消息表情回应。
type EmojiReaction struct {
	EmojiID string
	Count   int64
}

// GetMessageReactions 读取消息的表情回应。
// NapCat 走 get_msg 的 emoji_likes_list,SnowLuma 走独立的 get_msg_emoji_likes;
// 未知后端不支持该能力,返回错误而不是猜测。
func (c *Client) GetMessageReactions(ctx context.Context, messageID int64) ([]EmojiReaction, error) {
	d := c.currentDialect(ctx)
	switch {
	case d.ReactionsViaGetMsg:
		var res struct {
			EmojiLikesList []struct {
				EmojiID  string `json:"emoji_id"`
				LikesCnt string `json:"likes_cnt"`
			} `json:"emoji_likes_list"`
		}
		if err := c.Call(ctx, "get_msg", map[string]any{"message_id": messageID}, &res); err != nil {
			return nil, fmt.Errorf("读取消息表情回应失败: %w", err)
		}
		out := make([]EmojiReaction, 0, len(res.EmojiLikesList))
		for _, item := range res.EmojiLikesList {
			count, err := strconv.ParseInt(item.LikesCnt, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("解析表情回应数量失败: %w", err)
			}
			out = append(out, EmojiReaction{EmojiID: item.EmojiID, Count: count})
		}
		return out, nil
	case d.EmojiLikesAction:
		var res []struct {
			EmojiID string `json:"emoji_id"`
			Count   int64  `json:"count"`
		}
		if err := c.Call(ctx, "get_msg_emoji_likes", map[string]any{"message_id": messageID}, &res); err != nil {
			return nil, fmt.Errorf("读取消息表情回应失败: %w", err)
		}
		out := make([]EmojiReaction, 0, len(res))
		for _, item := range res {
			out = append(out, EmojiReaction{EmojiID: item.EmojiID, Count: item.Count})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("当前后端(%s)不支持读取消息表情回应", d.Name)
	}
}

// SetFriendRequest 处理好友请求。
// approve 为 true 时通过请求;remark 仅 NapCat 支持,其余后端静默忽略。
func (c *Client) SetFriendRequest(ctx context.Context, flag string, approve bool, remark string) error {
	d := c.currentDialect(ctx)
	params := map[string]any{
		"flag":    flag,
		"approve": approve,
	}
	if remark != "" && d.FriendRemark {
		params["remark"] = remark
	}
	return c.Call(ctx, "set_friend_add_request", params, nil)
}

// 代码由 napcatgen 根据 NapCat OpenAPI 4.18.33 生成，请勿手动修改.

package api

import "context"

// DotHandleQuickOperation 处理快速操作
func (c *Client) DotHandleQuickOperation(ctx context.Context, req DotHandleQuickOperationRequest) (*DotHandleQuickOperationResponse, error) {
	var out DotHandleQuickOperationResponse
	if err := c.caller.Call(ctx, string(ActionDotHandleQuickOperation), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DotOcrImage 图片 OCR 识别 (内部)
func (c *Client) DotOcrImage(ctx context.Context, req DotOcrImageRequest) (*DotOcrImageResponse, error) {
	var out DotOcrImageResponse
	if err := c.caller.Call(ctx, string(ActionDotOcrImage), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ArkShareGroup 分享群 (Ark)
func (c *Client) ArkShareGroup(ctx context.Context, req ArkShareGroupRequest) (*ArkShareGroupResponse, error) {
	var out ArkShareGroupResponse
	if err := c.caller.Call(ctx, string(ActionArkShareGroup), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ArkSharePeer 分享用户 (Ark)
func (c *Client) ArkSharePeer(ctx context.Context, req ArkSharePeerRequest) (*ArkSharePeerResponse, error) {
	var out ArkSharePeerResponse
	if err := c.caller.Call(ctx, string(ActionArkSharePeer), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnderscoreDelGroupNotice 删除群公告
func (c *Client) UnderscoreDelGroupNotice(ctx context.Context, req UnderscoreDelGroupNoticeRequest) (*UnderscoreDelGroupNoticeResponse, error) {
	var out UnderscoreDelGroupNoticeResponse
	if err := c.caller.Call(ctx, string(ActionUnderscoreDelGroupNotice), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnderscoreGetGroupNotice 获取群公告
func (c *Client) UnderscoreGetGroupNotice(ctx context.Context, req UnderscoreGetGroupNoticeRequest) (*UnderscoreGetGroupNoticeResponse, error) {
	var out UnderscoreGetGroupNoticeResponse
	if err := c.caller.Call(ctx, string(ActionUnderscoreGetGroupNotice), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnderscoreGetModelShow 获取机型显示
func (c *Client) UnderscoreGetModelShow(ctx context.Context, req UnderscoreGetModelShowRequest) (*UnderscoreGetModelShowResponse, error) {
	var out UnderscoreGetModelShowResponse
	if err := c.caller.Call(ctx, string(ActionUnderscoreGetModelShow), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnderscoreMarkAllAsRead 标记所有消息已读
func (c *Client) UnderscoreMarkAllAsRead(ctx context.Context, req UnderscoreMarkAllAsReadRequest) (*UnderscoreMarkAllAsReadResponse, error) {
	var out UnderscoreMarkAllAsReadResponse
	if err := c.caller.Call(ctx, string(ActionUnderscoreMarkAllAsRead), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnderscoreSendGroupNotice 发送群公告
func (c *Client) UnderscoreSendGroupNotice(ctx context.Context, req UnderscoreSendGroupNoticeRequest) (*UnderscoreSendGroupNoticeResponse, error) {
	var out UnderscoreSendGroupNoticeResponse
	if err := c.caller.Call(ctx, string(ActionUnderscoreSendGroupNotice), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnderscoreSetModelShow 设置机型
func (c *Client) UnderscoreSetModelShow(ctx context.Context, req UnderscoreSetModelShowRequest) (*UnderscoreSetModelShowResponse, error) {
	var out UnderscoreSetModelShowResponse
	if err := c.caller.Call(ctx, string(ActionUnderscoreSetModelShow), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddCustomFace 添加自定义表情
func (c *Client) AddCustomFace(ctx context.Context, req AddCustomFaceRequest) (*AddCustomFaceResponse, error) {
	var out AddCustomFaceResponse
	if err := c.caller.Call(ctx, string(ActionAddCustomFace), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BotExit 退出登录
func (c *Client) BotExit(ctx context.Context, req BotExitRequest) (*BotExitResponse, error) {
	var out BotExitResponse
	if err := c.caller.Call(ctx, string(ActionBotExit), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CanSendImage 是否可以发送图片
func (c *Client) CanSendImage(ctx context.Context, req CanSendImageRequest) (*CanSendImageResponse, error) {
	var out CanSendImageResponse
	if err := c.caller.Call(ctx, string(ActionCanSendImage), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CanSendRecord 是否可以发送语音
func (c *Client) CanSendRecord(ctx context.Context, req CanSendRecordRequest) (*CanSendRecordResponse, error) {
	var out CanSendRecordResponse
	if err := c.caller.Call(ctx, string(ActionCanSendRecord), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelGroupAlbumMediaLike 取消点赞群相册媒体
func (c *Client) CancelGroupAlbumMediaLike(ctx context.Context, req CancelGroupAlbumMediaLikeRequest) (*CancelGroupAlbumMediaLikeResponse, error) {
	var out CancelGroupAlbumMediaLikeResponse
	if err := c.caller.Call(ctx, string(ActionCancelGroupAlbumMediaLike), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelGroupTodo 取消群待办
func (c *Client) CancelGroupTodo(ctx context.Context, req CancelGroupTodoRequest) (*CancelGroupTodoResponse, error) {
	var out CancelGroupTodoResponse
	if err := c.caller.Call(ctx, string(ActionCancelGroupTodo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelOnlineFile 取消在线文件
func (c *Client) CancelOnlineFile(ctx context.Context, req CancelOnlineFileRequest) (*CancelOnlineFileResponse, error) {
	var out CancelOnlineFileResponse
	if err := c.caller.Call(ctx, string(ActionCancelOnlineFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CheckURLSafely 检查URL安全性
func (c *Client) CheckURLSafely(ctx context.Context, req CheckURLSafelyRequest) (*CheckURLSafelyResponse, error) {
	var out CheckURLSafelyResponse
	if err := c.caller.Call(ctx, string(ActionCheckURLSafely), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CleanCache 清理缓存
func (c *Client) CleanCache(ctx context.Context, req CleanCacheRequest) (*CleanCacheResponse, error) {
	var out CleanCacheResponse
	if err := c.caller.Call(ctx, string(ActionCleanCache), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CleanStreamTempFile 清理流式传输临时文件
func (c *Client) CleanStreamTempFile(ctx context.Context, req CleanStreamTempFileRequest) (*CleanStreamTempFileResponse, error) {
	var out CleanStreamTempFileResponse
	if err := c.caller.Call(ctx, string(ActionCleanStreamTempFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ClickInlineKeyboardButton 点击内联键盘按钮
func (c *Client) ClickInlineKeyboardButton(ctx context.Context, req ClickInlineKeyboardButtonRequest) (*ClickInlineKeyboardButtonResponse, error) {
	var out ClickInlineKeyboardButtonResponse
	if err := c.caller.Call(ctx, string(ActionClickInlineKeyboardButton), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CompleteGroupTodo 完成群待办
func (c *Client) CompleteGroupTodo(ctx context.Context, req CompleteGroupTodoRequest) (*CompleteGroupTodoResponse, error) {
	var out CompleteGroupTodoResponse
	if err := c.caller.Call(ctx, string(ActionCompleteGroupTodo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCollection 创建收藏
func (c *Client) CreateCollection(ctx context.Context, req CreateCollectionRequest) (*CreateCollectionResponse, error) {
	var out CreateCollectionResponse
	if err := c.caller.Call(ctx, string(ActionCreateCollection), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateFlashTask 创建闪传任务
func (c *Client) CreateFlashTask(ctx context.Context, req CreateFlashTaskRequest) (*CreateFlashTaskResponse, error) {
	var out CreateFlashTaskResponse
	if err := c.caller.Call(ctx, string(ActionCreateFlashTask), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateGroupFileFolder 创建群文件目录
func (c *Client) CreateGroupFileFolder(ctx context.Context, req CreateGroupFileFolderRequest) (*CreateGroupFileFolderResponse, error) {
	var out CreateGroupFileFolderResponse
	if err := c.caller.Call(ctx, string(ActionCreateGroupFileFolder), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DelGroupAlbumMedia 删除群相册媒体
func (c *Client) DelGroupAlbumMedia(ctx context.Context, req DelGroupAlbumMediaRequest) (*DelGroupAlbumMediaResponse, error) {
	var out DelGroupAlbumMediaResponse
	if err := c.caller.Call(ctx, string(ActionDelGroupAlbumMedia), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFace 删除自定义表情
func (c *Client) DeleteCustomFace(ctx context.Context, req DeleteCustomFaceRequest) (*DeleteCustomFaceResponse, error) {
	var out DeleteCustomFaceResponse
	if err := c.caller.Call(ctx, string(ActionDeleteCustomFace), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteEssenceMsg 移出精华消息
func (c *Client) DeleteEssenceMsg(ctx context.Context, req DeleteEssenceMsgRequest) (*DeleteEssenceMsgResponse, error) {
	var out DeleteEssenceMsgResponse
	if err := c.caller.Call(ctx, string(ActionDeleteEssenceMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteFriend 删除好友
func (c *Client) DeleteFriend(ctx context.Context, req DeleteFriendRequest) (*DeleteFriendResponse, error) {
	var out DeleteFriendResponse
	if err := c.caller.Call(ctx, string(ActionDeleteFriend), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteGroupFile 删除群文件
func (c *Client) DeleteGroupFile(ctx context.Context, req DeleteGroupFileRequest) (*DeleteGroupFileResponse, error) {
	var out DeleteGroupFileResponse
	if err := c.caller.Call(ctx, string(ActionDeleteGroupFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteGroupFolder 删除群文件目录
func (c *Client) DeleteGroupFolder(ctx context.Context, req DeleteGroupFolderRequest) (*DeleteGroupFolderResponse, error) {
	var out DeleteGroupFolderResponse
	if err := c.caller.Call(ctx, string(ActionDeleteGroupFolder), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMsg 撤回消息
func (c *Client) DeleteMsg(ctx context.Context, req DeleteMsgRequest) (*DeleteMsgResponse, error) {
	var out DeleteMsgResponse
	if err := c.caller.Call(ctx, string(ActionDeleteMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQzoneMsg 删除QQ空间说说
func (c *Client) DeleteQzoneMsg(ctx context.Context, req DeleteQzoneMsgRequest) (*DeleteQzoneMsgResponse, error) {
	var out DeleteQzoneMsgResponse
	if err := c.caller.Call(ctx, string(ActionDeleteQzoneMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DoGroupAlbumComment 发表群相册评论
func (c *Client) DoGroupAlbumComment(ctx context.Context, req DoGroupAlbumCommentRequest) (*DoGroupAlbumCommentResponse, error) {
	var out DoGroupAlbumCommentResponse
	if err := c.caller.Call(ctx, string(ActionDoGroupAlbumComment), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadFile 下载文件
func (c *Client) DownloadFile(ctx context.Context, req DownloadFileRequest) (*DownloadFileResponse, error) {
	var out DownloadFileResponse
	if err := c.caller.Call(ctx, string(ActionDownloadFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadFileImageStream 下载图片文件流
func (c *Client) DownloadFileImageStream(ctx context.Context, req DownloadFileImageStreamRequest) (*DownloadFileImageStreamResponse, error) {
	var out DownloadFileImageStreamResponse
	if err := c.caller.Call(ctx, string(ActionDownloadFileImageStream), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadFileRecordStream 下载语音文件流
func (c *Client) DownloadFileRecordStream(ctx context.Context, req DownloadFileRecordStreamRequest) (*DownloadFileRecordStreamResponse, error) {
	var out DownloadFileRecordStreamResponse
	if err := c.caller.Call(ctx, string(ActionDownloadFileRecordStream), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadFileStream 下载文件流
func (c *Client) DownloadFileStream(ctx context.Context, req DownloadFileStreamRequest) (*DownloadFileStreamResponse, error) {
	var out DownloadFileStreamResponse
	if err := c.caller.Call(ctx, string(ActionDownloadFileStream), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadFileset 下载文件集
func (c *Client) DownloadFileset(ctx context.Context, req DownloadFilesetRequest) (*DownloadFilesetResponse, error) {
	var out DownloadFilesetResponse
	if err := c.caller.Call(ctx, string(ActionDownloadFileset), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchCustomFace 获取自定义表情
func (c *Client) FetchCustomFace(ctx context.Context, req FetchCustomFaceRequest) (*FetchCustomFaceResponse, error) {
	var out FetchCustomFaceResponse
	if err := c.caller.Call(ctx, string(ActionFetchCustomFace), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchCustomFaceDetail 获取自定义表情详情
func (c *Client) FetchCustomFaceDetail(ctx context.Context, req FetchCustomFaceDetailRequest) (*FetchCustomFaceDetailResponse, error) {
	var out FetchCustomFaceDetailResponse
	if err := c.caller.Call(ctx, string(ActionFetchCustomFaceDetail), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchEmojiLike 获取表情点赞详情
func (c *Client) FetchEmojiLike(ctx context.Context, req FetchEmojiLikeRequest) (*FetchEmojiLikeResponse, error) {
	var out FetchEmojiLikeResponse
	if err := c.caller.Call(ctx, string(ActionFetchEmojiLike), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FetchPttText 获取语音转文字结果
func (c *Client) FetchPttText(ctx context.Context, req FetchPttTextRequest) (*FetchPttTextResponse, error) {
	var out FetchPttTextResponse
	if err := c.caller.Call(ctx, string(ActionFetchPttText), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ForwardFriendSingleMsg 转发单条消息
func (c *Client) ForwardFriendSingleMsg(ctx context.Context, req ForwardFriendSingleMsgRequest) (*ForwardFriendSingleMsgResponse, error) {
	var out ForwardFriendSingleMsgResponse
	if err := c.caller.Call(ctx, string(ActionForwardFriendSingleMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ForwardGroupSingleMsg 转发单条消息
func (c *Client) ForwardGroupSingleMsg(ctx context.Context, req ForwardGroupSingleMsgRequest) (*ForwardGroupSingleMsgResponse, error) {
	var out ForwardGroupSingleMsgResponse
	if err := c.caller.Call(ctx, string(ActionForwardGroupSingleMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FriendPoke 发送戳一戳
func (c *Client) FriendPoke(ctx context.Context, req FriendPokeRequest) (*FriendPokeResponse, error) {
	var out FriendPokeResponse
	if err := c.caller.Call(ctx, string(ActionFriendPoke), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAiCharacters 获取AI角色列表
func (c *Client) GetAiCharacters(ctx context.Context, req GetAiCharactersRequest) (*GetAiCharactersResponse, error) {
	var out GetAiCharactersResponse
	if err := c.caller.Call(ctx, string(ActionGetAiCharacters), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAiRecord 获取 AI 语音
func (c *Client) GetAiRecord(ctx context.Context, req GetAiRecordRequest) (*GetAiRecordResponse, error) {
	var out GetAiRecordResponse
	if err := c.caller.Call(ctx, string(ActionGetAiRecord), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetClientkey 获取ClientKey
func (c *Client) GetClientkey(ctx context.Context, req GetClientkeyRequest) (*GetClientkeyResponse, error) {
	var out GetClientkeyResponse
	if err := c.caller.Call(ctx, string(ActionGetClientkey), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCollectionList 获取收藏列表
func (c *Client) GetCollectionList(ctx context.Context, req GetCollectionListRequest) (*GetCollectionListResponse, error) {
	var out GetCollectionListResponse
	if err := c.caller.Call(ctx, string(ActionGetCollectionList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCookies 获取 Cookies
func (c *Client) GetCookies(ctx context.Context, req GetCookiesRequest) (*GetCookiesResponse, error) {
	var out GetCookiesResponse
	if err := c.caller.Call(ctx, string(ActionGetCookies), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCredentials 获取登录凭证
func (c *Client) GetCredentials(ctx context.Context, req GetCredentialsRequest) (*GetCredentialsResponse, error) {
	var out GetCredentialsResponse
	if err := c.caller.Call(ctx, string(ActionGetCredentials), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCsrfToken 获取 CSRF Token
func (c *Client) GetCsrfToken(ctx context.Context, req GetCsrfTokenRequest) (*GetCsrfTokenResponse, error) {
	var out GetCsrfTokenResponse
	if err := c.caller.Call(ctx, string(ActionGetCsrfToken), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDoubtFriendsAddRequest 获取可疑好友申请
func (c *Client) GetDoubtFriendsAddRequest(ctx context.Context, req GetDoubtFriendsAddRequestRequest) (*GetDoubtFriendsAddRequestResponse, error) {
	var out GetDoubtFriendsAddRequestResponse
	if err := c.caller.Call(ctx, string(ActionGetDoubtFriendsAddRequest), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEmojiLikes 获取消息表情点赞列表
func (c *Client) GetEmojiLikes(ctx context.Context, req GetEmojiLikesRequest) (*GetEmojiLikesResponse, error) {
	var out GetEmojiLikesResponse
	if err := c.caller.Call(ctx, string(ActionGetEmojiLikes), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEssenceMsgList 获取群精华消息
func (c *Client) GetEssenceMsgList(ctx context.Context, req GetEssenceMsgListRequest) (*GetEssenceMsgListResponse, error) {
	var out GetEssenceMsgListResponse
	if err := c.caller.Call(ctx, string(ActionGetEssenceMsgList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFile 获取文件
func (c *Client) GetFile(ctx context.Context, req GetFileRequest) (*GetFileResponse, error) {
	var out GetFileResponse
	if err := c.caller.Call(ctx, string(ActionGetFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFilesetID 获取文件集 ID
func (c *Client) GetFilesetID(ctx context.Context, req GetFilesetIDRequest) (*GetFilesetIDResponse, error) {
	var out GetFilesetIDResponse
	if err := c.caller.Call(ctx, string(ActionGetFilesetID), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFilesetInfo 获取文件集信息
func (c *Client) GetFilesetInfo(ctx context.Context, req GetFilesetInfoRequest) (*GetFilesetInfoResponse, error) {
	var out GetFilesetInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetFilesetInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFlashFileList 获取闪传文件列表
func (c *Client) GetFlashFileList(ctx context.Context, req GetFlashFileListRequest) (*GetFlashFileListResponse, error) {
	var out GetFlashFileListResponse
	if err := c.caller.Call(ctx, string(ActionGetFlashFileList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFlashFileURL 获取闪传文件链接
func (c *Client) GetFlashFileURL(ctx context.Context, req GetFlashFileURLRequest) (*GetFlashFileURLResponse, error) {
	var out GetFlashFileURLResponse
	if err := c.caller.Call(ctx, string(ActionGetFlashFileURL), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetForwardMsg 获取合并转发消息
func (c *Client) GetForwardMsg(ctx context.Context, req GetForwardMsgRequest) (*GetForwardMsgResponse, error) {
	var out GetForwardMsgResponse
	if err := c.caller.Call(ctx, string(ActionGetForwardMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFriendList 获取好友列表
func (c *Client) GetFriendList(ctx context.Context, req GetFriendListRequest) (*GetFriendListResponse, error) {
	var out GetFriendListResponse
	if err := c.caller.Call(ctx, string(ActionGetFriendList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFriendMsgHistory 获取好友历史消息
func (c *Client) GetFriendMsgHistory(ctx context.Context, req GetFriendMsgHistoryRequest) (*GetFriendMsgHistoryResponse, error) {
	var out GetFriendMsgHistoryResponse
	if err := c.caller.Call(ctx, string(ActionGetFriendMsgHistory), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFriendsWithCategory 获取带分组的好友列表
func (c *Client) GetFriendsWithCategory(ctx context.Context, req GetFriendsWithCategoryRequest) (*GetFriendsWithCategoryResponse, error) {
	var out GetFriendsWithCategoryResponse
	if err := c.caller.Call(ctx, string(ActionGetFriendsWithCategory), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupAlbumMediaList 获取群相册媒体列表
func (c *Client) GetGroupAlbumMediaList(ctx context.Context, req GetGroupAlbumMediaListRequest) (*GetGroupAlbumMediaListResponse, error) {
	var out GetGroupAlbumMediaListResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupAlbumMediaList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupAtAllRemain 获取群艾特全体剩余次数
func (c *Client) GetGroupAtAllRemain(ctx context.Context, req GetGroupAtAllRemainRequest) (*GetGroupAtAllRemainResponse, error) {
	var out GetGroupAtAllRemainResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupAtAllRemain), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupDetailInfo 获取群详细信息
func (c *Client) GetGroupDetailInfo(ctx context.Context, req GetGroupDetailInfoRequest) (*GetGroupDetailInfoResponse, error) {
	var out GetGroupDetailInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupDetailInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupFileSystemInfo 获取群文件系统信息
func (c *Client) GetGroupFileSystemInfo(ctx context.Context, req GetGroupFileSystemInfoRequest) (*GetGroupFileSystemInfoResponse, error) {
	var out GetGroupFileSystemInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupFileSystemInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupFileURL 获取群文件URL
func (c *Client) GetGroupFileURL(ctx context.Context, req GetGroupFileURLRequest) (*GetGroupFileURLResponse, error) {
	var out GetGroupFileURLResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupFileURL), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupFilesByFolder 获取群文件夹文件列表
func (c *Client) GetGroupFilesByFolder(ctx context.Context, req GetGroupFilesByFolderRequest) (*GetGroupFilesByFolderResponse, error) {
	var out GetGroupFilesByFolderResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupFilesByFolder), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupHonorInfo 获取群荣誉信息
func (c *Client) GetGroupHonorInfo(ctx context.Context, req GetGroupHonorInfoRequest) (*GetGroupHonorInfoResponse, error) {
	var out GetGroupHonorInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupHonorInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupIgnoreAddRequest 获取群被忽略的加群请求
func (c *Client) GetGroupIgnoreAddRequest(ctx context.Context, req GetGroupIgnoreAddRequestRequest) (*GetGroupIgnoreAddRequestResponse, error) {
	var out GetGroupIgnoreAddRequestResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupIgnoreAddRequest), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupIgnoredNotifies 获取群忽略通知
func (c *Client) GetGroupIgnoredNotifies(ctx context.Context, req GetGroupIgnoredNotifiesRequest) (*GetGroupIgnoredNotifiesResponse, error) {
	var out GetGroupIgnoredNotifiesResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupIgnoredNotifies), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupInfo 获取群信息
func (c *Client) GetGroupInfo(ctx context.Context, req GetGroupInfoRequest) (*GetGroupInfoResponse, error) {
	var out GetGroupInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupInfoEx 获取群详细信息 (扩展)
func (c *Client) GetGroupInfoEx(ctx context.Context, req GetGroupInfoExRequest) (*GetGroupInfoExResponse, error) {
	var out GetGroupInfoExResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupInfoEx), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupList 获取群列表
func (c *Client) GetGroupList(ctx context.Context, req GetGroupListRequest) (*GetGroupListResponse, error) {
	var out GetGroupListResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupMemberInfo 获取群成员信息
func (c *Client) GetGroupMemberInfo(ctx context.Context, req GetGroupMemberInfoRequest) (*GetGroupMemberInfoResponse, error) {
	var out GetGroupMemberInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupMemberInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupMemberList 获取群成员列表
func (c *Client) GetGroupMemberList(ctx context.Context, req GetGroupMemberListRequest) (*GetGroupMemberListResponse, error) {
	var out GetGroupMemberListResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupMemberList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupMsgHistory 获取群历史消息
func (c *Client) GetGroupMsgHistory(ctx context.Context, req GetGroupMsgHistoryRequest) (*GetGroupMsgHistoryResponse, error) {
	var out GetGroupMsgHistoryResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupMsgHistory), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupRootFiles 获取群根目录文件列表
func (c *Client) GetGroupRootFiles(ctx context.Context, req GetGroupRootFilesRequest) (*GetGroupRootFilesResponse, error) {
	var out GetGroupRootFilesResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupRootFiles), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupShareLink 获取群分享链接
func (c *Client) GetGroupShareLink(ctx context.Context, req GetGroupShareLinkRequest) (*GetGroupShareLinkResponse, error) {
	var out GetGroupShareLinkResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupShareLink), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupShutList 获取群禁言列表
func (c *Client) GetGroupShutList(ctx context.Context, req GetGroupShutListRequest) (*GetGroupShutListResponse, error) {
	var out GetGroupShutListResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupShutList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupSignedList 获取群组今日打卡列表
func (c *Client) GetGroupSignedList(ctx context.Context, req GetGroupSignedListRequest) (*GetGroupSignedListResponse, error) {
	var out GetGroupSignedListResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupSignedList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupSystemMsg 获取群系统消息
func (c *Client) GetGroupSystemMsg(ctx context.Context, req GetGroupSystemMsgRequest) (*GetGroupSystemMsgResponse, error) {
	var out GetGroupSystemMsgResponse
	if err := c.caller.Call(ctx, string(ActionGetGroupSystemMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGuildList 获取频道列表
func (c *Client) GetGuildList(ctx context.Context, req GetGuildListRequest) (*GetGuildListResponse, error) {
	var out GetGuildListResponse
	if err := c.caller.Call(ctx, string(ActionGetGuildList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGuildServiceProfile 获取频道个人信息
func (c *Client) GetGuildServiceProfile(ctx context.Context, req GetGuildServiceProfileRequest) (*GetGuildServiceProfileResponse, error) {
	var out GetGuildServiceProfileResponse
	if err := c.caller.Call(ctx, string(ActionGetGuildServiceProfile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetImage 获取图片
func (c *Client) GetImage(ctx context.Context, req GetImageRequest) (*GetImageResponse, error) {
	var out GetImageResponse
	if err := c.caller.Call(ctx, string(ActionGetImage), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLoginInfo 获取登录号信息
func (c *Client) GetLoginInfo(ctx context.Context, req GetLoginInfoRequest) (*GetLoginInfoResponse, error) {
	var out GetLoginInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetLoginInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMiniAppArk 获取小程序 Ark
func (c *Client) GetMiniAppArk(ctx context.Context, req GetMiniAppArkRequestUnion) (*GetMiniAppArkResponse, error) {
	var out GetMiniAppArkResponse
	if err := c.caller.Call(ctx, string(ActionGetMiniAppArk), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMsg 获取消息
func (c *Client) GetMsg(ctx context.Context, req GetMsgRequest) (*GetMsgResponse, error) {
	var out GetMsgResponse
	if err := c.caller.Call(ctx, string(ActionGetMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOnlineClients 获取在线客户端
func (c *Client) GetOnlineClients(ctx context.Context, req GetOnlineClientsRequest) (*GetOnlineClientsResponse, error) {
	var out GetOnlineClientsResponse
	if err := c.caller.Call(ctx, string(ActionGetOnlineClients), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOnlineFileMsg 获取在线文件消息
func (c *Client) GetOnlineFileMsg(ctx context.Context, req GetOnlineFileMsgRequest) (*GetOnlineFileMsgResponse, error) {
	var out GetOnlineFileMsgResponse
	if err := c.caller.Call(ctx, string(ActionGetOnlineFileMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPrivateFileURL 获取私聊文件URL
func (c *Client) GetPrivateFileURL(ctx context.Context, req GetPrivateFileURLRequest) (*GetPrivateFileURLResponse, error) {
	var out GetPrivateFileURLResponse
	if err := c.caller.Call(ctx, string(ActionGetPrivateFileURL), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetProfileLike 获取资料点赞
func (c *Client) GetProfileLike(ctx context.Context, req GetProfileLikeRequest) (*GetProfileLikeResponse, error) {
	var out GetProfileLikeResponse
	if err := c.caller.Call(ctx, string(ActionGetProfileLike), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQunAlbumList 获取群相册列表
func (c *Client) GetQunAlbumList(ctx context.Context, req GetQunAlbumListRequest) (*GetQunAlbumListResponse, error) {
	var out GetQunAlbumListResponse
	if err := c.caller.Call(ctx, string(ActionGetQunAlbumList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRecentContact 获取最近会话
func (c *Client) GetRecentContact(ctx context.Context, req GetRecentContactRequest) (*GetRecentContactResponse, error) {
	var out GetRecentContactResponse
	if err := c.caller.Call(ctx, string(ActionGetRecentContact), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRecord 获取语音
func (c *Client) GetRecord(ctx context.Context, req GetRecordRequest) (*GetRecordResponse, error) {
	var out GetRecordResponse
	if err := c.caller.Call(ctx, string(ActionGetRecord), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRkey 获取扩展 RKey
func (c *Client) GetRkey(ctx context.Context, req GetRkeyRequest) (*GetRkeyResponse, error) {
	var out GetRkeyResponse
	if err := c.caller.Call(ctx, string(ActionGetRkey), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRkeyServer 获取 RKey 服务器
func (c *Client) GetRkeyServer(ctx context.Context, req GetRkeyServerRequest) (*GetRkeyServerResponse, error) {
	var out GetRkeyServerResponse
	if err := c.caller.Call(ctx, string(ActionGetRkeyServer), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRobotUINRange 获取机器人 UIN 范围
func (c *Client) GetRobotUINRange(ctx context.Context, req GetRobotUINRangeRequest) (*GetRobotUINRangeResponse, error) {
	var out GetRobotUINRangeResponse
	if err := c.caller.Call(ctx, string(ActionGetRobotUINRange), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShareLink 获取文件分享链接
func (c *Client) GetShareLink(ctx context.Context, req GetShareLinkRequest) (*GetShareLinkResponse, error) {
	var out GetShareLinkResponse
	if err := c.caller.Call(ctx, string(ActionGetShareLink), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetStatus 获取运行状态
func (c *Client) GetStatus(ctx context.Context, req GetStatusRequest) (*GetStatusResponse, error) {
	var out GetStatusResponse
	if err := c.caller.Call(ctx, string(ActionGetStatus), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetStrangerInfo 获取陌生人信息
func (c *Client) GetStrangerInfo(ctx context.Context, req GetStrangerInfoRequest) (*GetStrangerInfoResponse, error) {
	var out GetStrangerInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetStrangerInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUnidirectionalFriendList 获取单向好友列表
func (c *Client) GetUnidirectionalFriendList(ctx context.Context, req GetUnidirectionalFriendListRequest) (*GetUnidirectionalFriendListResponse, error) {
	var out GetUnidirectionalFriendListResponse
	if err := c.caller.Call(ctx, string(ActionGetUnidirectionalFriendList), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVersionInfo 获取版本信息
func (c *Client) GetVersionInfo(ctx context.Context, req GetVersionInfoRequest) (*GetVersionInfoResponse, error) {
	var out GetVersionInfoResponse
	if err := c.caller.Call(ctx, string(ActionGetVersionInfo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GroupPoke 发送戳一戳
func (c *Client) GroupPoke(ctx context.Context, req GroupPokeRequest) (*GroupPokeResponse, error) {
	var out GroupPokeResponse
	if err := c.caller.Call(ctx, string(ActionGroupPoke), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MarkGroupMsgAsRead 标记群聊已读
func (c *Client) MarkGroupMsgAsRead(ctx context.Context, req MarkGroupMsgAsReadRequest) (*MarkGroupMsgAsReadResponse, error) {
	var out MarkGroupMsgAsReadResponse
	if err := c.caller.Call(ctx, string(ActionMarkGroupMsgAsRead), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MarkMsgAsRead 标记消息已读 (Go-CQHTTP)
func (c *Client) MarkMsgAsRead(ctx context.Context, req MarkMsgAsReadRequest) (*MarkMsgAsReadResponse, error) {
	var out MarkMsgAsReadResponse
	if err := c.caller.Call(ctx, string(ActionMarkMsgAsRead), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MarkPrivateMsgAsRead 标记私聊已读
func (c *Client) MarkPrivateMsgAsRead(ctx context.Context, req MarkPrivateMsgAsReadRequest) (*MarkPrivateMsgAsReadResponse, error) {
	var out MarkPrivateMsgAsReadResponse
	if err := c.caller.Call(ctx, string(ActionMarkPrivateMsgAsRead), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MoveGroupFile 移动群文件
func (c *Client) MoveGroupFile(ctx context.Context, req MoveGroupFileRequest) (*MoveGroupFileResponse, error) {
	var out MoveGroupFileResponse
	if err := c.caller.Call(ctx, string(ActionMoveGroupFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// NcGetPacketStatus 获取Packet状态
func (c *Client) NcGetPacketStatus(ctx context.Context, req NcGetPacketStatusRequest) (*NcGetPacketStatusResponse, error) {
	var out NcGetPacketStatusResponse
	if err := c.caller.Call(ctx, string(ActionNcGetPacketStatus), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// NcGetRkey 获取 RKey
func (c *Client) NcGetRkey(ctx context.Context, req NcGetRkeyRequest) (*NcGetRkeyResponse, error) {
	var out NcGetRkeyResponse
	if err := c.caller.Call(ctx, string(ActionNcGetRkey), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// NcGetUserStatus 获取用户在线状态
func (c *Client) NcGetUserStatus(ctx context.Context, req NcGetUserStatusRequest) (*NcGetUserStatusResponse, error) {
	var out NcGetUserStatusResponse
	if err := c.caller.Call(ctx, string(ActionNcGetUserStatus), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// OcrImage 图片 OCR 识别
func (c *Client) OcrImage(ctx context.Context, req OcrImageRequest) (*OcrImageResponse, error) {
	var out OcrImageResponse
	if err := c.caller.Call(ctx, string(ActionOcrImage), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReceiveOnlineFile 接收在线文件
func (c *Client) ReceiveOnlineFile(ctx context.Context, req ReceiveOnlineFileRequest) (*ReceiveOnlineFileResponse, error) {
	var out ReceiveOnlineFileResponse
	if err := c.caller.Call(ctx, string(ActionReceiveOnlineFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RefuseOnlineFile 拒绝在线文件
func (c *Client) RefuseOnlineFile(ctx context.Context, req RefuseOnlineFileRequest) (*RefuseOnlineFileResponse, error) {
	var out RefuseOnlineFileResponse
	if err := c.caller.Call(ctx, string(ActionRefuseOnlineFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RenameGroupFile 重命名群文件
func (c *Client) RenameGroupFile(ctx context.Context, req RenameGroupFileRequest) (*RenameGroupFileResponse, error) {
	var out RenameGroupFileResponse
	if err := c.caller.Call(ctx, string(ActionRenameGroupFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendArkShare 分享用户 (Ark)
func (c *Client) SendArkShare(ctx context.Context, req SendArkShareRequest) (*SendArkShareResponse, error) {
	var out SendArkShareResponse
	if err := c.caller.Call(ctx, string(ActionSendArkShare), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendFlashMsg 发送闪传消息
func (c *Client) SendFlashMsg(ctx context.Context, req SendFlashMsgRequest) (*SendFlashMsgResponse, error) {
	var out SendFlashMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendFlashMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendForwardMsg 发送合并转发消息
func (c *Client) SendForwardMsg(ctx context.Context, req SendForwardMsgRequest) (*SendForwardMsgResponse, error) {
	var out SendForwardMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendForwardMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendGroupAiRecord 发送群 AI 语音
func (c *Client) SendGroupAiRecord(ctx context.Context, req SendGroupAiRecordRequest) (*SendGroupAiRecordResponse, error) {
	var out SendGroupAiRecordResponse
	if err := c.caller.Call(ctx, string(ActionSendGroupAiRecord), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendGroupArkShare 分享群 (Ark)
func (c *Client) SendGroupArkShare(ctx context.Context, req SendGroupArkShareRequest) (*SendGroupArkShareResponse, error) {
	var out SendGroupArkShareResponse
	if err := c.caller.Call(ctx, string(ActionSendGroupArkShare), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendGroupForwardMsg 发送群合并转发消息
func (c *Client) SendGroupForwardMsg(ctx context.Context, req SendGroupForwardMsgRequest) (*SendGroupForwardMsgResponse, error) {
	var out SendGroupForwardMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendGroupForwardMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendGroupMsg 发送群消息
func (c *Client) SendGroupMsg(ctx context.Context, req SendGroupMsgRequest) (*SendGroupMsgResponse, error) {
	var out SendGroupMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendGroupMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendGroupSign 群打卡
func (c *Client) SendGroupSign(ctx context.Context, req SendGroupSignRequest) (*SendGroupSignResponse, error) {
	var out SendGroupSignResponse
	if err := c.caller.Call(ctx, string(ActionSendGroupSign), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendLike 点赞
func (c *Client) SendLike(ctx context.Context, req SendLikeRequest) (*SendLikeResponse, error) {
	var out SendLikeResponse
	if err := c.caller.Call(ctx, string(ActionSendLike), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendMsg 发送消息
func (c *Client) SendMsg(ctx context.Context, req SendMsgRequest) (*SendMsgResponse, error) {
	var out SendMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendOnlineFile 发送在线文件
func (c *Client) SendOnlineFile(ctx context.Context, req SendOnlineFileRequest) (*SendOnlineFileResponse, error) {
	var out SendOnlineFileResponse
	if err := c.caller.Call(ctx, string(ActionSendOnlineFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendOnlineFolder 发送在线文件夹
func (c *Client) SendOnlineFolder(ctx context.Context, req SendOnlineFolderRequest) (*SendOnlineFolderResponse, error) {
	var out SendOnlineFolderResponse
	if err := c.caller.Call(ctx, string(ActionSendOnlineFolder), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendPacket 发送原始数据包
func (c *Client) SendPacket(ctx context.Context, req SendPacketRequest) (*SendPacketResponseVariant1, error) {
	var out SendPacketResponseVariant1
	if err := c.caller.Call(ctx, string(ActionSendPacket), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendPoke 发送戳一戳
func (c *Client) SendPoke(ctx context.Context, req SendPokeRequest) (*SendPokeResponse, error) {
	var out SendPokeResponse
	if err := c.caller.Call(ctx, string(ActionSendPoke), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendPrivateForwardMsg 发送私聊合并转发消息
func (c *Client) SendPrivateForwardMsg(ctx context.Context, req SendPrivateForwardMsgRequest) (*SendPrivateForwardMsgResponse, error) {
	var out SendPrivateForwardMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendPrivateForwardMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendPrivateMsg 发送私聊消息
func (c *Client) SendPrivateMsg(ctx context.Context, req SendPrivateMsgRequest) (*SendPrivateMsgResponse, error) {
	var out SendPrivateMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendPrivateMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendQzoneMsg 发表QQ空间说说
func (c *Client) SendQzoneMsg(ctx context.Context, req SendQzoneMsgRequest) (*SendQzoneMsgResponse, error) {
	var out SendQzoneMsgResponse
	if err := c.caller.Call(ctx, string(ActionSendQzoneMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetCustomFaceDesc 修改自定义表情描述
func (c *Client) SetCustomFaceDesc(ctx context.Context, req SetCustomFaceDescRequest) (*SetCustomFaceDescResponse, error) {
	var out SetCustomFaceDescResponse
	if err := c.caller.Call(ctx, string(ActionSetCustomFaceDesc), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetDiyOnlineStatus 设置自定义在线状态
func (c *Client) SetDiyOnlineStatus(ctx context.Context, req SetDiyOnlineStatusRequest) (*SetDiyOnlineStatusResponse, error) {
	var out SetDiyOnlineStatusResponse
	if err := c.caller.Call(ctx, string(ActionSetDiyOnlineStatus), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetDoubtFriendsAddRequest 处理可疑好友申请
func (c *Client) SetDoubtFriendsAddRequest(ctx context.Context, req SetDoubtFriendsAddRequestRequest) (*SetDoubtFriendsAddRequestResponse, error) {
	var out SetDoubtFriendsAddRequestResponse
	if err := c.caller.Call(ctx, string(ActionSetDoubtFriendsAddRequest), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetEssenceMsg 设置精华消息
func (c *Client) SetEssenceMsg(ctx context.Context, req SetEssenceMsgRequest) (*SetEssenceMsgResponse, error) {
	var out SetEssenceMsgResponse
	if err := c.caller.Call(ctx, string(ActionSetEssenceMsg), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetFriendAddRequest 处理加好友请求
func (c *Client) SetFriendAddRequest(ctx context.Context, req SetFriendAddRequestRequest) (*SetFriendAddRequestResponse, error) {
	var out SetFriendAddRequestResponse
	if err := c.caller.Call(ctx, string(ActionSetFriendAddRequest), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetFriendRemark 设置好友备注
func (c *Client) SetFriendRemark(ctx context.Context, req SetFriendRemarkRequest) (*SetFriendRemarkResponse, error) {
	var out SetFriendRemarkResponse
	if err := c.caller.Call(ctx, string(ActionSetFriendRemark), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupAddOption 设置群加群选项
func (c *Client) SetGroupAddOption(ctx context.Context, req SetGroupAddOptionRequest) (*SetGroupAddOptionResponse, error) {
	var out SetGroupAddOptionResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupAddOption), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupAddRequest 处理加群请求
func (c *Client) SetGroupAddRequest(ctx context.Context, req SetGroupAddRequestRequest) (*SetGroupAddRequestResponse, error) {
	var out SetGroupAddRequestResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupAddRequest), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupAdmin 设置群管理员
func (c *Client) SetGroupAdmin(ctx context.Context, req SetGroupAdminRequest) (*SetGroupAdminResponse, error) {
	var out SetGroupAdminResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupAdmin), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupAlbumMediaLike 点赞群相册媒体
func (c *Client) SetGroupAlbumMediaLike(ctx context.Context, req SetGroupAlbumMediaLikeRequest) (*SetGroupAlbumMediaLikeResponse, error) {
	var out SetGroupAlbumMediaLikeResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupAlbumMediaLike), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupBan 群组禁言
func (c *Client) SetGroupBan(ctx context.Context, req SetGroupBanRequest) (*SetGroupBanResponse, error) {
	var out SetGroupBanResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupBan), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupCard 设置群名片
func (c *Client) SetGroupCard(ctx context.Context, req SetGroupCardRequest) (*SetGroupCardResponse, error) {
	var out SetGroupCardResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupCard), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupKick 群组踢人
func (c *Client) SetGroupKick(ctx context.Context, req SetGroupKickRequest) (*SetGroupKickResponse, error) {
	var out SetGroupKickResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupKick), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupKickMembers 批量踢出群成员
func (c *Client) SetGroupKickMembers(ctx context.Context, req SetGroupKickMembersRequest) (*SetGroupKickMembersResponse, error) {
	var out SetGroupKickMembersResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupKickMembers), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupLeave 退出群组
func (c *Client) SetGroupLeave(ctx context.Context, req SetGroupLeaveRequest) (*SetGroupLeaveResponse, error) {
	var out SetGroupLeaveResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupLeave), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupMemberInvitePolicy 设置群成员邀请策略
func (c *Client) SetGroupMemberInvitePolicy(ctx context.Context, req SetGroupMemberInvitePolicyRequest) (*SetGroupMemberInvitePolicyResponse, error) {
	var out SetGroupMemberInvitePolicyResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupMemberInvitePolicy), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupMemberPermissions 设置群成员功能权限
func (c *Client) SetGroupMemberPermissions(ctx context.Context, req SetGroupMemberPermissionsRequest) (*SetGroupMemberPermissionsResponse, error) {
	var out SetGroupMemberPermissionsResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupMemberPermissions), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupName 设置群名称
func (c *Client) SetGroupName(ctx context.Context, req SetGroupNameRequest) (*SetGroupNameResponse, error) {
	var out SetGroupNameResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupName), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupNewMemberHistoryVisibility 设置新成员历史消息可见性
func (c *Client) SetGroupNewMemberHistoryVisibility(ctx context.Context, req SetGroupNewMemberHistoryVisibilityRequest) (*SetGroupNewMemberHistoryVisibilityResponse, error) {
	var out SetGroupNewMemberHistoryVisibilityResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupNewMemberHistoryVisibility), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupPortrait 设置群头像
func (c *Client) SetGroupPortrait(ctx context.Context, req SetGroupPortraitRequest) (*SetGroupPortraitResponse, error) {
	var out SetGroupPortraitResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupPortrait), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupRemark 设置群备注
func (c *Client) SetGroupRemark(ctx context.Context, req SetGroupRemarkRequest) (*SetGroupRemarkResponse, error) {
	var out SetGroupRemarkResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupRemark), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupRobotAddOption 设置群机器人加群选项
func (c *Client) SetGroupRobotAddOption(ctx context.Context, req SetGroupRobotAddOptionRequest) (*SetGroupRobotAddOptionResponse, error) {
	var out SetGroupRobotAddOptionResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupRobotAddOption), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupSearch 设置群搜索选项
func (c *Client) SetGroupSearch(ctx context.Context, req SetGroupSearchRequest) (*SetGroupSearchResponse, error) {
	var out SetGroupSearchResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupSearch), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupSign 群打卡
func (c *Client) SetGroupSign(ctx context.Context, req SetGroupSignRequest) (*SetGroupSignResponse, error) {
	var out SetGroupSignResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupSign), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupSpecialTitle 设置专属头衔
func (c *Client) SetGroupSpecialTitle(ctx context.Context, req SetGroupSpecialTitleRequest) (*SetGroupSpecialTitleResponse, error) {
	var out SetGroupSpecialTitleResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupSpecialTitle), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupTodo 设置群待办
func (c *Client) SetGroupTodo(ctx context.Context, req SetGroupTodoRequest) (*SetGroupTodoResponse, error) {
	var out SetGroupTodoResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupTodo), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetGroupWholeBan 全员禁言
func (c *Client) SetGroupWholeBan(ctx context.Context, req SetGroupWholeBanRequest) (*SetGroupWholeBanResponse, error) {
	var out SetGroupWholeBanResponse
	if err := c.caller.Call(ctx, string(ActionSetGroupWholeBan), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetInputStatus 设置输入状态
func (c *Client) SetInputStatus(ctx context.Context, req SetInputStatusRequest) (*SetInputStatusResponse, error) {
	var out SetInputStatusResponse
	if err := c.caller.Call(ctx, string(ActionSetInputStatus), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetMsgEmojiLike 设置消息表情点赞
func (c *Client) SetMsgEmojiLike(ctx context.Context, req SetMsgEmojiLikeRequest) (*SetMsgEmojiLikeResponse, error) {
	var out SetMsgEmojiLikeResponse
	if err := c.caller.Call(ctx, string(ActionSetMsgEmojiLike), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetOnlineStatus 设置在线状态
func (c *Client) SetOnlineStatus(ctx context.Context, req SetOnlineStatusRequest) (*SetOnlineStatusResponse, error) {
	var out SetOnlineStatusResponse
	if err := c.caller.Call(ctx, string(ActionSetOnlineStatus), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetQQAvatar 设置QQ头像
func (c *Client) SetQQAvatar(ctx context.Context, req SetQQAvatarRequest) (*SetQQAvatarResponse, error) {
	var out SetQQAvatarResponse
	if err := c.caller.Call(ctx, string(ActionSetQQAvatar), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetQQProfile 设置QQ资料
func (c *Client) SetQQProfile(ctx context.Context, req SetQQProfileRequest) (*SetQQProfileResponse, error) {
	var out SetQQProfileResponse
	if err := c.caller.Call(ctx, string(ActionSetQQProfile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetRestart 重启服务
func (c *Client) SetRestart(ctx context.Context, req SetRestartRequest) (*SetRestartResponse, error) {
	var out SetRestartResponse
	if err := c.caller.Call(ctx, string(ActionSetRestart), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetSelfLongnick 设置个性签名
func (c *Client) SetSelfLongnick(ctx context.Context, req SetSelfLongnickRequest) (*SetSelfLongnickResponse, error) {
	var out SetSelfLongnickResponse
	if err := c.caller.Call(ctx, string(ActionSetSelfLongnick), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TestDownloadStream 测试下载流
func (c *Client) TestDownloadStream(ctx context.Context, req TestDownloadStreamRequest) (*TestDownloadStreamResponse, error) {
	var out TestDownloadStreamResponse
	if err := c.caller.Call(ctx, string(ActionTestDownloadStream), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TransGroupFile 传输群文件
func (c *Client) TransGroupFile(ctx context.Context, req TransGroupFileRequest) (*TransGroupFileResponse, error) {
	var out TransGroupFileResponse
	if err := c.caller.Call(ctx, string(ActionTransGroupFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TranslateEn2zh 英文单词翻译
func (c *Client) TranslateEn2zh(ctx context.Context, req TranslateEn2zhRequest) (*TranslateEn2zhResponse, error) {
	var out TranslateEn2zhResponse
	if err := c.caller.Call(ctx, string(ActionTranslateEn2zh), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadFileStream 上传文件流
func (c *Client) UploadFileStream(ctx context.Context, req UploadFileStreamRequest) (*UploadFileStreamResponse, error) {
	var out UploadFileStreamResponse
	if err := c.caller.Call(ctx, string(ActionUploadFileStream), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadGroupFile 上传群文件
func (c *Client) UploadGroupFile(ctx context.Context, req UploadGroupFileRequest) (*UploadGroupFileResponse, error) {
	var out UploadGroupFileResponse
	if err := c.caller.Call(ctx, string(ActionUploadGroupFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadImageToQunAlbum 上传图片到群相册
func (c *Client) UploadImageToQunAlbum(ctx context.Context, req UploadImageToQunAlbumRequest) (*UploadImageToQunAlbumResponse, error) {
	var out UploadImageToQunAlbumResponse
	if err := c.caller.Call(ctx, string(ActionUploadImageToQunAlbum), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadImagesToQunAlbum 上传多张图片到群相册
func (c *Client) UploadImagesToQunAlbum(ctx context.Context, req UploadImagesToQunAlbumRequest) (*UploadImagesToQunAlbumResponse, error) {
	var out UploadImagesToQunAlbumResponse
	if err := c.caller.Call(ctx, string(ActionUploadImagesToQunAlbum), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadPrivateFile 上传私聊文件
func (c *Client) UploadPrivateFile(ctx context.Context, req UploadPrivateFileRequest) (*UploadPrivateFileResponse, error) {
	var out UploadPrivateFileResponse
	if err := c.caller.Call(ctx, string(ActionUploadPrivateFile), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadVideoToQunAlbum 上传视频到群相册
func (c *Client) UploadVideoToQunAlbum(ctx context.Context, req UploadVideoToQunAlbumRequest) (*UploadVideoToQunAlbumResponse, error) {
	var out UploadVideoToQunAlbumResponse
	if err := c.caller.Call(ctx, string(ActionUploadVideoToQunAlbum), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

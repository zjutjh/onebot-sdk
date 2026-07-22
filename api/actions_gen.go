// 代码由 napcatgen 根据 NapCat OpenAPI 4.18.13 生成；请勿手动修改。

package api

// Action 是 NapCat API action 名称。
type Action string

const (
	// ActionDotHandleQuickOperation 处理快速操作
	ActionDotHandleQuickOperation Action = ".handle_quick_operation"
	// ActionDotOcrImage 图片 OCR 识别 (内部)
	ActionDotOcrImage Action = ".ocr_image"
	// ActionArkShareGroup 分享群 (Ark)
	ActionArkShareGroup Action = "ArkShareGroup"
	// ActionArkSharePeer 分享用户 (Ark)
	ActionArkSharePeer Action = "ArkSharePeer"
	// ActionUnderscoreDelGroupNotice 删除群公告
	ActionUnderscoreDelGroupNotice Action = "_del_group_notice"
	// ActionUnderscoreGetGroupNotice 获取群公告
	ActionUnderscoreGetGroupNotice Action = "_get_group_notice"
	// ActionUnderscoreGetModelShow 获取机型显示
	ActionUnderscoreGetModelShow Action = "_get_model_show"
	// ActionUnderscoreMarkAllAsRead 标记所有消息已读
	ActionUnderscoreMarkAllAsRead Action = "_mark_all_as_read"
	// ActionUnderscoreSendGroupNotice 发送群公告
	ActionUnderscoreSendGroupNotice Action = "_send_group_notice"
	// ActionUnderscoreSetModelShow 设置机型
	ActionUnderscoreSetModelShow Action = "_set_model_show"
	// ActionAddCustomFace 添加自定义表情
	ActionAddCustomFace Action = "add_custom_face"
	// ActionBotExit 退出登录
	ActionBotExit Action = "bot_exit"
	// ActionCanSendImage 是否可以发送图片
	ActionCanSendImage Action = "can_send_image"
	// ActionCanSendRecord 是否可以发送语音
	ActionCanSendRecord Action = "can_send_record"
	// ActionCancelGroupAlbumMediaLike 取消点赞群相册媒体
	ActionCancelGroupAlbumMediaLike Action = "cancel_group_album_media_like"
	// ActionCancelGroupTodo 取消群待办
	ActionCancelGroupTodo Action = "cancel_group_todo"
	// ActionCancelOnlineFile 取消在线文件
	ActionCancelOnlineFile Action = "cancel_online_file"
	// ActionCheckURLSafely 检查URL安全性
	ActionCheckURLSafely Action = "check_url_safely"
	// ActionCleanCache 清理缓存
	ActionCleanCache Action = "clean_cache"
	// ActionCleanStreamTempFile 清理流式传输临时文件
	ActionCleanStreamTempFile Action = "clean_stream_temp_file"
	// ActionClickInlineKeyboardButton 点击内联键盘按钮
	ActionClickInlineKeyboardButton Action = "click_inline_keyboard_button"
	// ActionCompleteGroupTodo 完成群待办
	ActionCompleteGroupTodo Action = "complete_group_todo"
	// ActionCreateCollection 创建收藏
	ActionCreateCollection Action = "create_collection"
	// ActionCreateFlashTask 创建闪传任务
	ActionCreateFlashTask Action = "create_flash_task"
	// ActionCreateGroupFileFolder 创建群文件目录
	ActionCreateGroupFileFolder Action = "create_group_file_folder"
	// ActionDelGroupAlbumMedia 删除群相册媒体
	ActionDelGroupAlbumMedia Action = "del_group_album_media"
	// ActionDeleteCustomFace 删除自定义表情
	ActionDeleteCustomFace Action = "delete_custom_face"
	// ActionDeleteEssenceMsg 移出精华消息
	ActionDeleteEssenceMsg Action = "delete_essence_msg"
	// ActionDeleteFriend 删除好友
	ActionDeleteFriend Action = "delete_friend"
	// ActionDeleteGroupFile 删除群文件
	ActionDeleteGroupFile Action = "delete_group_file"
	// ActionDeleteGroupFolder 删除群文件目录
	ActionDeleteGroupFolder Action = "delete_group_folder"
	// ActionDeleteMsg 撤回消息
	ActionDeleteMsg Action = "delete_msg"
	// ActionDoGroupAlbumComment 发表群相册评论
	ActionDoGroupAlbumComment Action = "do_group_album_comment"
	// ActionDownloadFile 下载文件
	ActionDownloadFile Action = "download_file"
	// ActionDownloadFileImageStream 下载图片文件流
	ActionDownloadFileImageStream Action = "download_file_image_stream"
	// ActionDownloadFileRecordStream 下载语音文件流
	ActionDownloadFileRecordStream Action = "download_file_record_stream"
	// ActionDownloadFileStream 下载文件流
	ActionDownloadFileStream Action = "download_file_stream"
	// ActionDownloadFileset 下载文件集
	ActionDownloadFileset Action = "download_fileset"
	// ActionFetchCustomFace 获取自定义表情
	ActionFetchCustomFace Action = "fetch_custom_face"
	// ActionFetchCustomFaceDetail 获取自定义表情详情
	ActionFetchCustomFaceDetail Action = "fetch_custom_face_detail"
	// ActionFetchEmojiLike 获取表情点赞详情
	ActionFetchEmojiLike Action = "fetch_emoji_like"
	// ActionFetchPttText 获取语音转文字结果
	ActionFetchPttText Action = "fetch_ptt_text"
	// ActionForwardFriendSingleMsg 转发单条消息
	ActionForwardFriendSingleMsg Action = "forward_friend_single_msg"
	// ActionForwardGroupSingleMsg 转发单条消息
	ActionForwardGroupSingleMsg Action = "forward_group_single_msg"
	// ActionFriendPoke 发送戳一戳
	ActionFriendPoke Action = "friend_poke"
	// ActionGetAiCharacters 获取AI角色列表
	ActionGetAiCharacters Action = "get_ai_characters"
	// ActionGetAiRecord 获取 AI 语音
	ActionGetAiRecord Action = "get_ai_record"
	// ActionGetClientkey 获取ClientKey
	ActionGetClientkey Action = "get_clientkey"
	// ActionGetCollectionList 获取收藏列表
	ActionGetCollectionList Action = "get_collection_list"
	// ActionGetCookies 获取 Cookies
	ActionGetCookies Action = "get_cookies"
	// ActionGetCredentials 获取登录凭证
	ActionGetCredentials Action = "get_credentials"
	// ActionGetCsrfToken 获取 CSRF Token
	ActionGetCsrfToken Action = "get_csrf_token"
	// ActionGetDoubtFriendsAddRequest 获取可疑好友申请
	ActionGetDoubtFriendsAddRequest Action = "get_doubt_friends_add_request"
	// ActionGetEmojiLikes 获取消息表情点赞列表
	ActionGetEmojiLikes Action = "get_emoji_likes"
	// ActionGetEssenceMsgList 获取群精华消息
	ActionGetEssenceMsgList Action = "get_essence_msg_list"
	// ActionGetFile 获取文件
	ActionGetFile Action = "get_file"
	// ActionGetFilesetID 获取文件集 ID
	ActionGetFilesetID Action = "get_fileset_id"
	// ActionGetFilesetInfo 获取文件集信息
	ActionGetFilesetInfo Action = "get_fileset_info"
	// ActionGetFlashFileList 获取闪传文件列表
	ActionGetFlashFileList Action = "get_flash_file_list"
	// ActionGetFlashFileURL 获取闪传文件链接
	ActionGetFlashFileURL Action = "get_flash_file_url"
	// ActionGetForwardMsg 获取合并转发消息
	ActionGetForwardMsg Action = "get_forward_msg"
	// ActionGetFriendList 获取好友列表
	ActionGetFriendList Action = "get_friend_list"
	// ActionGetFriendMsgHistory 获取好友历史消息
	ActionGetFriendMsgHistory Action = "get_friend_msg_history"
	// ActionGetFriendsWithCategory 获取带分组的好友列表
	ActionGetFriendsWithCategory Action = "get_friends_with_category"
	// ActionGetGroupAlbumMediaList 获取群相册媒体列表
	ActionGetGroupAlbumMediaList Action = "get_group_album_media_list"
	// ActionGetGroupAtAllRemain 获取群艾特全体剩余次数
	ActionGetGroupAtAllRemain Action = "get_group_at_all_remain"
	// ActionGetGroupDetailInfo 获取群详细信息
	ActionGetGroupDetailInfo Action = "get_group_detail_info"
	// ActionGetGroupFileSystemInfo 获取群文件系统信息
	ActionGetGroupFileSystemInfo Action = "get_group_file_system_info"
	// ActionGetGroupFileURL 获取群文件URL
	ActionGetGroupFileURL Action = "get_group_file_url"
	// ActionGetGroupFilesByFolder 获取群文件夹文件列表
	ActionGetGroupFilesByFolder Action = "get_group_files_by_folder"
	// ActionGetGroupHonorInfo 获取群荣誉信息
	ActionGetGroupHonorInfo Action = "get_group_honor_info"
	// ActionGetGroupIgnoreAddRequest 获取群被忽略的加群请求
	ActionGetGroupIgnoreAddRequest Action = "get_group_ignore_add_request"
	// ActionGetGroupIgnoredNotifies 获取群忽略通知
	ActionGetGroupIgnoredNotifies Action = "get_group_ignored_notifies"
	// ActionGetGroupInfo 获取群信息
	ActionGetGroupInfo Action = "get_group_info"
	// ActionGetGroupInfoEx 获取群详细信息 (扩展)
	ActionGetGroupInfoEx Action = "get_group_info_ex"
	// ActionGetGroupList 获取群列表
	ActionGetGroupList Action = "get_group_list"
	// ActionGetGroupMemberInfo 获取群成员信息
	ActionGetGroupMemberInfo Action = "get_group_member_info"
	// ActionGetGroupMemberList 获取群成员列表
	ActionGetGroupMemberList Action = "get_group_member_list"
	// ActionGetGroupMsgHistory 获取群历史消息
	ActionGetGroupMsgHistory Action = "get_group_msg_history"
	// ActionGetGroupRootFiles 获取群根目录文件列表
	ActionGetGroupRootFiles Action = "get_group_root_files"
	// ActionGetGroupShutList 获取群禁言列表
	ActionGetGroupShutList Action = "get_group_shut_list"
	// ActionGetGroupSignedList 获取群组今日打卡列表
	ActionGetGroupSignedList Action = "get_group_signed_list"
	// ActionGetGroupSystemMsg 获取群系统消息
	ActionGetGroupSystemMsg Action = "get_group_system_msg"
	// ActionGetGuildList 获取频道列表
	ActionGetGuildList Action = "get_guild_list"
	// ActionGetGuildServiceProfile 获取频道个人信息
	ActionGetGuildServiceProfile Action = "get_guild_service_profile"
	// ActionGetImage 获取图片
	ActionGetImage Action = "get_image"
	// ActionGetLoginInfo 获取登录号信息
	ActionGetLoginInfo Action = "get_login_info"
	// ActionGetMiniAppArk 获取小程序 Ark
	ActionGetMiniAppArk Action = "get_mini_app_ark"
	// ActionGetMsg 获取消息
	ActionGetMsg Action = "get_msg"
	// ActionGetOnlineClients 获取在线客户端
	ActionGetOnlineClients Action = "get_online_clients"
	// ActionGetOnlineFileMsg 获取在线文件消息
	ActionGetOnlineFileMsg Action = "get_online_file_msg"
	// ActionGetPrivateFileURL 获取私聊文件URL
	ActionGetPrivateFileURL Action = "get_private_file_url"
	// ActionGetProfileLike 获取资料点赞
	ActionGetProfileLike Action = "get_profile_like"
	// ActionGetQunAlbumList 获取群相册列表
	ActionGetQunAlbumList Action = "get_qun_album_list"
	// ActionGetRecentContact 获取最近会话
	ActionGetRecentContact Action = "get_recent_contact"
	// ActionGetRecord 获取语音
	ActionGetRecord Action = "get_record"
	// ActionGetRkey 获取扩展 RKey
	ActionGetRkey Action = "get_rkey"
	// ActionGetRkeyServer 获取 RKey 服务器
	ActionGetRkeyServer Action = "get_rkey_server"
	// ActionGetRobotUINRange 获取机器人 UIN 范围
	ActionGetRobotUINRange Action = "get_robot_uin_range"
	// ActionGetShareLink 获取文件分享链接
	ActionGetShareLink Action = "get_share_link"
	// ActionGetStatus 获取运行状态
	ActionGetStatus Action = "get_status"
	// ActionGetStrangerInfo 获取陌生人信息
	ActionGetStrangerInfo Action = "get_stranger_info"
	// ActionGetUnidirectionalFriendList 获取单向好友列表
	ActionGetUnidirectionalFriendList Action = "get_unidirectional_friend_list"
	// ActionGetVersionInfo 获取版本信息
	ActionGetVersionInfo Action = "get_version_info"
	// ActionGroupPoke 发送戳一戳
	ActionGroupPoke Action = "group_poke"
	// ActionMarkGroupMsgAsRead 标记群聊已读
	ActionMarkGroupMsgAsRead Action = "mark_group_msg_as_read"
	// ActionMarkMsgAsRead 标记消息已读 (Go-CQHTTP)
	ActionMarkMsgAsRead Action = "mark_msg_as_read"
	// ActionMarkPrivateMsgAsRead 标记私聊已读
	ActionMarkPrivateMsgAsRead Action = "mark_private_msg_as_read"
	// ActionMoveGroupFile 移动群文件
	ActionMoveGroupFile Action = "move_group_file"
	// ActionNcGetPacketStatus 获取Packet状态
	ActionNcGetPacketStatus Action = "nc_get_packet_status"
	// ActionNcGetRkey 获取 RKey
	ActionNcGetRkey Action = "nc_get_rkey"
	// ActionNcGetUserStatus 获取用户在线状态
	ActionNcGetUserStatus Action = "nc_get_user_status"
	// ActionOcrImage 图片 OCR 识别
	ActionOcrImage Action = "ocr_image"
	// ActionReceiveOnlineFile 接收在线文件
	ActionReceiveOnlineFile Action = "receive_online_file"
	// ActionRefuseOnlineFile 拒绝在线文件
	ActionRefuseOnlineFile Action = "refuse_online_file"
	// ActionRenameGroupFile 重命名群文件
	ActionRenameGroupFile Action = "rename_group_file"
	// ActionSendArkShare 分享用户 (Ark)
	ActionSendArkShare Action = "send_ark_share"
	// ActionSendFlashMsg 发送闪传消息
	ActionSendFlashMsg Action = "send_flash_msg"
	// ActionSendForwardMsg 发送合并转发消息
	ActionSendForwardMsg Action = "send_forward_msg"
	// ActionSendGroupAiRecord 发送群 AI 语音
	ActionSendGroupAiRecord Action = "send_group_ai_record"
	// ActionSendGroupArkShare 分享群 (Ark)
	ActionSendGroupArkShare Action = "send_group_ark_share"
	// ActionSendGroupForwardMsg 发送群合并转发消息
	ActionSendGroupForwardMsg Action = "send_group_forward_msg"
	// ActionSendGroupMsg 发送群消息
	ActionSendGroupMsg Action = "send_group_msg"
	// ActionSendGroupSign 群打卡
	ActionSendGroupSign Action = "send_group_sign"
	// ActionSendLike 点赞
	ActionSendLike Action = "send_like"
	// ActionSendMsg 发送消息
	ActionSendMsg Action = "send_msg"
	// ActionSendOnlineFile 发送在线文件
	ActionSendOnlineFile Action = "send_online_file"
	// ActionSendOnlineFolder 发送在线文件夹
	ActionSendOnlineFolder Action = "send_online_folder"
	// ActionSendPacket 发送原始数据包
	ActionSendPacket Action = "send_packet"
	// ActionSendPoke 发送戳一戳
	ActionSendPoke Action = "send_poke"
	// ActionSendPrivateForwardMsg 发送私聊合并转发消息
	ActionSendPrivateForwardMsg Action = "send_private_forward_msg"
	// ActionSendPrivateMsg 发送私聊消息
	ActionSendPrivateMsg Action = "send_private_msg"
	// ActionSetCustomFaceDesc 修改自定义表情描述
	ActionSetCustomFaceDesc Action = "set_custom_face_desc"
	// ActionSetDiyOnlineStatus 设置自定义在线状态
	ActionSetDiyOnlineStatus Action = "set_diy_online_status"
	// ActionSetDoubtFriendsAddRequest 处理可疑好友申请
	ActionSetDoubtFriendsAddRequest Action = "set_doubt_friends_add_request"
	// ActionSetEssenceMsg 设置精华消息
	ActionSetEssenceMsg Action = "set_essence_msg"
	// ActionSetFriendAddRequest 处理加好友请求
	ActionSetFriendAddRequest Action = "set_friend_add_request"
	// ActionSetFriendRemark 设置好友备注
	ActionSetFriendRemark Action = "set_friend_remark"
	// ActionSetGroupAddOption 设置群加群选项
	ActionSetGroupAddOption Action = "set_group_add_option"
	// ActionSetGroupAddRequest 处理加群请求
	ActionSetGroupAddRequest Action = "set_group_add_request"
	// ActionSetGroupAdmin 设置群管理员
	ActionSetGroupAdmin Action = "set_group_admin"
	// ActionSetGroupAlbumMediaLike 点赞群相册媒体
	ActionSetGroupAlbumMediaLike Action = "set_group_album_media_like"
	// ActionSetGroupBan 群组禁言
	ActionSetGroupBan Action = "set_group_ban"
	// ActionSetGroupCard 设置群名片
	ActionSetGroupCard Action = "set_group_card"
	// ActionSetGroupKick 群组踢人
	ActionSetGroupKick Action = "set_group_kick"
	// ActionSetGroupKickMembers 批量踢出群成员
	ActionSetGroupKickMembers Action = "set_group_kick_members"
	// ActionSetGroupLeave 退出群组
	ActionSetGroupLeave Action = "set_group_leave"
	// ActionSetGroupName 设置群名称
	ActionSetGroupName Action = "set_group_name"
	// ActionSetGroupPortrait 设置群头像
	ActionSetGroupPortrait Action = "set_group_portrait"
	// ActionSetGroupRemark 设置群备注
	ActionSetGroupRemark Action = "set_group_remark"
	// ActionSetGroupRobotAddOption 设置群机器人加群选项
	ActionSetGroupRobotAddOption Action = "set_group_robot_add_option"
	// ActionSetGroupSearch 设置群搜索选项
	ActionSetGroupSearch Action = "set_group_search"
	// ActionSetGroupSign 群打卡
	ActionSetGroupSign Action = "set_group_sign"
	// ActionSetGroupSpecialTitle 设置专属头衔
	ActionSetGroupSpecialTitle Action = "set_group_special_title"
	// ActionSetGroupTodo 设置群待办
	ActionSetGroupTodo Action = "set_group_todo"
	// ActionSetGroupWholeBan 全员禁言
	ActionSetGroupWholeBan Action = "set_group_whole_ban"
	// ActionSetInputStatus 设置输入状态
	ActionSetInputStatus Action = "set_input_status"
	// ActionSetMsgEmojiLike 设置消息表情点赞
	ActionSetMsgEmojiLike Action = "set_msg_emoji_like"
	// ActionSetOnlineStatus 设置在线状态
	ActionSetOnlineStatus Action = "set_online_status"
	// ActionSetQQAvatar 设置QQ头像
	ActionSetQQAvatar Action = "set_qq_avatar"
	// ActionSetQQProfile 设置QQ资料
	ActionSetQQProfile Action = "set_qq_profile"
	// ActionSetRestart 重启服务
	ActionSetRestart Action = "set_restart"
	// ActionSetSelfLongnick 设置个性签名
	ActionSetSelfLongnick Action = "set_self_longnick"
	// ActionTestDownloadStream 测试下载流
	ActionTestDownloadStream Action = "test_download_stream"
	// ActionTransGroupFile 传输群文件
	ActionTransGroupFile Action = "trans_group_file"
	// ActionTranslateEn2zh 英文单词翻译
	ActionTranslateEn2zh Action = "translate_en2zh"
	// ActionUploadFileStream 上传文件流
	ActionUploadFileStream Action = "upload_file_stream"
	// ActionUploadGroupFile 上传群文件
	ActionUploadGroupFile Action = "upload_group_file"
	// ActionUploadImageToQunAlbum 上传图片到群相册
	ActionUploadImageToQunAlbum Action = "upload_image_to_qun_album"
	// ActionUploadPrivateFile 上传私聊文件
	ActionUploadPrivateFile Action = "upload_private_file"
)

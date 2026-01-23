package second_ecomInsight_getTaskResult_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightGetTaskResultResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondEcomInsightGetTaskResultData `json:"data"`
}
type Usage struct {
	// 总token数
	TotalTokens int64 `json:"total_tokens"`
	// 提示token数
	PromptTokens int64 `json:"prompt_tokens"`
	// 完成token数
	CompletionTokens int64 `json:"completion_tokens"`
	// prompt_tokens详情
	PromptTokensDetails *PromptTokensDetails `json:"prompt_tokens_details"`
	// 完成token详情
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details"`
}
type VideoUnderstandResponse struct {
	// 请求花费
	Usage *Usage `json:"usage"`
	// 预留对象
	Object string `json:"object"`
	// 结果信息
	Choices []ChoicesItem `json:"choices"`
	// 创建时间
	Created int64 `json:"created"`
	// 默认值
	ServiceTier string `json:"service_tier"`
	// 任务id
	Id string `json:"id"`
	// 任务模型
	Model string `json:"model"`
}
type VideoByPicRefAndVideoResponse struct {
	// 完整视频url
	FullVideo string `json:"full_video"`
	// 分段视频结构
	VideoSegments []VideoSegmentsItem `json:"video_segments"`
}
type SecondEcomInsightGetTaskResultData struct {
	// 消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
	// 文生视频返回值
	VideoByTextResponse *VideoByTextResponse `json:"video_by_text_response"`
	// 视频生成-基于首帧
	VideoByPicFirstFrameResponse *VideoByPicFirstFrameResponse `json:"video_by_pic_first_frame_response"`
	// 视频生成-基于首尾帧
	VideoByPicFirstLastFrameResponse *VideoByPicFirstLastFrameResponse `json:"video_by_pic_first_last_frame_response"`
	// 视频理解
	VideoUnderstandResponse *VideoUnderstandResponse `json:"video_understand_response"`
	// VideoPilot返回结构
	VideoByPicRefAndVideoResponse *VideoByPicRefAndVideoResponse `json:"video_by_pic_ref_and_video_response"`
	// 文生图结构
	PicByTextResponse *PicByTextResponse `json:"pic_by_text_response"`
}
type VideoByTextResponse struct {
	// 文件url
	FileUrl string `json:"file_url"`
	// 视频url
	VideoUrl string `json:"video_url"`
	// 尾帧url
	LastFrameUrl string `json:"last_frame_url"`
}
type VideoByPicFirstFrameResponse struct {
	// 文件url
	FileUrl string `json:"file_url"`
	// 视频url
	VideoUrl string `json:"video_url"`
	// 尾帧url
	LastFrameUrl string `json:"last_frame_url"`
}
type PromptTokensDetails struct {
	// 缓存tokens
	CachedTokens int64 `json:"cached_tokens"`
	// reason tokens
	ReasoningTokens int64 `json:"reasoning_tokens"`
}
type Message struct {
	// 内容
	Content string `json:"content"`
	// 原因内容
	ReasoningContent string `json:"reasoning_content"`
	// 角色
	Role string `json:"role"`
}
type ChoicesItem struct {
	// 结果在列表中的索引
	Index int32 `json:"index"`
	// 消息体
	Message *Message `json:"message"`
	// 对数概率信息（通常为 null）
	Logprobs string `json:"logprobs"`
	// 生成停止的原因，例如 "stop"
	FinishReason string `json:"finish_reason"`
}
type VideoSegmentsItem struct {
	// 视频url
	VideoUrl string `json:"video_url"`
	// 关键帧url
	KeyframeUrl string `json:"keyframe_url"`
	// 视频索引
	IndexInVideo int32 `json:"index_in_video"`
	// 生成的视频片段 URL
	SegmentId string `json:"segment_id"`
	// 自动生成或优化后的场景描述提示词
	ScenePrompt string `json:"scene_prompt"`
	// 当前片段的版本号
	VersionNum int32 `json:"version_num"`
}
type DataItem struct {
	// 图片url
	Url string `json:"url"`
	// 图片大小
	Size string `json:"size"`
}
type PicByTextResponse struct {
	// 详细信息
	Data []DataItem `json:"data"`
	// 模型
	Model string `json:"model"`
	// 花费token信息
	Usage *Usage `json:"usage"`
	// 创建时间
	Created int64 `json:"created"`
}
type VideoByPicFirstLastFrameResponse struct {
	// 文件url
	FileUrl string `json:"file_url"`
	// 视频url
	VideoUrl string `json:"video_url"`
	// 尾帧url
	LastFrameUrl string `json:"last_frame_url"`
}
type CompletionTokensDetails struct {
	// 缓存tokens
	CachedTokens int64 `json:"cached_tokens"`
	// reason tokens
	ReasoningTokens int64 `json:"reasoning_tokens"`
}

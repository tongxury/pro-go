package second_ecomInsight_createTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/second_ecomInsight_createTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightCreateTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SecondEcomInsightCreateTaskParam
}

func (c *SecondEcomInsightCreateTaskRequest) GetUrlPath() string {
	return "/second/ecomInsight/createTask"
}

func New() *SecondEcomInsightCreateTaskRequest {
	request := &SecondEcomInsightCreateTaskRequest{
		Param: &SecondEcomInsightCreateTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SecondEcomInsightCreateTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*second_ecomInsight_createTask_response.SecondEcomInsightCreateTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &second_ecomInsight_createTask_response.SecondEcomInsightCreateTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SecondEcomInsightCreateTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SecondEcomInsightCreateTaskRequest) GetParams() *SecondEcomInsightCreateTaskParam {
	return c.Param
}

type VideoUnderstandRequest struct {
	// 基础信息
	Base *Base `json:"base"`
	// 视频理解的帧数，会根据输入的帧数抽帧进行视频理解
	Fps *string `json:"fps"`
	// 理解视频素材
	VideoAsset *VideoAsset `json:"video_asset"`
	// 描述
	Text *string `json:"text"`
}
type VideoByPicRefAndVideoRequest struct {
	// 视频分辨率（例如 “1080p”“720p”）
	VideoResolution *string `json:"video_resolution"`
	// 模仿策略，例如'imitative' （模仿） 或 'creative' （创造）
	ImitationSetting *string `json:"imitation_setting"`
	// 提供的图片素材
	RefImageAssets []RefImageAssetsItem `json:"ref_image_assets"`
	// 视频宽高比（例如 “1:1”“16:9”）
	VideoRatio *string `json:"video_ratio"`
	// 可选（取值 1、2 或 3）分配的算力或质量等级，数值越高，可生成更长时长或更高质量的结果，耗费的时间和积分也更多
	TimeBudget *string `json:"time_budget"`
	// 视频素材
	VideoAsset *VideoAsset `json:"video_asset"`
	// 基础信息
	Base *Base `json:"base"`
	// 需要进行的视频风格的描述
	Text *string `json:"text"`
}
type ImageItem struct {
	// 文件，内容包括视频url、商品id、图片url、素材id等，根据type传入
	FileElement *string `json:"file_element"`
	// asset(素材，fileElement使用素材id)，path（图片或视频url，必须是公网可访问的），video_id（抖音视频id）， product_id（商品id（主图视频））
	Type *string `json:"type"`
}
type PicByTextRequest struct {
	// 基础信息
	Base *Base `json:"base"`
	// 需要进行的图片风格的描述
	Prompt *string `json:"prompt"`
	// 生成图片的大小：2K
	Size *string `json:"size"`
	// 是否生成水印
	Watermark *bool `json:"watermark"`
	// 是否生成组图 disabled - 不生成，auto - 生成
	SequentialImageGeneration *string `json:"sequential_image_generation"`
	// 图片素材
	Image []ImageItem `json:"image"`
}
type FirstFrameAsset struct {
	// 文件，内容包括视频url、商品id、图片url、素材id等，根据type传入
	FileElement *string `json:"file_element"`
	// asset(素材，fileElement使用素材id)，path（图片或视频url，必须是公网可访问的），video_id（抖音视频id）， product_id（商品id（主图视频））
	Type *string `json:"type"`
}
type Base struct {
	// 任务类型
	TaskType *string `json:"task_type"`
	// 任务名（不可重复）
	Name *string `json:"name"`
	// 首帧生成视频支持模型：Doubao-Seedance-1.0-pro-fast
	Model *string `json:"model"`
}
type RefImageAssetsItem struct {
	// asset(素材，fileElement使用素材id)，path（图片或视频url，必须是公网可访问的），video_id（抖音视频id）， product_id（商品id（主图视频））
	Type *string `json:"type"`
	// 文件，内容包括视频url、商品id、图片url、素材id等，根据type传入
	FileElement *string `json:"file_element"`
}
type VideoByTextRequest struct {
	// 基础信息
	Base *Base `json:"base"`
	// 文生视频、图生视频中在text后添加--参数的方式添加对应视频参数 --rt 比例、 --dur 视频时长、--fps 视频帧数、--rs 清晰度  --wm 是否包含水印 -- cf 是否固定摄像头
	Text *string `json:"text"`
}
type SecondEcomInsightCreateTaskParam struct {
	// 视频生成-基于首帧
	VideoByPicFirstFrameRequest *VideoByPicFirstFrameRequest `json:"video_by_pic_first_frame_request"`
	// 视频生成-基于首尾帧
	VideoByPicFirstLastFrameRequest *VideoByPicFirstLastFrameRequest `json:"video_by_pic_first_last_frame_request"`
	// 视频理解
	VideoUnderstandRequest *VideoUnderstandRequest `json:"video_understand_request"`
	// VideoPilot生成视频参数
	VideoByPicRefAndVideoRequest *VideoByPicRefAndVideoRequest `json:"video_by_pic_ref_and_video_request"`
	// 文生图
	PicByTextRequest *PicByTextRequest `json:"pic_by_text_request"`
	// 文生视频
	VideoByTextRequest *VideoByTextRequest `json:"video_by_text_request"`
}
type VideoByPicFirstFrameRequest struct {
	// 帧数
	Fps *string `json:"fps"`
	// 文生视频、图生视频中在text后添加--参数的方式添加对应视频参数。   --rt 比例 --dur 视频时长  --fps 视频帧数  --rs 清晰度  --wm 是否包含水印 -- cf 是否固定摄像头
	Text *string `json:"text"`
	// 首帧素材
	FirstFrameAsset *FirstFrameAsset `json:"first_frame_asset"`
	// 基础信息
	Base *Base `json:"base"`
}
type LastFrameAsset struct {
	// 文件，内容包括视频url、商品id、图片url、素材id等，根据type传入
	FileElement *string `json:"file_element"`
	// asset(素材，fileElement使用素材id)，path（图片或视频url，必须是公网可访问的），video_id（抖音视频id）， product_id（商品id（主图视频））
	Type *string `json:"type"`
}
type VideoByPicFirstLastFrameRequest struct {
	// 基础信息
	Base *Base `json:"base"`
	// 首帧素材
	FirstFrameAsset *FirstFrameAsset `json:"first_frame_asset"`
	// 尾帧素材
	LastFrameAsset *LastFrameAsset `json:"last_frame_asset"`
	// 帧数
	Fps *string `json:"fps"`
	// 文生视频、图生视频中在text后添加--参数的方式添加对应视频参数。   --rt 比例 --dur 视频时长  --fps 视频帧数  --rs 清晰度  --wm 是否包含水印 -- cf 是否固定摄像头
	Text *string `json:"text"`
}
type VideoAsset struct {
	// asset(素材，fileElement使用素材id)，path（图片或视频url，必须是公网可访问的），video_id（抖音视频id）， product_id（商品id（主图视频））
	Type *string `json:"type"`
	// 文件，内容包括视频url、商品id、图片url、素材id等，根据type传入
	FileElement *string `json:"file_element"`
}

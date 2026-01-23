package service

import (
	"context"
	"encoding/json"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
)

func (t ProjService) AnalyzeSubSegment() {

	ctx := context.Background()

	hits, err := t.data.Mongo.TemplateSegment.List(ctx, bson.M{"status": "processing_tagged"}, options.Find().SetLimit(1))
	if err != nil {
		log.Errorw("Search err", err, "index", "items")
		return
	}

	if len(hits) == 0 {
		return
	}

	//t.analyzeResourceSegmentBySeed(ctx, hits[0])
	wg.WaitGroup(ctx, hits, t.analyzeSubSegment)
}

func (t ProjService) analyzeSubSegment(ctx context.Context, x *projpb.ResourceSegment) error {

	return t.analyzeSubSegmentByGemini(ctx, x)
}

// 输出要求:
// - 详细描述短片的运镜(参考运镜示例)
//
// ===
// 运镜示例:
//
// 一、基础动态运镜:
//
// 01. 无限拉远
// 提示词： 无限拉远，空中视角，超精细
// 应用： 结尾引出全景、宇宙感叙事
// 效果： 人物→星空无缝递进，制造空间转场震撼
//
// 02. 镜头穿窗
// 提示词： 镜头穿过窗户，无缝过渡
// 应用： 空间转换、回忆/新世界切入
// 效果： 如入幻境，打破空间壁垒
//
// 03. 时间静止环绕
// 提示词： 子弹时间，时间静止，环绕镜头
// 应用： 关键时刻、情绪爆发节点
// 效果： 画面定格+环绕，制造高潮张力
//
// 04. 穿越四季
// 提示词： 时间流逝，镜头前移，四季变化
// 应用： 叙事转场、成长表达
// 效果： 诗意传递时间流转感
//
// 05. 城市穿梭
// 提示词： 无人机穿城，城市峡谷，电影感
// 应用： 都市故事开场、切入
// 效果： 科技感+动感冲击，展现城市脉络
//
// 06. 第一人称视角
// 提示词： 第一人称，沉浸视角，平滑镜头
// 应用： 代入剧情、互动叙事
// 效果： 强沉浸感，让观众“亲历”故事
//
// 07. 镜面穿越
// 提示词： 镜像世界，传送门转场
// 应用： 梦境、意识、幻想世界穿越
// 效果： 科技+魔法混合风，创意转场
//
// 08. 向下坠落
// 提示词： 从天而降，垂直俯冲
// 应用： 高潮、惊吓场面
// 效果： 营造失重、紧张冲突感
//
// 09. 超近特写推进
// 提示词： 快速推进，戏剧特写
// 应用： 惊讶/悲痛等情绪表达
// 效果： 聚焦表情，强化情感共鸣
//
// 10. 螺旋上升
// 提示词： 螺旋上升，旋涡动感
// 应用： 突破、觉醒、升华象征
// 效果： 强视觉冲击，适配结尾飞升场景
//
// 二、创意叙事镜头
//
// 11. 灵魂出窍
// 提示词： 灵魂离体，星体投射
// 应用： 梦境、灵性、失落感场景
// 效果： 哲思+梦幻视觉，传递精神抽离感
//
// 12. 万花筒转场
// 提示词： 万花筒，镜像对称，旋转
// 应用： 幻觉表达、艺术剪辑
// 效果： 绚烂神秘，具设计感的转场
//
// 13. 拉远再推进
// 提示词： 拉远再拉近，动态焦点切换
// 应用： 信息重启、故事翻篇
// 效果： 视觉节奏鲜明，强化叙事层次
//
// 14. 火焰视角
// 提示词： 点燃火焰，火花飞溅，主观视角
// 应用： 激烈情绪转折点
// 效果： 强视觉冲击，放大情绪张力
//
// 15. 碎裂视角
// 提示词： 玻璃破碎，情绪崩裂
// 应用： 心理崩溃、打击场景
// 效果： 碎片隐喻，强化冲击感
//
// 16. 记忆倒退
// 提示词： 倒带回忆，逆向运镜
// 应用： 回忆穿插、旧时光叙事
// 效果： 瞬间时空转换，增强代入
//
// 17. 梦境模糊
// 提示词： 梦境模糊，发光虚焦，漂浮光点
// 应用： 入梦/出梦、昏迷场景
// 效果： 温柔唯美，营造超现实氛围
//
// 18. 飞跃古城
// 提示词： 俯瞰古城，空中视角
// 应用： 历史、民族、遗迹题材
// 效果： 震撼大气，提升格局感
//
// 19. 空间碎片漂浮
// 提示词： 漂浮碎片，环绕粒子
// 应用： 意识流、宇宙感构图
// 效果： 科幻+信息化交融，构建宏观意境
//
// 20. 空间撕裂
// 提示词： 空间裂缝，故障转场
// 应用： 科技冲击、超现实变换
// 效果： 打破线性空间，制造次元突破感
func (t ProjService) analyzeSubSegmentByGemini(ctx context.Context, x *projpb.ResourceSegment) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "analyzeSubSegmentByGemini",
		"item", x.XId,
	))

	genaiClient := t.data.GenaiFactory.Get()

	segment, err := videoz.GetSegmentByUrl(ctx, x.Root.Url, x.TimeStart, x.TimeEnd)
	if err != nil {
		logger.Errorw("Get segment by url", x.Root.Url, "err", err)
		return err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeObject,
			Required: []string{"segments"},
			Properties: map[string]*genai.Schema{
				"formula": {
					Type:        genai.TypeString,
					Description: "爆款公式",
				},
				"segments": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type:        genai.TypeObject,
						Description: "分镜结果",
						Required:    []string{"cameraMovementDesc", "formula", "timeStart", "timeEnd"},
						Properties: map[string]*genai.Schema{
							"timeStart": {
								Type:        genai.TypeNumber,
								Description: "当前分段脚本的开始时间戳(秒, 示例：1.22)",
							},
							"timeEnd": {
								Type:        genai.TypeNumber,
								Description: "当前分段脚本的结束时间戳(秒, 示例：1.23)",
							},
							"formula": {
								Type:        genai.TypeString,
								Description: "当前分段所属的爆款公式环节",
							},
							"cameraMovementDesc": {
								Type:        genai.TypeString,
								Description: "运镜描述",
							},
						},
					},
				},
			},
		},
	}

	content, err := genaiClient.AnalyzeVideo(ctx, gemini.AnalyzeVideoRequest{
		VideoBytes: segment.Content,
		Prompt: `

作为一名资深短视频编导，请帮我分析拆解视频爆款内容并且拆解包含相关场景、道具、提及人物，总结爆款公式。
并根据爆款公式帮我将视频拆分成一系列独立的短片。
`,

		Config: config,
	})
	if err != nil {
		logger.Errorw("AnalyzeVideo err ", err)
		return err
	}

	var res projpb.ResourceSegment
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		return err
	}

	segments := len(res.Segments)

	if segments > 6 {
		logger.Errorw("Segments count err", err, "segments", segments)
		return fmt.Errorf("invalid segments number: %d", segments)
	}
	if segments == 0 {
		logger.Errorw("Segments count err", err, "segments", segments)
		return fmt.Errorf("invalid segments number: %d", segments)
	}

	for i := range res.Segments {

		xx := res.Segments[i]

		fmt.Println(xx.TimeStart, xx.Description)

		frame, err := videoz.GetFrame(segment.Content, xx.TimeStart+0.1)
		if err != nil {
			return err
		}

		res.Segments[i].StartFrame, err = t.data.TOS.Put(ctx, tos.PutRequest{
			Bucket:  "yoozyres",
			Content: frame,
			Key:     helper.CreateUUID() + ".jpg",
		})

		if err != nil {
			logger.Errorw("Put err ", err)
			return err
		}

		endFrame, err := videoz.GetFrame(segment.Content, xx.TimeEnd-0.1)
		if err != nil {
			return err
		}

		res.Segments[i].EndFrame, err = t.data.TOS.Put(ctx, tos.PutRequest{
			Bucket:  "yoozyres",
			Content: endFrame,
			Key:     helper.CreateUUID() + ".jpg",
		})

		if err != nil {
			logger.Errorw("Put err ", err)
			return err
		}

	}

	updated, err := t.data.Mongo.TemplateSegment.UpdateByIDXX(ctx, x.XId, bson.M{
		"$set": bson.M{
			"status":   "completed",
			"segments": res.Segments,
		},
	})

	logger.Debugw("updated", "", "item", x.XId, "updated", updated, "segments", segments)

	if err != nil {
		logger.Errorw("UpdateByIDXX err", err, "index", x.XId)
		return err
	}

	return nil
}

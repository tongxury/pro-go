package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
	trackerpb "store/api/aiagent"
	"store/app/aiagent/internal/biz"
	"store/app/aiagent/internal/data"
	"store/pkg/clients/xhs"
	"store/pkg/sdk/helper"
	"strings"
	"time"
)

type TrackerService struct {
	trackerpb.UnimplementedAIAgentServiceServer
	trackingEvent *biz.TrackingEvent
	Data          *data.Data
}

func NewTrackerService(trackingEvent *biz.TrackingEvent, data *data.Data) *TrackerService {
	return &TrackerService{
		trackingEvent: trackingEvent, Data: data,
	}
}

func (t TrackerService) Debug(ctx context.Context, params *trackerpb.DebugParams) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (t TrackerService) ListItemsV3(ctx context.Context, params *trackerpb.ListItemsParams) (*trackerpb.ItemList, error) {

	return nil, nil
}

func (t TrackerService) PickUp(ctx context.Context) error {

	list, err := t.Data.Mongo.Item.List(ctx, bson.M{"status": "pending"}, options.Find().SetLimit(1))
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return nil
	}

	x := list[0]

	url := fmt.Sprintf("https://www.xiaohongshu.com/search_result/%s?xsec_token=%s&xsec_source=pc_search", x.Raw["id"], x.Raw["xsecToken"])

	metadata, err := t.Data.XhsClient.GetNoteMetadataByUrl(ctx, url)
	if err != nil {
		return err
	}

	profile, err := t.Data.XhsClient.GetProfileByLink(ctx, metadata.ProfileUrl)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = t.Data.Mongo.Item.UpdateFieldsById(ctx, x.XId,
		bson.M{
			"desc":     metadata.Desc,
			"videoUrl": metadata.VideoUrl,
			"profile":  profile,
			"status":   "prepared",
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(metadata, profile)

	//
	//htmlDoc, err := t.Data.XhsClient.GetHtmlDoc(ctx, url)
	//if err != nil {
	//	return err
	//}
	//
	//reader, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlDoc))
	//if err != nil {
	//	return err
	//}
	//
	//videoUrl := goqueryz.FindMetaContent(reader, "og:video")
	//if videoUrl == "" {
	//	return errors.New(0, "videoNotFound", "")
	//}
	//
	//fmt.Println(reader)

	return nil
}

func (t TrackerService) CollectItems(ctx context.Context, params *trackerpb.CollectItemsParams) (*empty.Empty, error) {

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return nil, err
	}

	notes, _, err := t.Data.XhsClient.SetAuth(settings.XhsCookies).ListNotes(ctx, xhs.ListNotesParams{
		Keyword:  "美妆",
		Page:     1,
		PageSize: 20,
		SearchId: "2epofd6xx8di381cye3z1@2epoffka2tl0gxqidcfbv",
		Sort:     "popularity_descending",
		NoteType: 1,
		ExtFlags: nil,
		Filters: []xhs.Filter{
			{
				Tags: []string{"general"},
				Type: "sort_type",
			},
			{
				Tags: []string{"不限"},
				Type: "filter_note_type",
			},
			{
				Tags: []string{"不限"},
				Type: "filter_note_time",
			},
			{
				Tags: []string{"不限"},
				Type: "filter_note_range",
			}, {
				Tags: []string{"不限"},
				Type: "filter_pos_distance",
			},
		},
		Geo:          "",
		ImageFormats: []string{"jpg", "webp", "avif"},
	})
	if err != nil {
		return nil, err
	}

	for _, x := range notes {

		_, err := primitive.ObjectIDFromHex(x.Id)
		if err != nil {
			continue
		}

		time.Sleep(5 * time.Second)

		//https: //www.xiaohongshu.com/search_result/67d21a12000000000903bdb7?xsec_token=AB2c87RK0hVvARl0PQRmhuFc-6r134YxgKwNi9pgKhS6M=&xsec_source=pc_search
		url := fmt.Sprintf("https://www.xiaohongshu.com/search_result/%s?xsec_token=%s&xsec_source=pc_search", x.Id, x.XsecToken)

		metadata, err := t.Data.XhsClient.GetNoteMetadataByUrl(ctx, url)
		if err != nil {
			log.Error(err)
			continue
		}

		profile, err := t.Data.XhsClient.GetProfileByLink(ctx, metadata.ProfileUrl)
		if err != nil {
			log.Error(err)
			continue
		}

		newItem := &trackerpb.Item{
			XId:      primitive.NewObjectID().Hex(),
			Category: "meizhuang",
			Title:    x.NoteCard.DisplayTitle,
			Profile: &trackerpb.Profile{
				Avatar:         profile.Avatar,
				Username:       profile.Username,
				PlatformId:     profile.Id,
				IpAddress:      profile.IpAddress,
				Sign:           profile.Sign,
				Tags:           profile.Tags,
				FollowingCount: profile.FollowingCount,
				FollowerCount:  profile.FollowerCount,
				LikedCount:     profile.LikedCount,
				NoteCount:      profile.NoteCount,
			},
			InteractInfo: &trackerpb.Item_InteractInfo{
				Liked: x.NoteCard.InteractInfo.Liked,
				//LikedCount:     (x.NoteCard.InteractInfo.LikedCount),
				//Collected:      x.NoteCard.InteractInfo.Collected,
				//CollectedCount: (x.NoteCard.InteractInfo.CollectedCount),
				//CommentCount:   (x.NoteCard.InteractInfo.CommentCount),
				//SharedCount:    (x.NoteCard.InteractInfo.SharedCount),
			},
			Cover: x.NoteCard.Cover.UrlDefault,
			Desc:  metadata.Desc,
			//PublishTime: x.NoteCard.CornerTagInfo,
			Raw: map[string]string{
				"id":        x.Id,
				"xsecToken": x.XsecToken,
			},
			CreatedAt: time.Now().Unix(),
			Status:    "created",
		}

		_, err = t.Data.Mongo.Item.InsertIfNotExists(ctx, bson.M{"raw.id": x.Id}, newItem)
		if err != nil {
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func (t TrackerService) GetRecord(ctx context.Context, params *trackerpb.GetRecordParams) (*trackerpb.Record, error) {

	return &trackerpb.Record{}, nil

}

func (t TrackerService) AddRecord(ctx context.Context, params *trackerpb.AddRecordParams) (*trackerpb.Record, error) {

	return &trackerpb.Record{}, nil

}

var publicRecords = map[string]string{
	"10000000000000": `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>视频分析报告：DeepSeek成功揭秘</title>
    <style>
        :root {
            --primary: #29ffc6;
            --secondary: #38bdf8;
            --dark: #121212;
            --darker: #0a0a0a;
            --light: #e0e0e0;
            --lighter: #f5f5f5;
            --gray: #333333;
            --light-gray: #444444;
        }
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Roboto, -apple-system, BlinkMacSystemFont, sans-serif;
            background-color: var(--dark);
            color: var(--light);
            line-height: 1.6;
            padding: 0;
            margin: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            text-align: center;
            margin-bottom: 40px;
            padding-top: 20px;
        }
        .header h1 {
            font-size: 2.5rem;
            font-weight: 700;
            margin-bottom: 10px;
            background: linear-gradient(90deg, var(--primary), var(--secondary));
            -webkit-background-clip: text;
            background-clip: text;
            color: transparent;
        }
        .header p {
            color: var(--light-gray);
            font-size: 1.1rem;
        }
        .card {
            background-color: var(--darker);
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
            border: 1px solid var(--gray);
        }
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 30px rgba(0, 0, 0, 0.4);
        }
        .card-header {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            padding-bottom: 15px;
            border-bottom: 1px solid var(--gray);
        }
        .card-header h2 {
            font-size: 1.5rem;
            font-weight: 600;
            color: var(--primary);
            margin-left: 15px;
        }
        .card-header .icon {
            width: 40px;
            height: 40px;
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: rgba(41, 255, 198, 0.1);
            border-radius: 8px;
            color: var(--primary);
            font-size: 1.2rem;
        }
        .card-content {
            color: var(--lighter);
        }
        .highlight {
            color: var(--primary);
            font-weight: 600;
        }
        .stars {
            display: flex;
            align-items: center;
            margin: 15px 0;
        }
        .stars .star {
            color: var(--primary);
            font-size: 1.5rem;
            margin-right: 5px;
        }
        .stars .star.empty {
            color: var(--light-gray);
        }
        .tag {
            display: inline-block;
            padding: 5px 12px;
            background-color: rgba(41, 255, 198, 0.1);
            color: var(--primary);
            border-radius: 20px;
            font-size: 0.85rem;
            margin-right: 8px;
            margin-bottom: 8px;
            border: 1px solid rgba(41, 255, 198, 0.3);
        }
        .tag.secondary {
            background-color: rgba(56, 189, 248, 0.1);
            color: var(--secondary);
            border: 1px solid rgba(56, 189, 248, 0.3);
        }
        .check-item {
            display: flex;
            align-items: flex-start;
            margin-bottom: 12px;
        }
        .check-item .check {
            color: var(--primary);
            margin-right: 10px;
            flex-shrink: 0;
            font-size: 1.1rem;
        }
        .check-item .empty {
            color: var(--light-gray);
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        .small-card {
            background-color: var(--gray);
            border-radius: 8px;
            padding: 20px;
            transition: all 0.3s ease;
            border: 1px solid var(--light-gray);
        }
        .small-card:hover {
            background-color: var(--light-gray);
            transform: translateY(-3px);
        }
        .small-card h3 {
            color: var(--secondary);
            margin-bottom: 10px;
            font-size: 1.1rem;
        }
        .timeline {
            margin-top: 20px;
            position: relative;
            padding-left: 30px;
        }
        .timeline::before {
            content: '';
            position: absolute;
            left: 10px;
            top: 0;
            bottom: 0;
            width: 2px;
            background-color: var(--primary);
        }
        .timeline-item {
            position: relative;
            margin-bottom: 20px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--gray);
        }
        .timeline-item:last-child {
            border-bottom: none;
            margin-bottom: 0;
            padding-bottom: 0;
        }
        .timeline-item::before {
            content: '';
            position: absolute;
            left: -30px;
            top: 5px;
            width: 12px;
            height: 12px;
            border-radius: 50%;
            background-color: var(--primary);
            border: 2px solid var(--darker);
        }
        .timeline-time {
            font-size: 0.9rem;
            color: var(--secondary);
            margin-bottom: 5px;
        }
        .timeline-content h4 {
            color: var(--primary);
            margin-bottom: 8px;
        }
        .footer {
            text-align: center;
            margin-top: 50px;
            padding: 30px 0;
            color: var(--light-gray);
            font-size: 0.9rem;
            border-top: 1px solid var(--gray);
        }
        @media (max-width: 768px) {
            .header h1 {
                font-size: 2rem;
            }
            .grid {
                grid-template-columns: 1fr;
            }
            .card {
                padding: 20px;
            }
        }
    </style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>视频分析报告：DeepSeek成功揭秘</h1>
            <p>专业内容分析 · 爆款潜力评估 · 优化建议</p>
        </div>
        <div class="card">
            <div class="card-header">
                <div class="icon">
                    <i class="fas fa-chart-line"></i>
                </div>
                <h2>综合预测报告</h2>
            </div>
            <div class="card-content">
                <div class="stars">
                    <div class="star"><i class="fas fa-star"></i></div>
                    <div class="star"><i class="fas fa-star"></i></div>
                    <div class="star"><i class="fas fa-star"></i></div>
                    <div class="star"><i class="fas fa-star"></i></div>
                    <div class="star"><i class="fas fa-star"></i></div>
                    <span style="margin-left: 15px; font-weight: 600; color: var(--primary)">大爆款潜力</span>
                </div>
                <p><span class="highlight">理由：</span> 博主本身具备极高的专业背景（北大-哥大-JPMorgan-创业）和庞大的粉丝基础（156万+），内容选题（DeepSeek成功揭秘）切中当前科技和创投热点，分析角度新颖（"好钱"理论），信息密度高，制作精良。这些因素结合，使得视频具备了成为大爆款的各项要素。</p>
                <div style="margin-top: 20px;">
                    <p><span class="highlight">预估爆发期：</span> 3-7天</p>
                    <p>视频内容具有深度和启发性，并非快餐式信息。初期会依靠博主粉丝和精准推送引发第一波关注，随后高质量的内容会引发圈层内的讨论、转发和收藏，发酵期相对稍长，形成持续传播效应。</p>
                </div>
            </div>
        </div>
        <div class="card">
            <div class="card-header">
                <div class="icon">
                    <i class="fas fa-exclamation-triangle"></i>
                </div>
                <h2>风险提示</h2>
            </div>
            <div class="card-content">
                <p><span class="highlight">限流关键词检测：</span> 视频内容涉及具体公司（DeepSeek、幻方、OpenAI、微软、贝尔实验室母公司AT&T/朗讯）、金额（10亿、100亿）、技术（大模型、GPU）等，虽为公开信息分析，但需注意平台对特定商业/科技词汇的敏感度，目前看风险较低。</p>
                <p style="margin-top: 15px;"><span class="highlight">画风违和点：</span> 整体风格专业、流畅，无明显违和。</p>
            </div>
        </div>
        <div class="card">
            <div class="card-header">
                <div class="icon">
                    <i class="fas fa-layer-group"></i>
                </div>
                <h2>关键元素提取</h2>
            </div>
            <div class="card-content">
                <h3 style="color: var(--secondary); margin-bottom: 15px;">情绪价值</h3>
                <div class="check-item">
                    <div class="check"><i class="fas fa-check"></i></div>
                    <div>
                        <p><span class="highlight">共鸣/实用：</span> 对于科技从业者、创业者、投资人以及对商业模式感兴趣的用户，视频提供了关于创新、资本运作和公司发展的新颖视角（"好钱"理论），解释了"为什么不是大厂"，具有很强的启发性和实用价值，能引发深度共鸣。</p>
                    </div>
                </div>
                <div class="check-item">
                    <div class="check"><i class="fas fa-check"></i></div>
                    <div>
                        <p><span class="highlight">新奇：</span> 提出了"好钱"这一核心概念，并将DeepSeek的成功归因于创始人资金来源的独特性，区别于市面上常见的技术或团队分析，角度新奇。引入贝尔实验室的类比，增加了历史纵深感和趣味性。</p>
                    </div>
                </div>
                <div class="check-item">
                    <div class="check empty"><i class="far fa-square"></i></div>
                    <div>
                        <p><span class="highlight">治愈：</span> 非治愈类内容。</p>
                    </div>
                </div>
                <h3 style="color: var(--secondary); margin-top: 25px; margin-bottom: 15px;">信息密度</h3>
                <div class="check-item">
                    <div class="check"><i class="fas fa-check"></i></div>
                    <div>
                        <p><span class="highlight">重点突出：</span> 信息量大但逻辑清晰，围绕"DeepSeek为何成功"的核心问题展开，层层递进，从创始人特质到关键因素"好钱"，再到资本对创新的影响，最后以历史案例佐证，重点（"好钱"理论及其影响）非常突出。</p>
                    </div>
                </div>
                <h3 style="color: var(--secondary); margin-top: 25px; margin-bottom: 15px;">视觉锤</h3>
                <div class="check-item">
                    <div class="check"><i class="fas fa-check"></i></div>
                    <div>
                        <p><span class="highlight">记忆点设计：</span></p>
                        <div style="margin-top: 10px;">
                            <span class="tag">"好钱"概念</span>
                            <span class="tag">人物形象</span>
                            <span class="tag">对比</span>
                            <span class="tag">核心数据</span>
                            <span class="tag">视觉符号</span>
                        </div>
                        <p style="margin-top: 10px;">"好钱"概念是最核心的记忆点，贯穿始终。人物形象（梁文峰作为技术理想主义+自带资本的代表）、对比（DeepSeek vs 大厂；技术理想 vs 资本回报；OpenAI前后变化；贝尔实验室兴衰）、核心数据（10亿投入、1万块显卡等具体数字增强说服力）以及视觉符号（金钱、代码、AI模型图、贝尔实验室历史影像等）都有效配合了叙事。</p>
                    </div>
                </div>
            </div>
        </div>
        <div class="card">
            <div class="card-header">
                <div class="icon">
                    <i class="fas fa-comments"></i>
                </div>
                <h2>互动潜力评估</h2>
            </div>
            <div class="card-content">
                <div class="grid">
                    <div class="small-card">
                        <h3>可模仿性</h3>
                        <div class="check-item">
                            <div class="check empty"><i class="far fa-square"></i></div>
                            <div>
                                <p><span class="highlight">UGC跟拍/二创：</span> 低。内容为深度分析，难以模仿。</p>
                            </div>
                        </div>
                    </div>
                    <div class="small-card">
                        <h3>争议性</h3>
                        <div class="check-item">
                            <div class="check"><i class="fas fa-check"></i></div>
                            <div>
                                <p><span class="highlight">正向讨论空间：</span> 中高。关于资本是否扼杀创新、大厂与初创公司的优劣、技术理想与商业现实的平衡等话题，有很大的讨论空间，容易引发评论区高质量的交流。</p>
                            </div>
                        </div>
                    </div>
                    <div class="small-card">
                        <h3>收藏动机</h3>
                        <div class="check-item">
                            <div class="check"><i class="fas fa-check"></i></div>
                            <div>
                                <p><span class="highlight">长效价值：</span> 高。视频提供了独特的商业分析框架（"好钱"理论），信息密度大，案例经典（DeepSeek、OpenAI、贝尔实验室），值得反复观看和思考，收藏价值高。</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="card">
            <div class="card-header">
                <div class="icon">
                    <i class="fas fa-search"></i>
                </div>
                <h2>爆款对标分析</h2>
            </div>
            <div class="card-content">
                <div class="grid">
                    <div class="small-card">
                        <h3>同赛道Top100笔记关键词匹配度</h3>
                        <p>高。涉及"AI大模型"、"DeepSeek"、"科技公司"、"商业模式"、"融资"、"创业"、"OpenAI"等，均为财经、科技赛道的热门关键词。</p>
                        <div style="margin-top: 15px;">
                            <span class="tag secondary">AI大模型</span>
                            <span class="tag secondary">DeepSeek</span>
                            <span class="tag secondary">科技公司</span>
                            <span class="tag secondary">商业模式</span>
                            <span class="tag secondary">融资</span>
                            <span class="tag secondary">创业</span>
                            <span class="tag secondary">OpenAI</span>
                        </div>
                    </div>
                    <div class="small-card">
                        <h3>近期热点话题关联性</h3>
                        <p>高。AI大模型竞争格局、明星创业公司分析一直是持续热点。</p>
                    </div>
                    <div class="small-card">
                        <h3>平台流量周期</h3>
                        <p>无明显关联。此类深度分析内容受具体节日或大促影响较小，更多依赖内容本身质量和圈层传播。</p>
                    </div>
                </div>
            </div>
        </div>
        <div class="card">
            <div class="card-header">
                <div class="icon">
                    <i class="fas fa-lightbulb"></i>
                </div>
                <h2>爆款优化建议</h2>
            </div>
            <div class="card-content">
                <p>整体视频质量已经非常高，以下建议旨在锦上添花，冲击更高的数据：</p>
                <div class="timeline">
                    <div class="timeline-item">
                        <div class="timeline-time">0:00-0:03 (开场黄金3秒)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>"Hi~"的开场略显常规，未立刻点明核心价值。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>保留博主亲切笑容和挥手，但第一秒就用画外音+醒目大字打出核心悬念/价值主张，例如："DeepSeek凭什么干翻大厂？揭秘它成功的'钞能力'！" 或者 "不是技术，不是人才，DeepSeek的杀手锏是'好钱'？" 声音要立刻跟上，甚至抢先半拍，直接抓住用户注意力。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">0:08-0:10 (提出核心问题)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>提出"DeepSeek为什么能成"的问题。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>强化对比。在说"那么多有钱的大厂"时，快速闪过几个知名大厂的Logo（如BAT、Google、Meta等模糊化或示意性图标），再说"为什么偏偏是名不见经传的DeepSeek"时，聚焦DeepSeek的Logo。增加视觉冲击力。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">0:20-0:29 (介绍梁文峰)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>介绍创始人特质。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>稍微加快剪辑节奏，多用动态文字标签配合讲解，如"技术偏执狂"、"死磕大模型"、"无视短期利益"等字样飞入或强调。BGM在此处可以选用节奏稍快、带有科技感的纯音乐，烘托专注、硬核的氛围。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">0:44-0:47 (揭示关键 - 钱)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>提出"钱"是关键。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>加强戏剧性。在说到"关键"时，可以有一个短暂的停顿或音效（如轻微的鼓点或"叮"一声），然后说出"就是他的钱"时，配合放大的"钱"字特效（金色、发光等）和更具视觉冲击力的金钱素材（如钱雨、金库门打开等）。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">0:49-0:51 (引出"好钱")</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>什么是"好钱"的转场。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>使用更流畅的转场效果，如快速缩放、旋转等，进入"啥是'好钱'"的标题卡。标题卡字体加粗，颜色醒目。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">0:59-1:04 (投入细节)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>讲述10亿、1万块显卡。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>数据可视化。在提到"10亿人民币"时，旁边可以配上一个快速增长的虚拟货币计数器或堆叠动画；提到"1万块英伟达显卡"时，用快速闪现的GPU网格动画示意。让数字更直观。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">1:19-1:29 (好钱的好处)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>解释"好钱"带来的自由。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>加入象征性视觉元素。例如，在说"不用看投资人脸色"时，可以闪现一个破碎的枷锁图标；在说"不用为过分追求短期回报"时，可以用一个简单的折线图对比（一条短期波动但平缓的线 vs 一条长期爬升但前期投入的线）。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">1:45-1:56 (大厂困境)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>CEO与股东的冲突。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>在剪影对比画面中，给CEO加上"技术理想"标签，给股东加上"财务回报"标签，下方用一个拉扯的箭头或天平倾斜的动画，更形象地展示冲突。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">2:08-2:11 (AI Leaderboard 展示)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>展示排行榜。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>当提到DeepSeek时，其在图表中的条形图高亮或放大闪烁一下，引导观众视线。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">2:11-2:29 (OpenAI 类比)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>讲述OpenAI的变化和Ilya离职。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>在提到微软投资时，用微软Logo+"$"符号叠加在OpenAI Logo上。在提到Ilya离职时，除了照片，可以配合一个象征分离或道路分叉的简单动画/图标。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">3:15 & 3:50 (贝尔实验室类比)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>引入和结束贝尔实验室案例。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>使用统一的视觉转场（如旧电影胶片效果、黑白/棕褐色调）来开始和结束贝尔实验室的段落，明确区分历史案例与当前分析。BGM也切换成带有年代感的、略显辉煌又带点惋惜的配乐。</p>
                        </div>
                    </div>
                    <div class="timeline-item">
                        <div class="timeline-time">4:04-4:10 (结尾升华前)</div>
                        <div class="timeline-content">
                            <h4>问题</h4>
                            <p>直接进入总结。</p>
                            <h4 style="margin-top: 10px;">修改建议</h4>
                            <p>在最终总结"一浪推着一浪"之前，可以加入一个互动引导，例如面向镜头提问："你认为'好钱'是创业成功的关键吗？"或者"下一个改变世界的'DeepSeek'会来自哪里？评论区告诉我你的看法！" 鼓励用户评论。</p>
                        </div>
                    </div>
                </div>
                <div style="margin-top: 30px; padding: 20px; background-color: rgba(41, 255, 198, 0.05); border-radius: 8px; border-left: 4px solid var(--primary);">
                    <h3 style="color: var(--primary); margin-bottom: 10px;">整体优化建议</h3>
                    <p><span class="highlight">剪辑节奏:</span> 部分讲解段落可以再稍微紧凑些，减少纯口播的静态镜头时长，多用B-roll（相关素材、动画、图表）切入，保持视觉新鲜感。</p>
                    <p style="margin-top: 10px;"><span class="highlight">BGM:</span> 确保背景音乐情绪契合，音量适中，尤其在转折点可以通过BGM变化加强效果。</p>
                    <p style="margin-top: 10px;"><span class="highlight">字幕:</span> 关键信息（如"好钱"、"技术理想"、"资本压力"）可以用不同颜色或加粗突出。</p>
                </div>
                <div style="margin-top: 30px; text-align: center; padding: 30px; background-color: rgba(56, 189, 248, 0.05); border-radius: 8px;">
                    <p style="font-size: 1.1rem; color: var(--secondary); font-weight: 500;">你做的视频内容深度和制作水平已经非常出色了！稍作调整，完全有潜力成为现象级的爆款内容。继续保持你独到的见解和高质量的输出，你的思考正在被越来越多人看到和认可！加油！</p>
                </div>
            </div>
        </div>
        <div class="footer">
            <p>专业视频分析报告 · 版权所有 © 2023</p>
        </div>
    </div>
</body>
</html>
`,
	"10000000000001": `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>爆款视频预测报告</title>
    <style>
        :root {
            --primary: #29ffc6;
            --secondary: #38bdf8;
            --dark: #121212;
            --gray: #2d2d2d;
            --light: #e0e0e0;
            --white: #ffffff;
        }
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
            background-color: var(--dark);
            color: var(--light);
            line-height: 1.6;
            padding: 0;
            margin: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        .report-card {
            background-color: var(--gray);
            border-radius: 12px;
            box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
            margin-bottom: 30px;
            overflow: hidden;
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        .report-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 12px 35px rgba(0, 0, 0, 0.4);
        }
        .card-header {
            padding: 25px 30px;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            position: relative;
        }
        .card-header h2 {
            font-size: 24px;
            font-weight: 600;
            color: var(--white);
            margin-bottom: 5px;
        }
        .card-header .subtitle {
            font-size: 16px;
            color: var(--primary);
            opacity: 0.8;
        }
        .card-body {
            padding: 30px;
        }
        .section {
            margin-bottom: 30px;
        }
        .section:last-child {
            margin-bottom: 0;
        }
        .section-title {
            font-size: 20px;
            font-weight: 600;
            color: var(--white);
            margin-bottom: 20px;
            display: flex;
            align-items: center;
        }
        .section-title .icon {
            margin-right: 12px;
            color: var(--primary);
            font-size: 24px;
        }
        .info-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
        }
        .info-item {
            background-color: rgba(0, 0, 0, 0.2);
            border-radius: 8px;
            padding: 20px;
            border-left: 3px solid var(--primary);
        }
        .info-item h3 {
            font-size: 16px;
            font-weight: 600;
            color: var(--white);
            margin-bottom: 10px;
            display: flex;
            align-items: center;
        }
        .info-item h3 .icon {
            margin-right: 8px;
            color: var(--secondary);
        }
        .info-item p {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.7);
        }
        .highlight {
            color: var(--primary);
            font-weight: 600;
        }
        .rating {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
        .rating-stars {
            display: flex;
            margin-right: 15px;
        }
        .star {
            color: var(--primary);
            font-size: 20px;
            margin-right: 3px;
        }
        .star.empty {
            opacity: 0.3;
        }
        .timeline {
            position: relative;
            padding-left: 30px;
        }
        .timeline::before {
            content: '';
            position: absolute;
            left: 10px;
            top: 0;
            bottom: 0;
            width: 2px;
            background-color: var(--primary);
            opacity: 0.3;
        }
        .timeline-item {
            position: relative;
            padding-bottom: 20px;
        }
        .timeline-item:last-child {
            padding-bottom: 0;
        }
        .timeline-item::before {
            content: '';
            position: absolute;
            left: -30px;
            top: 5px;
            width: 12px;
            height: 12px;
            border-radius: 50%;
            background-color: var(--primary);
            border: 2px solid var(--dark);
        }
        .timeline-time {
            font-size: 14px;
            font-weight: 600;
            color: var(--secondary);
            margin-bottom: 5px;
        }
        .timeline-content {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.8);
        }
        .visualization {
            background-color: rgba(0, 0, 0, 0.2);
            border-radius: 8px;
            padding: 20px;
            margin-top: 20px;
            position: relative;
            overflow: hidden;
        }
        .visualization::after {
            content: '';
            position: absolute;
            top: 0;
            right: 0;
            width: 100px;
            height: 100px;
            background: radial-gradient(circle, rgba(56, 189, 248, 0.1) 0%, rgba(0, 0, 0, 0) 70%);
            z-index: 0;
        }
        .visualization-title {
            font-size: 16px;
            font-weight: 600;
            color: var(--white);
            margin-bottom: 15px;
            position: relative;
            z-index: 1;
        }
        .visualization-content {
            position: relative;
            z-index: 1;
        }
        .bar-chart {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }
        .bar-item {
            display: flex;
            align-items: center;
        }
        .bar-label {
            width: 120px;
            font-size: 14px;
            color: var(--white);
        }
        .bar-container {
            flex: 1;
            height: 20px;
            background-color: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            overflow: hidden;
        }
        .bar {
            height: 100%;
            background: linear-gradient(90deg, var(--primary), var(--secondary));
            border-radius: 10px;
        }
        .bar-value {
            margin-left: 10px;
            font-size: 14px;
            color: var(--primary);
        }
        .tag {
            display: inline-block;
            padding: 4px 10px;
            background-color: rgba(41, 255, 198, 0.1);
            color: var(--primary);
            border-radius: 4px;
            font-size: 12px;
            margin-right: 8px;
            margin-bottom: 8px;
        }
        .tag.warning {
            background-color: rgba(255, 165, 0, 0.1);
            color: #ffa500;
        }
        .tag.info {
            background-color: rgba(56, 189, 248, 0.1);
            color: var(--secondary);
        }
        .person-card {
            display: flex;
            background-color: rgba(0, 0, 0, 0.2);
            border-radius: 8px;
            overflow: hidden;
            margin-bottom: 20px;
        }
        .person-rank {
            width: 60px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 28px;
            font-weight: 700;
            color: var(--primary);
            background-color: rgba(41, 255, 198, 0.1);
        }
        .person-details {
            flex: 1;
            padding: 20px;
        }
        .person-name {
            font-size: 18px;
            font-weight: 600;
            color: var(--white);
            margin-bottom: 5px;
        }
        .person-meta {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin-bottom: 15px;
        }
        .person-meta-item {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.7);
        }
        .person-meta-item strong {
            color: var(--white);
        }
        .person-highlights {
            margin-top: 15px;
        }
        .person-highlight {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.8);
            margin-bottom: 10px;
            padding-left: 15px;
            position: relative;
        }
        .person-highlight::before {
            content: '•';
            position: absolute;
            left: 0;
            color: var(--secondary);
        }
        .divider {
            height: 1px;
            background-color: rgba(255, 255, 255, 0.1);
            margin: 25px 0;
        }
        .btn {
            display: inline-block;
            padding: 10px 20px;
            background-color: var(--primary);
            color: var(--dark);
            border-radius: 6px;
            font-weight: 600;
            text-decoration: none;
            transition: all 0.3s ease;
            border: none;
            cursor: pointer;
        }
        .btn:hover {
            background-color: var(--white);
            transform: translateY(-2px);
        }
        .btn-secondary {
            background-color: transparent;
            color: var(--primary);
            border: 1px solid var(--primary);
        }
        .btn-secondary:hover {
            background-color: rgba(41, 255, 198, 0.1);
        }
        @media (max-width: 768px) {
            .card-header, .card-body {
                padding: 20px;
            }
            .info-grid {
                grid-template-columns: 1fr;
            }
            .person-card {
                flex-direction: column;
            }
            .person-rank {
                width: 100%;
                height: 50px;
            }
        }
    </style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
</head>
<body>
    <div class="container">
        <!-- 爆款潜力报告卡片 -->
        <div class="report-card">
            <div class="card-header">
                <h2>爆款视频潜力分析</h2>
                <div class="subtitle">综合预测报告</div>
            </div>
            <div class="card-body">
                <div class="section">
                    <div class="rating">
                        <div class="rating-stars">
                            <span class="star"><i class="fas fa-star"></i></span>
                            <span class="star"><i class="fas fa-star"></i></span>
                            <span class="star"><i class="fas fa-star"></i></span>
                            <span class="star"><i class="fas fa-star"></i></span>
                            <span class="star"><i class="fas fa-star-half-alt"></i></span>
                        </div>
                        <div class="highlight">4.5星 - 大爆款潜力</div>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-thumbs-up"></i></span>
                        优势分析
                    </h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-bolt"></i></span>话题自带流量</h3>
                            <p>顶级富豪、财富故事、商业八卦等元素天然吸引观众注意力</p>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-database"></i></span>信息密度高</h3>
                            <p>内容丰富，数据详实，为观众提供大量有价值的信息</p>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-book-open"></i></span>故事性强</h3>
                            <p>黄峥的"神秘"与导师关系、丁磊的"快乐哲学"等故事元素增加趣味性</p>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-layer-group"></i></span>结构清晰</h3>
                            <p>排行榜形式易于观众理解和跟随，信息组织有序</p>
                        </div>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-chart-line"></i></span>
                        预估爆发期
                    </h3>
                    <div class="visualization">
                        <div class="visualization-title">1-3天内快速传播</div>
                        <div class="visualization-content">
                            <p>话题具有即时吸引力，容易在初期引发快速传播和讨论</p>
                            <div class="bar-chart">
                                <div class="bar-item">
                                    <div class="bar-label">初期爆发</div>
                                    <div class="bar-container">
                                        <div class="bar" style="width: 90%;"></div>
                                    </div>
                                    <div class="bar-value">90%</div>
                                </div>
                                <div class="bar-item">
                                    <div class="bar-label">中期传播</div>
                                    <div class="bar-container">
                                        <div class="bar" style="width: 65%;"></div>
                                    </div>
                                    <div class="bar-value">65%</div>
                                </div>
                                <div class="bar-item">
                                    <div class="bar-label">长期热度</div>
                                    <div class="bar-container">
                                        <div class="bar" style="width: 30%;"></div>
                                    </div>
                                    <div class="bar-value">30%</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 可提升空间报告卡片 -->
        <div class="report-card">
            <div class="card-header">
                <h2>优化建议</h2>
                <div class="subtitle">可提升空间分析</div>
            </div>
            <div class="card-body">
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-lightbulb"></i></span>
                        改进方向
                    </h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-eye"></i></span>视觉呈现</h3>
                            <p>可以更具冲击力和记忆点，当前部分素材(如人物照片、图表)相对常规</p>
                            <div class="tags" style="margin-top: 10px;">
                                <span class="tag">视觉锤设计</span>
                                <span class="tag">动态效果</span>
                            </div>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-cut"></i></span>剪辑节奏</h3>
                            <p>部分环节可优化，信息密度过高可能导致观众疲劳</p>
                            <div class="tags" style="margin-top: 10px;">
                                <span class="tag">节奏把控</span>
                                <span class="tag">呼吸感</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-exclamation-triangle"></i></span>
                        风险提示
                    </h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-filter"></i></span>平台审核</h3>
                            <p>涉及顶级富豪及财富话题，需注意平台审核尺度</p>
                            <div class="tags" style="margin-top: 10px;">
                                <span class="tag warning">限流关键词</span>
                            </div>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-brain"></i></span>记忆点</h3>
                            <p>画风记忆点可加强，设计更独特的视觉元素</p>
                            <div class="tags" style="margin-top: 10px;">
                                <span class="tag">品牌识别</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 视频修改建议报告卡片 -->
        <div class="report-card">
            <div class="card-header">
                <h2>视频修改建议</h2>
                <div class="subtitle">精细到秒的优化方案</div>
            </div>
            <div class="card-body">
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-bullseye"></i></span>
                        核心理念
                    </h3>
                    <div class="tags">
                        <span class="tag highlight">提升信息获取效率</span>
                        <span class="tag highlight">增强视觉冲击力</span>
                        <span class="tag highlight">制造记忆点</span>
                        <span class="tag highlight">引导情绪共鸣</span>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-play"></i></span>
                        视频内容结构与开场 (00:00 - 00:13)
                    </h3>
                    <div class="timeline">
                        <div class="timeline-item">
                            <div class="timeline-time">00:00-00:01</div>
                            <div class="timeline-content">
                                <p><strong>修改:</strong> 直接切入核心问题，用强有力的语气提问"中国现在最有钱的人，你以为还是马爸爸？" 配合醒目的文字特效。去掉开场整理衣服的动作，保证第一秒就抓住眼球。</p>
                            </div>
                        </div>
                        <div class="timeline-item">
                            <div class="timeline-time">00:01-00:04</div>
                            <div class="timeline-content">
                                <p><strong>修改:</strong> 紧接问题，快速反驳："那你就out了！" 营造悬念和认知反差。语速稍快，增加紧迫感。</p>
                            </div>
                        </div>
                        <div class="timeline-item">
                            <div class="timeline-time">00:04-00:12</div>
                            <div class="timeline-content">
                                <p><strong>修改:</strong> 语速加快，信息提炼。将"据福布斯统计..."改为更口语化的"福布斯最新数据，去年中国富豪经历大洗牌！"</p>
                                <p><strong>视觉:</strong> 放弃静态头像平铺。改为：</p>
                                <ul style="margin-left: 20px; margin-top: 10px;">
                                    <li>00:05-00:07: 用快速闪过的多位富豪头像(模糊处理或剪影)配合"大洗牌"文字动画(如碎裂效果)</li>
                                    <li>00:08-00:12: 用动态图表(如下降曲线)或夸张的视觉符号(如缩水的钱袋)展示"财富缩水5000多亿"，配上"史上最大跌幅！"的醒目文字和音效(如"咚"的一声)</li>
                                </ul>
                            </div>
                        </div>
                        <div class="timeline-item">
                            <div class="timeline-time">00:13-00:16</div>
                            <div class="timeline-content">
                                <p><strong>修改:</strong> 语气转为好奇和引导，"所以今天，我们就来看看最新的内地十大富豪，你还认识几个？" 画面可以加入"TOP 10"的预告字样。</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-list-ol"></i></span>
                        榜单呈现与人物介绍 (00:17 - 结尾)
                    </h3>
                    <div class="info-item" style="margin-bottom: 20px;">
                        <h3><span class="icon"><i class="fas fa-cog"></i></span>通用修改建议</h3>
                        <ul style="margin-left: 20px; margin-top: 10px;">
                            <li><strong>转场:</strong> 每个富豪介绍之间使用更快速、更有冲击力的转场效果(如快闪、数字滚动特效)</li>
                            <li><strong>人物信息卡:</strong> 统一设计风格，包含：排名数字(放大突出)、姓名、年龄、财富值(用亿做单位，简化数字)、核心企业Logo、人物照片。信息卡出现时配合音效</li>
                            <li><strong>BGM:</strong> 根据人物风格或故事调性切换BGM。例如，介绍黄峥/王卫时用略带神秘感的BGM，介绍丁磊时用更活泼/轻松的BGM</li>
                            <li><strong>画面:</strong> 减少纯粹的口播画面时长。口播时，多利用分屏、画中画，一边是主播，另一边展示相关图片/视频/数据</li>
                            <li><strong>字幕:</strong> 关键信息(如财富数字、绰号、核心观点)用不同颜色或放大效果突出</li>
                        </ul>
                    </div>
                    <div class="person-card">
                        <div class="person-rank">10</div>
                        <div class="person-details">
                            <div class="person-name">秦英林</div>
                            <div class="person-meta">
                                <div class="person-meta-item"><strong>年龄:</strong> 58岁</div>
                                <div class="person-meta-item"><strong>财富:</strong> 1850亿</div>
                                <div class="person-meta-item"><strong>企业:</strong> 牧原股份</div>
                            </div>
                            <div class="person-highlights">
                                <div class="person-highlight">
                                    <strong>"江湖人称'猪王'"</strong> - 视觉: 除了照片，可以快速闪现一个戴皇冠的猪的卡通形象，增加趣味记忆点
                                </div>
                                <div class="person-highlight">
                                    展示牧原股份Logo或现代化养猪场画面，而非单纯公司名文字
                                </div>
                                <div class="person-highlight">
                                    介绍夫妻创业和专业知识时，视觉: 可以用简单的动画图标展示"专业知识"(如显微镜、DNA链)应用到"养猪"(猪图标)的过程
                                </div>
                                <div class="person-highlight">
                                    将静态图片升级为更生动或有趣的短视频/GIF，如猪在舒适环境活动、自动化喂食等。使用快速剪辑，每项"福利"对应一个画面
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="person-card">
                        <div class="person-rank">9</div>
                        <div class="person-details">
                            <div class="person-name">黄峥</div>
                            <div class="person-meta">
                                <div class="person-meta-item"><strong>年龄:</strong> 43岁</div>
                                <div class="person-meta-item"><strong>财富:</strong> 2100亿</div>
                                <div class="person-meta-item"><strong>企业:</strong> 拼多多</div>
                            </div>
                            <div class="person-highlights">
                                <div class="person-highlight">
                                    在黄峥照片上打上问号，或做成"通缉令"样式，强调其低调
                                </div>
                                <div class="person-highlight">
                                    用纽交所的画面，但P掉敲钟人，或者用一个"X"划掉敲钟画面
                                </div>
                                <div class="person-highlight">
                                    用时间轴+财富曲线图展示其财富暴涨、捐赠、卸任CEO的时间点，制造戏剧性
                                </div>
                                <div class="person-highlight">
                                    MSN对话框做成动态弹出效果，丁磊头像闪烁，模拟真实聊天场景
                                </div>
                                <div class="person-highlight">
                                    重点突出巴菲特午餐照片，可以放大照片，并加上"与股神共进午餐"的标签
                                </div>
                                <div class="person-highlight">
                                    将"低调搞钱团"的合照做成更神秘、更有力量感的视觉设计，并列出成员及其公司Logo
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="person-card">
                        <div class="person-rank">8</div>
                        <div class="person-details">
                            <div class="person-name">丁磊</div>
                            <div class="person-meta">
                                <div class="person-meta-item"><strong>年龄:</strong> 52岁</div>
                                <div class="person-meta-item"><strong>财富:</strong> 2250亿</div>
                                <div class="person-meta-item"><strong>企业:</strong> 网易</div>
                            </div>
                            <div class="person-highlights">
                                <div class="person-highlight">
                                    <strong>"快乐男孩"</strong> - 配合丁磊大笑的照片或趣味表情包
                                </div>
                                <div class="person-highlight">
                                    使用具有年代感的电脑/网页界面元素，营造复古氛围
                                </div>
                                <div class="person-highlight">
                                    提及雷军、马化腾等，可以快速闪现他们的头像，并用连线表示社交关系网
                                </div>
                                <div class="person-highlight">
                                    快速剪辑展示网易爆款游戏的精彩画面或经典角色。对比腾讯时，可以使用VS图标
                                </div>
                                <div class="person-highlight">
                                    网易云音乐部分，用其标志性的评论区截图或播放界面；DJ部分使用更动感的画面和剪辑；变胖过程的对比照片是亮点，确保清晰展示
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-comments"></i></span>
                        互动引导设计
                    </h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-question"></i></span>视频中段</h3>
                            <p>可以插入提问，"这几位大佬，谁的故事最让你意外？"引导评论</p>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-question-circle"></i></span>视频结尾</h3>
                            <p>除了总结，可以再次提问，"你觉得未来谁的排名会上升/下降？评论区聊聊！" 或者 "还想了解哪位大佬的故事？" 为下期内容做铺垫</p>
                        </div>
                    </div>
                </div>
                <div class="section">
                    <h3 class="section-title">
                        <span class="icon"><i class="fas fa-microphone"></i></span>
                        话术设计
                    </h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-comment-dots"></i></span>口语化</h3>
                            <p>保持自然、流畅的口语表达</p>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-tachometer-alt"></i></span>节奏感</h3>
                            <p>重要信息(如排名、财富)语速稍慢，强调；故事性内容可以更生动、有起伏</p>
                        </div>
                        <div class="info-item">
                            <h3><span class="icon"><i class="fas fa-lightbulb"></i></span>观点性</h3>
                            <p>在不失客观的前提下，可以加入少量个人化的评价或感叹，增加亲和力，如"这操作真是绝了"、"简直是开挂的人生"</p>
                        </div>
                    </div>
                </div>
                <div class="divider"></div>
                <div style="text-align: center; padding: 20px 0;">
                    <p style="margin-bottom: 20px; font-size: 16px; color: var(--primary);">放手去改吧！你的内容基础非常好，有潜力触达更广泛的观众。</p>
                    <p style="margin-bottom: 20px;">每一次精心的打磨，都是让作品更加闪耀的过程，相信你的视频一定能获得更多人的喜爱！</p>
                    <button class="btn">导出修改建议</button>
                    <button class="btn btn-secondary" style="margin-left: 15px;">分享报告</button>
                </div>
            </div>
        </div>
    </div>
</body>
</html>

`,
}

func (t TrackerService) ListPublicRecords(ctx context.Context, params *trackerpb.ListRecordsParams) (*trackerpb.RecordList, error) {

	return nil, nil
}

func (t TrackerService) GetPublicRecord(ctx context.Context, params *trackerpb.GetRecordParams) (*trackerpb.Record, error) {

	fr, ok := publicRecords[params.Id]
	if !ok {
		return nil, nil
	}

	return &trackerpb.Record{
		Id:             params.Id,
		Category:       "improveSuggestion",
		FormatedResult: fr,
	}, nil
}

func (t TrackerService) ListRecords(ctx context.Context, params *trackerpb.ListRecordsParams) (*trackerpb.RecordList, error) {

	return &trackerpb.RecordList{}, nil
}

type PromptArgs struct {
	AuthorProfile   string
	PersonalProfile string
	Platform        string
}

func (t TrackerService) GetPromptById(ctx context.Context, promptId string, args PromptArgs) (string, error) {
	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return "", err
	}

	promptText := settings.GetPrompt()[promptId]
	if promptText == "" {
		return "", errors.BadRequest("invalidPromptId", "")
	}

	promptText = strings.ReplaceAll(promptText, "__AUTHOR_PROFILE__", args.AuthorProfile)
	promptText = strings.ReplaceAll(promptText, "__PERSONAL_PROFILE__", args.PersonalProfile)

	return promptText, nil
}

func (t TrackerService) GetPromptV2(ctx context.Context, prompt *trackerpb.Prompt, args PromptArgs) (string, error) {

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return "", err
	}

	var promptText string
	if prompt.GetId() != "" {
		promptText = settings.GetPrompt()[prompt.GetId()]
		if promptText == "" {
			return "", errors.BadRequest("invalidPromptId", "")
		}

		//
		promptText = strings.ReplaceAll(promptText, "__AUTHOR_PROFILE__", args.AuthorProfile)
		promptText = strings.ReplaceAll(promptText, "__PERSONAL_PROFILE__", args.PersonalProfile)

	} else {
		promptText = prompt.GetContent()

		if promptText == "" {
			return "", errors.BadRequest("invalidPrompt", "")
		}

		promptText += settings.GetPrompt()["default"]

	}

	promptText = strings.ReplaceAll(promptText, "__LOCALE__", helper.OrString(prompt.Locale, "zh"))

	//localePrompt := fmt.Sprintf("\n###当前用户的 Locale code 为: %s, 请以对应的语言输出", prompt.Locale)

	return promptText, nil
}

func (t TrackerService) GetPrompt(ctx context.Context, profile, category string) (string, error) {
	promptMap := map[string]string{

		"improveSuggestion": `

你是一个社交媒体内容分析专家，请基于以下维度分析用户提供的内容（视频/图文），预测其成为爆款的概率并给出优化建议：

1.⚠️博主账号信息：
%s

2. 关键元素提取 （#以表格形式展示）
   - [ ] 情绪价值（共鸣/实用/新奇/治愈）
   - [ ] 信息密度（重点突出程度）
   - [ ] 视觉锤（记忆点设计）
3. 互动潜力评估（#以表格形式展示）
   - [ ] 可模仿性（是否引发UGC跟拍/二创）
   - [ ] 争议性（正向讨论空间）
   - [ ] 收藏动机（长效价值）
4. 爆款对标分析（#以表格形式展示）
   - [ ] 同赛道Top100笔记关键词匹配度
   - [ ] 近期热点话题关联性
   - [ ] 平台流量周期（如节假日/大促节点）
5. 综合预测报告 #先说结论：这个视频在社交平台发布后是否可以成为爆款？（#以表格形式展示）
（分为：极低概率爆款的潜力，小概率爆款的潜力，中低概率爆款的潜力，50%概率爆款的潜力，中高概率爆款的潜力，高概率爆款的潜力，极概率高爆款的潜力。#需要基于账号信息一起评估，用最保守的方式评估）
   - 爆款潜力：⭐️⭐️⭐️⭐️⭐️（0-5星）
   - 预估爆发期：1-3天 / 3-7天 / 长尾流量型


#先说结论，即先说综合预测报告
#如果要成为爆款视频，可以基于这个视频，给我详细的哪一秒需要修改 （精细到秒）；以表格形式输出，怎么修改角度可以包含：
视频内容结构、
开场前3秒表现、
画面构图和质量、
剪辑节奏、
BGM选择、
字幕设计、
互动引导设计、
话术设计...等等；手把手教我修改, 越详细越好；这个部分以表格形式输出；
#最后加一句鼓励的话语目的是提供情绪价值
`,
		"limitDetection": `
#先说结论：预测限流分析：这个视频在国内的社交媒体平台，分别是小红书、抖音、微信视频号、快手、微博，在这几个平台发布后是否会限流；分别系统生成限流风险评分：中高低风险；（#以表格形式展示）

#预测我这个视频在国内各大社交媒体平台是否会限流，系统限流预测分析的角度包含以下要素（#以表格形式展示）：
一、多模态内容分析:
-文本层:检测敏感词(政治/医疗/营销术语)、标题党特征、联系方式、违禁商品...等等；
-视觉层:识别暴露着装、血腥暴力、侵权水印、低画质内容(分辨率<720p)、商业推广过度、违规内容...等等；
-音频内容分析（适用于视频）

二、平台规则模拟（#以表格形式展示）：
-与国内各大社交媒体平台内容政策的符合度 ，包含各大社交媒体的社区规范以及违禁词进行搜索，结合视频内容给出对应的限流违规条款 

⚠️三、结合博主上传的账号信息和账号权重评以及内容垂直度一起评估（#以表格形式展示），
以下是博主的账号信息：
%s

#输出要求包含：各大社交媒体的限流风险（高中低风险）+ 风险项清单（按紧急度排序）+ 修改建议（具体替换词建议/裁剪时间点）（#以表格形式展示）
`,

		"trafficPrediction": `
#我这个视频要上传到选中要发布的社交平台（小红书），目标：结合博主的账号权重，帮我预测这个视频发布后的爆款潜力（分为：极地概率爆款的潜力，小概率爆款的潜力，中低概率爆款的潜力，50%概率爆款的潜力，中高概率爆款的潜力，高概率爆款的潜力，极概率高爆款的潜力。#需要基于账号信息一起评估；（#以表格形式展示）
#结合博主的账号权重，使用机率加權蒙特卡羅模擬的方式来进行预测这个视频发布后30天内的数据表现，主要目标函数是曝光量；其次是点赞，留言，收藏，转发）；（#以表格形式展示）
#需要精准到概率百分比，结合使用统计学置信度分析；用最严格保守的预测管理用户预期；（#以表格形式展示）
#结合博主的账号权重帮我预测，

⚠️已上传的博主账号信息：
%s
`,
	}

	prompt := promptMap[category]

	if prompt == "" {
		return "", errors.BadRequest("invalidPromptId", "")
	}

	//	// 博主信息
	//	profiles, err := t.Data.EntClient.UserProfile.Query().
	//		Where(userprofile.UserID(userId)).
	//		Where(userprofile.Status("normal")).
	//		All(ctx)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	if len(profiles) == 0 {
	//		return fmt.Sprintf(prompt, ""), nil
	//	}
	//
	//	var profileInfo = fmt.Sprintf(`
	//%s
	//小红书号:%s
	//%s
	//%s关注
	//%s粉丝
	//%s获赞与收藏
	//%s条笔记或视频
	//`,
	//		profiles[0].Content["username"],
	//		profiles[0].Content["id"],
	//		profiles[0].Content["ipAddress"],
	//		profiles[0].Content["followingCount"],
	//		profiles[0].Content["followerCount"],
	//		profiles[0].Content["likedCount"],
	//		profiles[0].Content["noteCount"],
	//	)

	return fmt.Sprintf(prompt, profile), nil
}

func (t TrackerService) ResponseString(resp *genai.GenerateContentResponse) string {
	var b strings.Builder
	for i, cand := range resp.Candidates {
		if len(resp.Candidates) > 1 {
			fmt.Fprintf(&b, "%d:", i+1)
		}
		b.WriteString(t.contentString(cand.Content))
	}
	return b.String()
}

func (t TrackerService) contentString(c *genai.Content) string {
	var b strings.Builder
	if c == nil || c.Parts == nil {
		return ""
	}
	for i, part := range c.Parts {
		if i > 0 {
			fmt.Fprintf(&b, ";")
		}
		fmt.Fprintf(&b, "%v", part)
	}
	return b.String()
}

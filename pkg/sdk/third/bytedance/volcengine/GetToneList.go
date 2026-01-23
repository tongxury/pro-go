package volcengine

import (
	"context"
)

func (t *Client) GetToneList(ctx context.Context, params GetToneListParams) (*Tones, error) {
	return &toneConfig, nil
}

type GetToneListParams struct {
}

type Tone struct {
	Id          string `json:"Id"`
	Title       string `json:"Title"`
	DownloadUrl string `json:"DownloadUrl"`
	Description string `json:"Description"`
}

type Tones []Tone

var toneConfig = Tones{
	{
		Id:          "0",
		Title:       "抖音 IP 小姐姐",
		Description: "机械女声，适用于通用场景",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV001_streaming_v2.mp3",
	},
	{
		Id:          "1",
		Title:       "抖音 IP 小哥哥",
		Description: "机械男声，适用于通用场景",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV002_streaming_v2.mp3",
	},
	{
		Id:          "2",
		Title:       "成熟女声",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV009_DPE_streaming_v2.mp3",
	},
	{
		Id:          "3",
		Title:       "青年女声",
		Description: "女声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV007_streaming_v2.mp3",
	},
	{
		Id:          "4",
		Title:       "稳重大叔",
		Description: "男声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV006_streaming_v2.mp3",
	},
	{
		Id:          "5",
		Title:       "青年男声",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV008_DPE_streaming_v2.mp3",
	},
	{
		Id:          "6",
		Title:       "新闻女声",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV011_streaming_v2.mp3",
	},
	{
		Id:          "7",
		Title:       "可爱少女",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV005_streaming_v2.mp3",
	},
	{
		Id:          "8",
		Title:       "新闻男声",
		Description: "男声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV012_streaming_v2.mp3",
	},
	{
		Id:          "9",
		Title:       "活力青年",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV056_streaming_v2.mp3",
	},
	{
		Id:          "10",
		Title:       "中英男声",
		Description: "男声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV033_ParaTaco_v2.mp3",
	},
	{
		Id:          "11",
		Title:       "东北老铁",
		Description: "男声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV021_streaming_v2.mp3",
	},
	{
		Id:          "12",
		Title:       "西安掌柜",
		Description: "女声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV210_streaming_v2.mp3",
	},
	{
		Id:          "13",
		Title:       "港剧男神",
		Description: "男声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV026_streaming_v2.mp3",
	},
	{
		Id:          "14",
		Title:       "甜美台妹",
		Description: "女声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV025_streaming_v2.mp3",
	},
	{
		Id:          "15",
		Title:       "相声演员",
		Description: "男声，推荐，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV212_streaming_v2.mp3",
	},
	{
		Id:          "16",
		Title:       "重庆小伙",
		Description: "男声，推荐，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV019_streaming_v2.mp3",
	},
	{
		Id:          "17",
		Title:       "二次元萝莉",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV064_streaming_v2.mp3",
	},
	{
		Id:          "18",
		Title:       "海绵宝宝",
		Description: "男声，推荐",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV063_streaming_v2.mp3",
	},
	{
		Id:          "19",
		Title:       "萌娃童声",
		Description: "男声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV051_streaming_v2.mp3",
	},
	{
		Id:          "20",
		Title:       "说书大叔",
		Description: "男声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV110_streaming_v2.mp3",
	},
	{
		Id:          "21",
		Title:       "阳光青年",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV123_streaming_v2.mp3",
	},
	{
		Id:          "22",
		Title:       "憨厚青年",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV120_streaming_v2.mp3",
	},
	{
		Id:          "23",
		Title:       "散漫赘婿",
		Description: "男声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV119_streaming_v2.mp3",
	},
	{
		Id:          "24",
		Title:       "霸气青叔",
		Description: "男声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV107_streaming_v2.mp3",
	},
	{
		Id:          "25",
		Title:       "质朴青年",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV100_streaming_v2.mp3",
	},
	{
		Id:          "26",
		Title:       "儒雅青年",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV102_streaming_v2.mp3",
	},
	{
		Id:          "27",
		Title:       "开朗青年",
		Description: "男声，青年",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV004_streaming_v2.mp3",
	},
	{
		Id:          "28",
		Title:       "温和少御",
		Description: "女声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV114_streaming_v2.mp3",
	},
	{
		Id:          "29",
		Title:       "平缓少御",
		Description: "女声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV113_streaming_v2.mp3",
	},
	{
		Id:          "30",
		Title:       "甜美女声",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV405_streaming_v2.mp3",
	},
	{
		Id:          "32",
		Title:       "活泼幼教",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV057_ParaTaco_streaming_v2.mp3",
	},
	{
		Id:          "33",
		Title:       "活泼女声",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV005_ParaTaco_streaming_v2.mp3",
	},
	{
		Id:          "34",
		Title:       "亲切女声",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV007_ParaTaco_streaming_v2.mp3",
	},
	{
		Id:          "35",
		Title:       "知性女声",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV009_DPE_ParaTaco_streaming_v2.mp3",
	},
	{
		Id:          "36",
		Title:       "知性男声",
		Description: "男声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV008_DPE_ParaTaco_streaming_v2.mp3",
	},
	{
		Id:          "37",
		Title:       "灿灿",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV700_streaming_v2.mp3",
	},
	{
		Id:          "38",
		Title:       "阳光男声",
		Description: "男声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV419_streaming_v2.mp3",
	},
	{
		Id:          "39",
		Title:       "天才童声",
		Description: "男声，虚拟人，推荐",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV061_v2.mp3",
	},
	{
		Id:          "40",
		Title:       "超自然 - 梓梓",
		Description: "女声，虚拟人",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV406_streaming_v2.mp3",
	},
	{
		Id:          "45",
		Title:       "动漫小新",
		Description: "男声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV050_streaming_v2.mp3",
	},
	{
		Id:          "46",
		Title:       "台普男声",
		Description: "男声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV227_streaming_v2.mp3",
	},
	{
		Id:          "47",
		Title:       "广西表哥",
		Description: "男声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV213_streaming_v2.mp3",
	},
	{
		Id:          "48",
		Title:       "温柔小哥",
		Description: "男声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV033_streaming_v2.mp3",
	},
	{
		Id:          "49",
		Title:       "影视解说小帅",
		Description: "男声，推荐，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV411_streaming_v2.mp3",
	},
	{
		Id:          "50",
		Title:       "活力解说男",
		Description: "男声，推荐，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV410_streaming_v2.mp3",
	},
	{
		Id:          "51",
		Title:       "译制片男声",
		Description: "男声，推荐，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV408_streaming_v2.mp3",
	},
	{
		Id:          "52",
		Title:       "擎苍",
		Description: "男声，推荐，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV701_streaming_v2.mp3",
	},
	{
		Id:          "53",
		Title:       "智慧老者",
		Description: "男声，推荐，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV158_streaming_v2.mp3",
	},
	{
		Id:          "54",
		Title:       "东北丫头",
		Description: "女声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV020_streaming_v2.mp3",
	},
	{
		Id:          "55",
		Title:       "长沙靓女",
		Description: "女声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV216_streaming_v2.mp3",
	},
	{
		Id:          "56",
		Title:       "沪上阿姐",
		Description: "女声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV217_streaming_v2.mp3",
	},
	{
		Id:          "57",
		Title:       "促销女声",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV402_streaming_v2.mp3",
	},
	{
		Id:          "58",
		Title:       "湖南妹坨",
		Description: "女声，方言",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV226_streaming_v2.mp3",
	},
	{
		Id:          "59",
		Title:       "天才少女",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV421_streaming_v2.mp3",
	},
	{
		Id:          "60",
		Title:       "鸡汤女声",
		Description: "女声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV403_streaming_v2.mp3",
	},
	{
		Id:          "61",
		Title:       "影视解说小美",
		Description: "女声，小说",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV412_streaming_v2.mp3",
	},
	{
		Id:          "62",
		Title:       "直播一姐",
		Description: "女声",
		DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV418_streaming_v2.mp3",
	},
}

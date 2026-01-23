package configs

import aiagentpb "store/api/aiagent"

var ModelFlash = "gemini-2.5-flash"

//var ModelFlash = "gemini-2.5-flash"

var pro = "gemini-2.5-pro-preview-05-06"

//var pro = "gemini-2.5-pro"

type Cost struct {
	Value int64
}
type QuotaConfig struct {
	Items []*aiagentpb.PromptSetting
}

var config = &QuotaConfig{
	Items: []*aiagentpb.PromptSetting{
		{
			//Scene:    "preAnalysis", // 爆款预测
			PromptId: "preAnalysis", // 视频预测
			Cost:     40,
			//Model:    flash,
			Model: ModelFlash,
		},
		{
			//Scene:    "preAnalysis",       // 爆款预测
			PromptId: "preAnalysisImages", // 视频预测
			Cost:     20,
			//Model:    flash,
			Model: ModelFlash,
		},
		{
			//Scene:    "limitAnalysis", // 限流预测
			PromptId: "limitAnalysis", // 视频预测
			Cost:     40,
			Model:    ModelFlash,
		},
		{
			//Scene:    "limitAnalysis",       // 限流预测
			PromptId: "limitAnalysisImages", // 视频预测
			Cost:     20,
			Model:    ModelFlash,
		},
		{
			//Scene:    "coverAnalysis",       // 封面预测
			PromptId: "coverAnalysisImages", // 封面预测
			Cost:     30,
			//Model:    pro,
			Model:    ModelFlash,
			MaxFiles: 3,
			//gemini-2.5-pro-preview-05-06
		},
		{
			//Scene:    "analysis", // 爆款分析
			PromptId: "analysis", // 视频分析
			Cost:     40,
			//Model:    pro,
			Model: ModelFlash,
		},
		{
			//Scene:    "analysis",       // 爆款分析
			PromptId: "analysisImages", // 视频分析
			Cost:     20,
			//Model:    pro,
			Model: ModelFlash,
		},
		{
			//Scene:    "duplicateScript", // 脚本复刻
			PromptId: "duplicateScript", // 脚本复刻
			Cost:     40,
			Model:    ModelFlash,
		},
		{
			//Scene:    "duplicateScript", // 脚本复刻
			PromptId: "duplicateScriptImages", // 脚本复刻
			Cost:     20,
			Model:    ModelFlash,
		},
		// 脚本优化
		{
			PromptId: "scriptOptimization",
			Cost:     20,
			Model:    ModelFlash,
		},
		// 文本提前
		{
			PromptId: "contentExtraction",
			Cost:     20,
			Model:    ModelFlash,
		},
		// 文本生成
		{
			PromptId: "contentGeneration",
			Cost:     20,
			Model:    ModelFlash,
		},
		// 账号诊断
		{
			PromptId: "accountAnalysis",
			Cost:     0,
			Model:    ModelFlash,
		},
	},
}

var defaultCost = &Cost{
	Value: 10,
}

func GetPromptSettings() []*aiagentpb.PromptSetting {
	return config.Items
}

func GetCost(scene, promptId string) int64 {

	for _, x := range config.Items {
		if x.PromptId == promptId {
			return x.Cost
		}
	}

	return 20
}

func GetModel(scene, promptId string) string {

	for _, x := range config.Items {
		if x.PromptId == promptId {
			return x.Model
		}
	}

	return ModelFlash
}

func GetModalPro() string {
	return pro
}

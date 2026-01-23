package configs

import (
	paymentpb "store/api/payment"
	"store/pkg/sdk/helper"
)

var plans = []*paymentpb.Metadata_Plan{

	{
		Id:             "l1-pkg",
		Alias:          []string{"com_veogo_l1_pkg", "l1_pkg"},
		Title:          "500积分包(30天)",
		CnyAmount:      38,
		Amount:         38,
		Months:         1,
		CreditPerMonth: 500,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 16,
			},
			{
				Name:  "analysisImages",
				Value: 25,
			},
			{
				Name:  "preAnalysisImages",
				Value: 25,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 25,
			},
			{
				Name:  "analysis",
				Value: 12,
			},
			{
				Name:  "limitAnalysis",
				Value: 12,
			},
			{
				Name:  "preAnalysis",
				Value: 6,
			},
			{
				Name:  "duplicateScript",
				Value: 6,
			},
		},
	},
	{
		Id:             "l2-pkg",
		Alias:          []string{"com_veogo_l2_pkg", "l2_pkg"},
		Title:          "1500积分包(30天)",
		CnyAmount:      98,
		Amount:         98,
		Months:         1,
		CreditPerMonth: 1500,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 50,
			},
			{
				Name:  "analysisImages",
				Value: 75,
			},
			{
				Name:  "preAnalysisImages",
				Value: 75,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 75,
			},
			{
				Name:  "analysis",
				Value: 37,
			},
			{
				Name:  "limitAnalysis",
				Value: 37,
			},
			{
				Name:  "preAnalysis",
				Value: 18,
			},
			{
				Name:  "duplicateScript",
				Value: 18,
			},
		},
	},
	{
		Id:             "l3-pkg",
		Alias:          []string{"com_veogo_l3_pkg", "l3_pkg"},
		Title:          "7000积分包(30天)",
		CnyAmount:      399,
		Amount:         399,
		Months:         1,
		CreditPerMonth: 7000,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 233,
			},
			{
				Name:  "analysisImages",
				Value: 350,
			},
			{
				Name:  "preAnalysisImages",
				Value: 350,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 350,
			},
			{
				Name:  "analysis",
				Value: 175,
			},
			{
				Name:  "limitAnalysis",
				Value: 175,
			},
			{
				Name:  "preAnalysis",
				Value: 87,
			},
			{
				Name:  "duplicateScript",
				Value: 87,
			},
		},
		Extra: "赠送: 自媒体起号操作分享(价值198元)",
	},
	{
		Id:             "l4-pkg",
		Alias:          []string{"com_veogo_l4_pkg", "l4_pkg"},
		Title:          "20000积分包(30天)",
		CnyAmount:      999,
		Amount:         999,
		Months:         1,
		CreditPerMonth: 20000,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 666,
			},
			{
				Name:  "analysisImages",
				Value: 1000,
			},
			{
				Name:  "preAnalysisImages",
				Value: 1000,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 1000,
			},
			{
				Name:  "analysis",
				Value: 500,
			},
			{
				Name:  "limitAnalysis",
				Value: 500,
			},
			{
				Name:  "preAnalysis",
				Value: 250,
			},
			{
				Name:  "duplicateScript",
				Value: 250,
			},
		},
		Extra: "赠送: 1v1账号深度分析咨询60分钟(价值800元)",
	},

	{
		Id:             "l1-monthly",
		Alias:          []string{"veogo_l1_monthly", "l1_monthly"},
		Title:          "500积分包(每月)",
		CnyAmount:      38,
		Amount:         38,
		Months:         1,
		CreditPerMonth: 500,
		Mode:           "oneOff",
		//Mode:           "recurring",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 16,
			},
			{
				Name:  "analysisImages",
				Value: 25,
			},
			{
				Name:  "preAnalysisImages",
				Value: 25,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 25,
			},
			{
				Name:  "analysis",
				Value: 12,
			},
			{
				Name:  "limitAnalysis",
				Value: 12,
			},
			{
				Name:  "preAnalysis",
				Value: 6,
			},
			{
				Name:  "duplicateScript",
				Value: 6,
			},
		},
	},
	{
		Id:             "l2-monthly",
		Alias:          []string{"veogo_l2_monthly", "l2_monthly"},
		Title:          "1500积分包(每月)",
		CnyAmount:      98,
		Amount:         98,
		Months:         1,
		CreditPerMonth: 1500,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 50,
			},
			{
				Name:  "analysisImages",
				Value: 75,
			},
			{
				Name:  "preAnalysisImages",
				Value: 75,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 75,
			},
			{
				Name:  "analysis",
				Value: 37,
			},
			{
				Name:  "limitAnalysis",
				Value: 37,
			},
			{
				Name:  "preAnalysis",
				Value: 18,
			},
			{
				Name:  "duplicateScript",
				Value: 18,
			},
		},
	},
	{
		Id:    "l3-monthly",
		Alias: []string{"veogo_l3_monthly", "l3_monthly"},

		Title:          "7000积分包(每月)",
		CnyAmount:      399,
		Amount:         399,
		Months:         1,
		CreditPerMonth: 7000,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 233,
			},
			{
				Name:  "analysisImages",
				Value: 350,
			},
			{
				Name:  "preAnalysisImages",
				Value: 350,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 350,
			},
			{
				Name:  "analysis",
				Value: 175,
			},
			{
				Name:  "limitAnalysis",
				Value: 175,
			},
			{
				Name:  "preAnalysis",
				Value: 87,
			},
			{
				Name:  "duplicateScript",
				Value: 87,
			},
		},
		Extra: "赠送: 自媒体起号操作分享(价值198元)",
	},
	{
		Id:             "l4-monthly",
		Alias:          []string{"veogo_l4_monthly", "l4_monthly"},
		Title:          "2000积分包(每月)",
		CnyAmount:      999,
		Amount:         999,
		Months:         1,
		CreditPerMonth: 20000,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 666,
			},
			{
				Name:  "analysisImages",
				Value: 1000,
			},
			{
				Name:  "preAnalysisImages",
				Value: 1000,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 1000,
			},
			{
				Name:  "analysis",
				Value: 500,
			},
			{
				Name:  "limitAnalysis",
				Value: 500,
			},
			{
				Name:  "preAnalysis",
				Value: 250,
			},
			{
				Name:  "duplicateScript",
				Value: 250,
			},
		},
		Extra: "赠送: 1v1账号深度分析咨询60分钟(价值800元)",
	},
	{
		Id:             "l1-month",
		Alias:          []string{"com_veogo_l1_month", "l1_month"},
		Title:          "500积分包(30天)",
		CnyAmount:      38,
		Amount:         38,
		Months:         1,
		CreditPerMonth: 500,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 16,
			},
			{
				Name:  "analysisImages",
				Value: 25,
			},
			{
				Name:  "preAnalysisImages",
				Value: 25,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 25,
			},
			{
				Name:  "analysis",
				Value: 12,
			},
			{
				Name:  "limitAnalysis",
				Value: 12,
			},
			{
				Name:  "preAnalysis",
				Value: 6,
			},
			{
				Name:  "duplicateScript",
				Value: 6,
			},
		},
	},
	{
		Id:             "l2-month",
		Alias:          []string{"com_veogo_l2_month", "l2_month"},
		Title:          "1500积分包(30天)",
		CnyAmount:      98,
		Amount:         98,
		Months:         1,
		CreditPerMonth: 1500,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 50,
			},
			{
				Name:  "analysisImages",
				Value: 75,
			},
			{
				Name:  "preAnalysisImages",
				Value: 75,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 75,
			},
			{
				Name:  "analysis",
				Value: 37,
			},
			{
				Name:  "limitAnalysis",
				Value: 37,
			},
			{
				Name:  "preAnalysis",
				Value: 18,
			},
			{
				Name:  "duplicateScript",
				Value: 18,
			},
		},
	},
	{
		Id:             "l3-month",
		Alias:          []string{"com_veogo_l3_month", "l3_month"},
		Title:          "7000积分包(30天)",
		CnyAmount:      399,
		Amount:         399,
		Months:         1,
		CreditPerMonth: 7000,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 233,
			},
			{
				Name:  "analysisImages",
				Value: 350,
			},
			{
				Name:  "preAnalysisImages",
				Value: 350,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 350,
			},
			{
				Name:  "analysis",
				Value: 175,
			},
			{
				Name:  "limitAnalysis",
				Value: 175,
			},
			{
				Name:  "preAnalysis",
				Value: 87,
			},
			{
				Name:  "duplicateScript",
				Value: 87,
			},
		},
		Extra: "赠送: 自媒体起号操作分享(价值198元)",
	},
	{
		Id:             "l4-month",
		Alias:          []string{"com_veogo_l4_month", "l4_month"},
		Title:          "20000积分包(30天)",
		CnyAmount:      999,
		Amount:         999,
		Months:         1,
		CreditPerMonth: 20000,
		Mode:           "oneOff",
		Features: []*paymentpb.Metadata_Feature{
			{
				Name:  "coverAnalysisImages",
				Value: 666,
			},
			{
				Name:  "analysisImages",
				Value: 1000,
			},
			{
				Name:  "preAnalysisImages",
				Value: 1000,
			},
			{
				Name:  "limitAnalysisImages",
				Value: 1000,
			},
			{
				Name:  "analysis",
				Value: 500,
			},
			{
				Name:  "limitAnalysis",
				Value: 500,
			},
			{
				Name:  "preAnalysis",
				Value: 250,
			},
			{
				Name:  "duplicateScript",
				Value: 250,
			},
		},
		Extra: "赠送: 1v1账号深度分析咨询60分钟(价值800元)",
	},

	{Id: "basic-month", Title: "1个月", CnyAmount: 38, Amount: 9.9, Months: 1, CreditPerMonth: 500, Duplicated: true},
	{Id: "basic-monthly", Title: "1个月", CnyAmount: 38, Amount: 9.9, Months: 1, CreditPerMonth: 500, Duplicated: true},

	{Id: "pro-month", Title: "1个月(无限)", CnyAmount: 98, Amount: 19.9, Months: 1, CreditPerMonth: 50000, Duplicated: true},
	{Id: "pro-monthly", Title: "1个月(无限)", CnyAmount: 98, Amount: 19.9, Months: 1, CreditPerMonth: 50000, Duplicated: true},

	{Id: "basic-annual", Title: "12个月", CnyAmount: 268, Amount: 95, OriginAmount: 336, Months: 12, CreditPerMonth: 500, Duplicated: true},
	{Id: "basic-annually", Title: "12个月", CnyAmount: 268, Amount: 95, OriginAmount: 336, Months: 12, CreditPerMonth: 500, Duplicated: true},

	{Id: "pro-annual", Title: "12个月(无限)", CnyAmount: 938, Amount: 191, OriginAmount: 1176, Months: 12, CreditPerMonth: 50000, Duplicated: true},
	{Id: "pro-annually", Title: "12个月(无限)", CnyAmount: 938, Amount: 191, OriginAmount: 1176, Months: 12, CreditPerMonth: 50000, Duplicated: true},
}

func GetPlans() []*paymentpb.Metadata_Plan {
	return plans
}

func GetPlanById(planId string) *paymentpb.Metadata_Plan {
	for _, x := range plans {
		if x.Id == planId {
			return x
		}

		if helper.InSlice(planId, x.Alias) {
			return x
		}
	}

	return nil
}

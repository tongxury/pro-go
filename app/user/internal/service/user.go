package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
	aiagentpb "store/api/aiagent"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	typepb "store/api/user/types"
	"store/app/user/internal/biz"
	"store/app/user/internal/data"
	"store/app/user/internal/data/enums"
	"store/app/user/internal/data/repo"
	"store/app/user/internal/data/repo/ent/usercreditincrement"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/mathz"
	"strings"
	"time"
)

type UserService struct {
	userpb.UnimplementedUserServer
	userBiz      *biz.UserBiz
	userCredit   *biz.UserCreditBiz
	promotionBiz *biz.PromotionBiz
	data         *data.Data
}

func NewUserService(
	userBiz *biz.UserBiz,
	userCredit *biz.UserCreditBiz,
	promotionBiz *biz.PromotionBiz,
	data *data.Data,
) *UserService {
	return &UserService{
		userBiz:      userBiz,
		userCredit:   userCredit,
		promotionBiz: promotionBiz,
		data:         data,
	}
}

func (t *UserService) CreateSurveyPayment(ctx context.Context, params *userpb.CreateSurveyPaymentParams) (*userpb.CreatePaymentResult, error) {

	return &userpb.CreatePaymentResult{}, nil

}

func (t *UserService) CallbackPayment(ctx context.Context, params *userpb.CallbackPaymentParams) (*emptypb.Empty, error) {

	log.Debugw("CallbackPayment", "", "params", params)

	var err error
	if strings.HasPrefix(params.OutTradeNo, "sur_") {

		_, err = t.data.GrpcClients.AiAgentClient.UpdateSurvey(ctx, &aiagentpb.UpdateSurveyParams{
			Id: params.OutTradeNo[4:],
		})

	} else {
		err = t.data.Repos.EntClient.UserCreditIncrement.Update().
			SetStatus(enums.UserCreditStatusActive).
			SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
			Where(usercreditincrement.ID(conv.Int64(params.OutTradeNo))).
			Exec(ctx)
	}

	if err != nil {

		r, _ := http.RequestFromServerContext(ctx)

		log.Error("CallbackPayment err: ", err, "params", params, "req", r.RequestURI)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *UserService) CreatePayment(ctx context.Context, params *userpb.CreatePaymentParams) (*userpb.CreatePaymentResult, error) {
	return &userpb.CreatePaymentResult{}, nil
}

func (t *UserService) CostCredit(ctx context.Context, params *userpb.CheckCreditParams) (*userpb.CheckCreditResult, error) {

	credit, err := t.userCredit.CostCredit(ctx, params)
	if err != nil {
		log.Errorw("cost credit", "err", err, "params", params)
		return nil, err
	}

	return credit, nil
}

//func (t *UserService) IncrCredit(ctx context.Context, params *userpb.IncrCreditParams) (*emptypb.Empty, error) {
//
//	err := t.userCredit.Incr(ctx, params.UserId, params.Amount)
//	if err != nil {
//		log.Errorw("Incr credit", "err", err, "params", params)
//		return nil, err
//	}
//
//	return &emptypb.Empty{}, nil
//}

func (t *UserService) AddFeedback(ctx context.Context, params *userpb.AddFeedbackParams) (*emptypb.Empty, error) {

	userId := krathelper.RequireUserId(ctx)

	_, err := t.data.Repos.EntClient.Feedback.Create().
		SetUserID(conv.Int64(userId)).
		SetCategory(params.Category).
		SetContent(params.Content).
		Save(ctx)
	if err != nil {
		log.Errorw("AddFeedback err", err, "params", params)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *UserService) GetChatSettings(ctx context.Context, params *userpb.GetChatSettingsParams) (*userpb.ChatSettings, error) {

	return &userpb.ChatSettings{
		Questions: []*userpb.ChatSettings_Q{
			{Text: "Tell me something about the Big Bang so that I can explain it to my 5-year-old child2"},
			{Text: "Tell me something about the Big Bang so that I can explain it to my 5-year-old childw"},
			{Text: "Tell me something about the Big Bang so that I can explain it to my 5-year-old childe"},
		},
	}, nil
}

func (t *UserService) GetAppVersion(ctx context.Context, params *userpb.GetAppVersionParams) (*userpb.AppVersion, error) {

	return &userpb.AppVersion{
		Version:     "1.0.0",
		ForceUpdate: false,
		Description: "",
		DownloadUrl: &userpb.AppVersion_DownloadUrl{
			Ios:             "",
			Android:         "",
			FallbackIos:     "",
			FallbackAndroid: "",
		},
	}, nil
}

func (t *UserService) GetAppSettings(ctx context.Context, req *userpb.GetAppSettingsParams) (*userpb.AppSettings, error) {

	//userId := krathelper.FindUserId(ctx)

	settings, err := t.data.GrpcClients.AiAgentClient.GetPromptSettings(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return &userpb.AppSettings{
		Scenes: []*aiagentpb.Scene{
			{Value: "coverAnalysis", Name: "封面预测", Icon: "image", IsPopular: true, Description: ""},
			{Value: "analysis", Name: "爆款分析", Icon: "fire", IsNew: false},
			{Value: "duplicateScript", Name: "脚本复刻", Icon: "copy"},
			{Value: "preAnalysis", Name: "爆款预测", Icon: "fire"},
			{Value: "limitAnalysis", Name: "限流预测", Icon: "shield-halved"},
		},
		Prompts: helper.Mapping(settings.Settings, func(x *aiagentpb.PromptSetting) *aiagentpb.PromptSetting {
			x.Model = ""
			return x
		}),
	}, nil
}

func (t *UserService) GetUserById(ctx context.Context, params *userpb.GetUserByIdParams) (*typepb.User, error) {

	x, err := t.data.Repos.User.GetById(ctx, params.Id, true)
	if err != nil {
		return nil, err
	}

	return &typepb.User{
		Id:         conv.Str(x.ID),
		Username:   x.Nickname,
		UserAvatar: x.Avatar,
		Email:      x.Email,
		CreatedAt:  x.CreatedAt.Unix(),
		Phone:      x.Phone,
	}, nil
}

func (t *UserService) GetUserByEmail(ctx context.Context, params *userpb.GetUserByEmailParams) (*typepb.User, error) {

	list, _, err := t.data.Repos.User.List(ctx, repo.ListParams{
		Emails: []string{params.Email},
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	x := list[0]

	return &typepb.User{
		Id:         conv.Str(x.ID),
		Username:   x.Nickname,
		UserAvatar: x.Avatar,
		Email:      x.Email,
		CreatedAt:  x.CreatedAt.Unix(),
		Phone:      x.Phone,
	}, nil
}

func (t *UserService) GetUserByEmailOrPhone(ctx context.Context, params *userpb.GetUserByEmailOrPhoneParams) (*typepb.User, error) {

	list, _, err := t.data.Repos.User.List(ctx, repo.ListParams{
		EmailOrPhone: params.Value,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	x := list[0]

	return &typepb.User{
		Id:         conv.Str(x.ID),
		Username:   x.Nickname,
		UserAvatar: x.Avatar,
		Email:      x.Email,
		CreatedAt:  x.CreatedAt.Unix(),
		Phone:      x.Phone,
	}, nil
}

func (t *UserService) GetUser(ctx context.Context, req *userpb.GetUserParams) (*typepb.User, error) {

	userId := krathelper.FindUserId(ctx)

	if userId == "" {
		return nil, nil
	}

	user, err := t.userBiz.GetUserDetail(ctx, userId, req.GetAuthPlatform())
	if err != nil {
		log.Error("GetUserDetail err", err)

		if errors.IsNotFound(err) {
			return nil, errors.Unauthorized("no user found by id: "+userId, "")
		}

		return nil, err
	}

	state, err := t.data.GrpcClients.PaymentClient.GetCreditState(ctx, &paymentpb.GetCreditStateParams{
		UserId: userId,
	})

	if err != nil {
		log.Error("GetCreditState err", err)
		return nil, err
	}

	//creditSummary, err := t.userCredit.GetCreditSummary(ctx, userId)
	//if err != nil {
	//	log.Error("GetUser GetCreditSummary err", err)
	//	return nil, err
	//}
	//
	//ongoingPlanCredits, err := t.userCredit.ListOngoingPlanCredits(ctx, userId)
	//if err != nil {
	//	return nil, err
	//}

	user.CreditSummary = &typepb.CreditSummary{
		Total:     state.Total,
		Remaining: state.Total - state.Used,
	}
	user.IsVip = state.PlanId != "free"

	//member, err := t.memberBiz.GetMember(ctx, userId)
	//if err != nil {
	//	log.Error("GetUser Member err", err)
	//	return nil, err
	//}
	//
	//user.Membership = member
	//
	//t.repairUsed(ctx, user.Membership)

	return user, nil
}
func (t *UserService) GetUserV2(ctx context.Context, req *userpb.GetUserV2Params) (*typepb.User, error) {

	userId := krathelper.FindUserId(ctx)

	if userId == "" {
		return nil, nil
	}

	user, err := t.userBiz.GetUserDetail(ctx, userId, "")
	if err != nil {
		log.Error("GetUserDetail err", err)

		if errors.IsNotFound(err) {
			return nil, errors.Unauthorized("no user found by id: "+userId, "")
		}

		return nil, err
	}

	return user, nil
}

// todo
func (t *UserService) repairUsed(ctx context.Context, member *typepb.Member) {

	now := time.Now()
	mm := now.Format("2006-01")
	dd := now.Format("2006-01-02")

	newUsed := map[string]int64{}

	for k, _ := range member.Used {

		if ml, found := member.Quota.ModelLimits[k]; found {
			if ml.Day > 0 {
				newUsed[k] = mathz.Min(ml.Day, member.Used[k+":"+dd])
			} else if ml.Month > 0 {
				newUsed[k] = mathz.Min(ml.Month, member.Used[k+":"+mm])
			} else if ml.Total > 0 {
				newUsed[k] = mathz.Min(ml.Total, member.Used[k])
			}
		}

		if fl, found := member.Quota.FunctionLimits[k]; found {
			if fl.Day > 0 {
				newUsed[k] = mathz.Min(fl.Day, member.Used[k+":"+dd])
			} else if fl.Month > 0 {
				newUsed[k] = mathz.Min(fl.Month, member.Used[k+":"+mm])
			} else if fl.Total > 0 {
				newUsed[k] = mathz.Min(fl.Total, member.Used[k])
			}
		}

	}
	member.Used = newUsed
}

func (t *UserService) ListUsers(ctx context.Context, params *userpb.ListUsersParams) (*userpb.ListUsersResult, error) {

	users, err := t.userBiz.ListUsers(ctx, biz.ListUsersParams{
		UserID:  params.UserID,
		Emails:  params.Emails,
		Page:    params.Page,
		Size:    params.Size,
		Keyword: params.Keyword,
	})
	if err != nil {
		log.Errorw("ListUsers err", err)
		return nil, err
	}

	var userIDs []int64
	for _, user := range users.List {
		userIDs = append(userIDs, conv.Int64(user.Id))
	}

	//memberships, err := t.memberBiz.ListMemberships(ctx, userIDs)
	//if err != nil {
	//	return nil, err
	//}

	//for _, user := range users.List {
	//	user.Membership = memberships[conv.Int64(user.Id)]
	//}

	return users, nil
}

func (t *UserService) GetUsersSummary(ctx context.Context, params *userpb.GetUsersSummaryParams) (*userpb.UserSummary, error) {

	users, err := t.userBiz.GetUsersSummary(ctx)
	if err != nil {
		log.Errorw("GetUsersSummary err", err)
		return nil, err
	}

	return users, nil
}

func (t *UserService) UpdateUser(ctx context.Context, params *userpb.UpdateUserParams) (*emptypb.Empty, error) {
	_, err := t.userBiz.UpdateUser(ctx, params)
	if err != nil {
		log.Errorw("UpdateUser err", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *UserService) DeleteUser(ctx context.Context, params *userpb.DeleteUserParams) (*emptypb.Empty, error) {
	_, err := t.userBiz.Delete(ctx, params)
	if err != nil {
		log.Errorw("Delete err", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// todo ListUsers ?
//func (t *UserService) FindUsers(ctx context.Context, params *userpb.FindUsersParams) (*userpb.ListUsersResult, error) {
//
//	users, err := t.userBiz.ListUsers(ctx, biz.ListUsersParams{
//		Keyword: params.Keyword,
//	})
//	if err != nil {
//		log.Errorw("ListUsers err", err, "params", params)
//		return nil, err
//	}
//
//	var userIDs []int64
//	for _, user := range users.List {
//		userIDs = append(userIDs, conv.Int64(user.Id))
//	}
//
//	memberships, err := t.memberBiz.ListMemberships(ctx, userIDs)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, user := range users.List {
//		user.Membership = memberships[conv.Int64(user.Id)]
//	}
//
//	return users, nil
//}

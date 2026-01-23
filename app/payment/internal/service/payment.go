package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/emptypb"
	aiagentpb "store/api/aiagent"
	paymentpb "store/api/payment"
	typepb "store/api/payment/types"
	responsepb "store/api/public/response"
	"store/app/payment/configs"
	"store/app/payment/internal/biz"
	"store/app/payment/internal/data"
	"store/app/payment/internal/data/repo/ent/payment"
	"store/pkg/enums"
	"store/pkg/events"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"strings"
	"time"
)

type PaymentService struct {
	paymentpb.UnimplementedPaymentServer
	payment *biz.PaymentBiz
	credit  *biz.CreditBiz
	data    *data.Data
}

func NewPaymentService(paymentBiz *biz.PaymentBiz, credit *biz.CreditBiz, data *data.Data) *PaymentService {
	return &PaymentService{
		payment: paymentBiz,
		credit:  credit,
		data:    data,
	}
}

func (t PaymentService) productEvent(ctx context.Context, event events.PaymentSuccessEvent) {

	helper.Go(ctx, func(ctx context.Context) {
		err := t.data.KafkaClient.W().Write(ctx,
			events.TopicPaymentSuccess,
			kafka.Message{
				Key:   []byte(event.UniqueId),
				Value: conv.S2B(event),
			})

		if err != nil {
			log.Errorw("kafka write", "err", err, "event", event)
		}
	})

}

func (t PaymentService) GetMetadata(ctx context.Context, _ *empty.Empty) (*paymentpb.Metadata, error) {

	plans := configs.GetPlans()

	plans = helper.Filter(plans, func(param *paymentpb.Metadata_Plan) bool {
		return !param.Duplicated
	})

	return &paymentpb.Metadata{
		Plans: plans,
	}, nil
}

func (t PaymentService) GeneratePayment(ctx context.Context, params *paymentpb.GeneratePaymentParams) (*empty.Empty, error) {

	payments, err := t.data.Repos.EntClient.Payment.Query().Where(
		payment.Or(
			payment.PlanIDEQ("pro-annually"),
			payment.PlanIDEQ("pro-annual"),
			payment.PlanIDEQ("basic-annually"),
			payment.PlanIDEQ("basic-annual"),
		),
		payment.CreatedAtGT(time.Now().Add(-time.Hour*24*30)), // 整体未过期
		payment.ExpireAtLTE(time.Now()),                       // 当前周期过期了
	).Limit(1).All(ctx)

	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		return &empty.Empty{}, nil
	}

	x := payments[0]

	tx, err := t.data.Repos.EntClient.Tx(ctx)
	if err != nil {
		return nil, err
	}

	_, err = tx.Payment.Create().
		SetStatus(enums.PaymentStatus_Complete).
		SetPlatform(x.Platform).
		SetAmount(x.Amount).
		SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
		SetPlanID(x.PlanID).
		SetUserID(x.UserID).
		SetCreatedAt(x.CreatedAt). // 保留初始的创建时间，便于生成逻辑
		Save(ctx)

	if err != nil {
		return nil, err
	}

	tx.Payment.Update().SetStatus(enums.PaymentStatus_Expired)

	t.data.Repos.LocalCache.Delete("ongoingPayment:" + x.UserID)

	return &empty.Empty{}, nil
}

func (t PaymentService) RollbackCredit(ctx context.Context, params *paymentpb.RollbackCreditParams) (*emptypb.Empty, error) {

	ongoingPayment, err := t.payment.GetOngoingPayment(ctx, params.UserId)
	if err != nil {
		return nil, err
	}

	if ongoingPayment == nil {
		return nil, nil
	}

	err = t.data.Repos.RedisClient.DecrBy(ctx, fmt.Sprintf("credit:used.%d_%s", ongoingPayment.ID, params.UserId), params.Amount).Err()
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

//
//func (t PaymentService) MockDecrCredit(ctx context.Context, params *paymentpb.MockDecrCreditParams) (*paymentpb.CreditState, error) {
//	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
//		Value: params.EmailOrPhone,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	ongoingPayment, err := t.payment.GetOngoingPayment(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	if ongoingPayment == nil {
//		return nil, nil
//	}
//
//	err = t.data.Repos.RedisClient.DecrBy(ctx, fmt.Sprintf("credit:used.%d_%s", ongoingPayment.ID, user.Id), params.Amount).Err()
//	if err != nil {
//		return nil, err
//	}
//
//	state, err := t.payment.GetCreditState(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	return state, nil
//}
//
//func (t PaymentService) MockDiscardPayment(ctx context.Context, params *paymentpb.MockDiscardPaymentParams) (*paymentpb.CreditState, error) {
//
//	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
//		Value: params.EmailOrPhone,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	err = t.data.Repos.EntClient.Payment.Update().
//		Where(payment.UserID(user.Id)).
//		SetExpireAt(time.Now()).
//		Exec(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	t.data.Repos.LocalCache.Delete("ongoingPayment:" + user.Id)
//
//	state, err := t.payment.GetCreditState(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	return state, nil
//}
//
//func (t PaymentService) MockPayment(ctx context.Context, params *paymentpb.MockPaymentParams) (*paymentpb.MockPaymentResult, error) {
//
//	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
//		Value: params.EmailOrPhone,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	plan := configs.GetPlanById(params.PlanId)
//
//	if plan == nil {
//		return nil, errors.BadRequest("invalidPlan", "")
//	}
//
//	pay, err := t.data.Repos.EntClient.Payment.Create().
//		SetStatus(enums.PaymentStatus_Complete).
//		SetPlatform("manual").
//		SetAmount(plan.CnyAmount).
//		SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
//		SetPlanID(params.PlanId).
//		SetUserID(user.Id).
//		//OnConflictColumns(payment.FieldPlatform, payment.FieldSubID).
//		//UpdateExpireAt().
//		Save(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	t.productEvent(ctx, events.PaymentSuccessEvent{
//		UniqueId: conv.Str(pay.ID),
//		UserID:   user.Id,
//		//DeviceID: deviceId,
//		Amount:   plan.Amount,
//		Platform: "manual",
//	})
//
//	t.data.Repos.LocalCache.Delete("ongoingPayment:" + user.Id)
//
//	state, err := t.payment.GetCreditState(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	return &paymentpb.MockPaymentResult{
//		PlanId:    params.PlanId,
//		Remaining: state.Total - state.Used,
//		PaymentId: conv.Str(pay.ID),
//		UserId:    user.Id,
//	}, nil
//
//}

func (t PaymentService) CostCredit(ctx context.Context, params *paymentpb.CheckCreditParams) (*paymentpb.CheckCreditResult, error) {

	state, err := t.payment.GetCreditState(ctx, params.UserId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if state.Total >= state.Used+params.Cost {

		t.data.Repos.RedisClient.IncrBy(ctx, fmt.Sprintf("credit:used.%s_%s", state.PaymentId, params.UserId), params.Cost)

		// 双写
		if params.Cost > 0 {

			err := t.credit.Cost(ctx, params.UserId, params.Key, params.Cost)

			//err := t.data.Repos.EntClient.CreditConsume.Create().
			//	SetKey(params.Key).
			//	SetUserID(params.UserId).
			//	SetAmount(params.Cost).
			//	Exec(ctx)

			if err != nil {
				log.Error("cost credit", err)
			}
		}

		return &paymentpb.CheckCreditResult{
			Ok:        true,
			Remaining: state.Total - state.Used - params.Cost,
			Amount:    state.Total,
		}, nil
	} else {
		return &paymentpb.CheckCreditResult{
			Ok:        false,
			Remaining: state.Total - state.Used,
			Amount:    state.Total,
		}, nil
	}
}

func (t PaymentService) CallbackPayment(ctx context.Context, params *paymentpb.CallbackPaymentParams) (*emptypb.Empty, error) {

	log.Debugw("CallbackPayment", "", "params", params)

	var err error
	if strings.HasPrefix(params.OutTradeNo, "sur_") {

		_, err = t.data.GrpcClients.AiAgentClient.UpdateSurvey(ctx, &aiagentpb.UpdateSurveyParams{
			Id: params.OutTradeNo[4:],
		})

	} else {
		err = t.data.Repos.EntClient.Payment.Update().
			Where(payment.ID(conv.Int64(params.OutTradeNo))).
			SetStatus(enums.PaymentStatus_Complete).
			SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
			Exec(ctx)

		// todo
		all, err := t.data.Repos.EntClient.Payment.Query().Where(payment.ID(conv.Int64(params.OutTradeNo))).All(ctx)
		if err != nil {
			return nil, err
		}

		if len(all) > 0 {
			t.data.Repos.LocalCache.Delete("ongoingPayment:" + all[0].UserID)
		}

	}

	if err != nil {

		r, _ := http.RequestFromServerContext(ctx)

		log.Error("CallbackPayment err: ", err, "params", params, "req", r.RequestURI)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t PaymentService) CreateSubscription(ctx context.Context, _ *emptypb.Empty) (*paymentpb.CreatePaymentResult, error) {

	t.payment.CreateSubscription(ctx, "", &paymentpb.CreateParams{
		Level: "l1", Cycle: "month",
	})

	return nil, nil
}

func (t PaymentService) CreatePayment(ctx context.Context, params *paymentpb.CreatePaymentParams) (*paymentpb.CreatePaymentResult, error) {

	return &paymentpb.CreatePaymentResult{
		//Url: create.PayInfo,
	}, nil
}

func (t PaymentService) Create(ctx context.Context, params *paymentpb.CreateParams) (*responsepb.RedirectResponse, error) {

	userId := krathelper.RequireUserId(ctx)

	url, err := t.payment.Create(ctx, userId, params)
	if err != nil {
		log.Errorw("Create err", err, "userId", userId, "params", params)
		return nil, err
	}

	return &responsepb.RedirectResponse{
		Url: url,
	}, nil
}

func (t PaymentService) CreateV2(ctx context.Context, params *paymentpb.CreateParams) (*responsepb.RedirectResponse, error) {

	userId := krathelper.RequireUserId(ctx)

	url, err := t.payment.CreateV2(ctx, userId, params)
	if err != nil {
		log.Errorw("Create err", err, "userId", userId, "params", params)
		return nil, err
	}

	return &responsepb.RedirectResponse{
		Url: url,
	}, nil
}

func (t PaymentService) CreateV3(ctx context.Context, params *paymentpb.CreateParams) (*paymentpb.CreateResult, error) {

	userId := krathelper.RequireUserId(ctx)

	url, err := t.payment.CreateSubscriptionEmbedded(ctx, userId, params)
	if err != nil {
		log.Errorw("Create err", err, "userId", userId, "params", params)
		return nil, err
	}

	return &paymentpb.CreateResult{
		ClientSecret: url,
	}, nil
}

func (t PaymentService) GetCustomer(ctx context.Context, req *paymentpb.GetCustomerParams) (*typepb.Customer, error) {

	url, err := t.payment.GetStripeBillUrl(ctx, req.Version, req.UserEmail)
	if err != nil {
		log.Errorw("GetStripeBillUrl err", err)
		return nil, err
	}

	return &typepb.Customer{
		//UserId:     req.GetUserId(),
		BillingUrl: url,
	}, nil
}

func (t PaymentService) OnEvent(ctx context.Context, body []byte, signature, callback string) error {

	err := t.payment.OnEvent(ctx, body, signature, callback)
	if err != nil {
		log.Errorw("OnEvent err", err, "callback", callback, "body", string(body))
		return err
	}

	return nil
}

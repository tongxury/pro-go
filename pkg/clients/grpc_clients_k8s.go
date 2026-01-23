package clients

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpc2 "google.golang.org/grpc"
	aiagentpb "store/api/aiagent"
	bipb "store/api/bi"
	chatpb "store/api/chat"
	databankpb "store/api/databank"
	dexpb "store/api/dex"
	mgmtpb "store/api/mgmt"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	"store/pkg/confcenter"
)

type GrpcClients struct {
	AiAgentClient   aiagentpb.AIAgentServiceClient
	UserClient      userpb.UserClient
	MemberClient    userpb.MemberClient
	PaymentClient   paymentpb.PaymentClient
	DatabankClient  databankpb.DatabankClient
	ChatClient      chatpb.ChatClient
	BiClient        bipb.BiClient
	MgmtClient      mgmtpb.MgmtClient
	TgBotClient     dexpb.TgbotServiceClient
	TgUserClient    dexpb.UserServiceClient
	DexTradeClient  dexpb.TradeServiceClient
	DexMarketClient dexpb.MarketServiceClient
	DexWalletClient dexpb.WalletServiceClient
}

type GrpcClientsConfig struct {
	Namespace string
	Port      string
}

func NewGrpcClients(conf GrpcClientsConfig) (*GrpcClients, error) {

	conf.Port = "8090"

	clients := &GrpcClients{}

	userConn, err := getConn(confcenter.ServiceName_User, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.UserClient = userpb.NewUserClient(userConn)
	clients.MemberClient = userpb.NewMemberClient(userConn)

	paymentConn, err := getConn(confcenter.ServiceName_Payment, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.PaymentClient = paymentpb.NewPaymentClient(paymentConn)

	databankConn, err := getConn(confcenter.ServiceName_Databank, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.DatabankClient = databankpb.NewDatabankClient(databankConn)

	chatConn, err := getConn(confcenter.ServiceName_Chat, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.ChatClient = chatpb.NewChatClient(chatConn)

	biConn, err := getConn(confcenter.ServiceName_Bi, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.BiClient = bipb.NewBiClient(biConn)

	mgmtConn, err := getConn(confcenter.ServiceName_Mgmt, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.MgmtClient = mgmtpb.NewMgmtClient(mgmtConn)

	tgbotConn, err := getConn(confcenter.ServiceName_TgBot, conf.Namespace, conf.Port)
	if err != nil {
		return nil, err
	}
	clients.TgBotClient = dexpb.NewTgbotServiceClient(tgbotConn)
	clients.TgUserClient = dexpb.NewUserServiceClient(tgbotConn)
	clients.DexTradeClient = dexpb.NewTradeServiceClient(tgbotConn)
	clients.DexMarketClient = dexpb.NewMarketServiceClient(tgbotConn)
	clients.DexWalletClient = dexpb.NewWalletServiceClient(tgbotConn)

	//marketConn, err := getConn(confcenter.ServiceName_Market, conf.Namespace, conf.Port)
	//if err != nil {
	//	return nil, err
	//}
	//clients.MarketClient = marketpb.NewMarketClient(marketConn)

	log.Debugw("NewGrpcClients", "done", "namespace", conf.Namespace, "port", conf.Port)

	return clients, nil
}

func getConn(service, namespace, port string) (*grpc2.ClientConn, error) {
	return grpc.DialInsecure(
		context.Background(),
		//grpc.WithTimeout(20*time.Second),
		grpc.WithEndpoint(fmt.Sprintf("%s.%s.svc.cluster.local:%s", service, namespace, port)),
		grpc.WithMiddleware(
			tracing.Client(),
		),
	)
}

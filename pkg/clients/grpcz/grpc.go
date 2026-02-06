package grpcz

import (
	"context"
	aiagentpb "store/api/aiagent"
	bipb "store/api/bi"
	chatpb "store/api/chat"
	creditpb "store/api/credit"
	databankpb "store/api/databank"
	dexpb "store/api/dex"
	dexxpb "store/api/dexx"
	mgmtpb "store/api/mgmt"
	mondaypb "store/api/monday"
	notepb "store/api/note"
	paymentpb "store/api/payment"
	projpb "store/api/proj"
	userpb "store/api/user"
	ucpb "store/api/usercenter"
	voiceagentpb "store/api/voiceagent"
	"store/pkg/sdk/conv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpc2 "google.golang.org/grpc"
)

type Clients struct {
	UserClient        userpb.UserClient
	MemberClient      userpb.MemberClient
	PaymentClient     paymentpb.PaymentClient
	DatabankClient    databankpb.DatabankClient
	ChatClient        chatpb.ChatClient
	BiClient          bipb.BiClient
	MgmtClient        mgmtpb.MgmtClient
	DexTgBotClient    dexpb.TgbotServiceClient
	DexTgUserClient   dexpb.UserServiceClient
	DexTradeClient    dexpb.TradeServiceClient
	DexMarketClient   dexpb.MarketServiceClient
	DexMarketClient2  dexxpb.MarketServiceClient
	DexTradeClient2   dexxpb.TradeServiceClient
	DexWalletClient   dexpb.WalletServiceClient
	DexAdminClient    dexpb.AdminServiceClient
	NoteClient        notepb.NoteServiceClient
	AiAgentClient     aiagentpb.AIAgentServiceClient
	ProjAdminClient   projpb.ProjAdminServiceClient
	ProjProClient     projpb.ProjProServiceClient
	UserCenterClient  ucpb.UserServiceClient
	MondayAdminClient mondaypb.MondayAdminServiceClient
	CreditClient      creditpb.CreditServiceClient
	VoiceAgentClient  voiceagentpb.LiveKitServiceClient
}

type Configs struct {
	Credit      *Config
	UserCenter  *Config
	MondayAdmin *Config
	ProjAdmin   *Config
	ProjPro     *Config
	AiAgent     *Config
	VoiceAgent  *Config
	User        *Config
	Member      *Config
	Payment     *Config
	Databank    *Config
	Chat        *Config
	Bi          *Config
	Mgmt        *Config
	TgBot       *Config
	TgUser      *Config
	DexTrade    *Config
	DexMarket   *Config
	DexMarket2  *Config
	DexWallet   *Config
	DexAdmin    *Config
	Note        *Config
}

type Config struct {
	Endpoint string
}

func NewClients(configs Configs) (*Clients, error) {

	//conf.Port = "8090"
	//fmt.Sprintf("%s.%s.svc.cluster.local:%s", service, namespace, port)

	clients := &Clients{}

	if configs.Credit != nil {
		userConn, err := getConn(configs.Credit.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.CreditClient = creditpb.NewCreditServiceClient(userConn)
	}

	if configs.MondayAdmin != nil {
		userConn, err := getConn(configs.MondayAdmin.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.MondayAdminClient = mondaypb.NewMondayAdminServiceClient(userConn)
	}

	if configs.UserCenter != nil {
		userConn, err := getConn(configs.UserCenter.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.UserCenterClient = ucpb.NewUserServiceClient(userConn)
	}

	if configs.ProjAdmin != nil {
		userConn, err := getConn(configs.ProjAdmin.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.ProjAdminClient = projpb.NewProjAdminServiceClient(userConn)
	}

	if configs.ProjPro != nil {
		userConn, err := getConn(configs.ProjPro.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.ProjProClient = projpb.NewProjProServiceClient(userConn)
	}

	if configs.AiAgent != nil {
		userConn, err := getConn(configs.AiAgent.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.AiAgentClient = aiagentpb.NewAIAgentServiceClient(userConn)
	}

	if configs.Note != nil {
		userConn, err := getConn(configs.Note.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.NoteClient = notepb.NewNoteServiceClient(userConn)
	}

	if configs.User != nil {
		userConn, err := getConn(configs.User.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.UserClient = userpb.NewUserClient(userConn)
		clients.MemberClient = userpb.NewMemberClient(userConn)
	}

	if configs.Payment != nil {
		paymentConn, err := getConn(configs.Payment.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.PaymentClient = paymentpb.NewPaymentClient(paymentConn)
	}

	if configs.Databank != nil {
		databankConn, err := getConn(configs.Databank.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.DatabankClient = databankpb.NewDatabankClient(databankConn)
	}

	if configs.Chat != nil {
		chatConn, err := getConn(configs.Chat.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.ChatClient = chatpb.NewChatClient(chatConn)
	}

	if configs.Bi != nil {
		biConn, err := getConn(configs.Bi.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.BiClient = bipb.NewBiClient(biConn)
	}

	if configs.Mgmt != nil {
		mgmtConn, err := getConn(configs.Mgmt.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.MgmtClient = mgmtpb.NewMgmtClient(mgmtConn)
	}

	if configs.TgBot != nil {
		tgbotConn, err := getConn(configs.TgBot.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.DexTgBotClient = dexpb.NewTgbotServiceClient(tgbotConn)
		clients.DexTgUserClient = dexpb.NewUserServiceClient(tgbotConn)
		clients.DexTradeClient = dexpb.NewTradeServiceClient(tgbotConn)
		clients.DexMarketClient = dexpb.NewMarketServiceClient(tgbotConn)
		clients.DexWalletClient = dexpb.NewWalletServiceClient(tgbotConn)
	}

	if configs.DexMarket2 != nil {
		market2Conn, err := getConn(configs.DexMarket2.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.DexMarketClient2 = dexxpb.NewMarketServiceClient(market2Conn)
		clients.DexTradeClient2 = dexxpb.NewTradeServiceClient(market2Conn)
	}

	if configs.DexAdmin != nil {
		dexAdminConn, err := getConn(configs.DexAdmin.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.DexAdminClient = dexpb.NewAdminServiceClient(dexAdminConn)
	}

	if configs.VoiceAgent != nil {
		voiceAgentConn, err := getConn(configs.VoiceAgent.Endpoint)
		if err != nil {
			return nil, err
		}
		clients.VoiceAgentClient = voiceagentpb.NewLiveKitServiceClient(voiceAgentConn)
	}

	log.Debugw("NewGrpcClients", "done", "configs", conv.S2J(configs))

	return clients, nil
}

func getConn(endpoint string) (*grpc2.ClientConn, error) {
	return grpc.DialInsecure(
		context.Background(),
		grpc.WithTimeout(120*time.Second),
		grpc.WithEndpoint(endpoint),
		grpc.WithMiddleware(
			tracing.Client(),
		),
	)
}

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
	"store/app/databank/configs"
	"store/pkg/confcenter"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
// Name is the name of the compiled software.
//Name = configcenter.ServiceName_User
// Version is the version of the compiled software.
//Version string

// id, _ = os.Hostname()
)

func newApp(logger log.Logger, meta confcenter.Meta, gs *grpc.Server, hs *http.Server) *kratos.App {

	return kratos.New(
		//kratos.ID(id),
		kratos.Name(meta.Appname),
		//kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
		//kratos.Registrar(clients.NewK8sRegistrar()),
		//kratos.Registrar(reg),
	)
}

func main() {

	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("failed get config: %v", err)
	}

	app, cleanup, err := wireApp(
		*config,
		config.Meta,
		config.Server,
		config.Logger,
	)
	if err != nil {
		log.Fatalf("failed create app: %v", err)

	}
	defer cleanup()

	if err := app.Run(); err != nil {
		log.Fatalf("failed start app: %v", err)
	}
}

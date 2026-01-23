package main

import (
	"os"
	"store/app/usercenter/configs"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "github.com/go-sql-driver/mysql"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

//func init() {
//	flag.StringVar(&flagconf, "conf", "app/admin/configs", "config path, eg: -conf app/admin/configs")
//}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs, hs,
		),
	)
}

func main() {

	cf, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("failed get config: %v", err)
		panic(err)
	}

	app, cleanup, err := wireApp(
		*cf,
		cf.Meta,
		cf.Server,
		cf.Logger,
	)
	if err != nil {
		log.Fatalf("failed create app: %v", err)

	}
	defer cleanup()

	if err := app.Run(); err != nil {
		log.Fatalf("failed start app: %v", err)
	}
}

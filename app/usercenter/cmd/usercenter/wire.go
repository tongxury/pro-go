//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"store/app/usercenter/configs"
	"store/app/usercenter/internal/biz"
	"store/app/usercenter/internal/data"
	"store/app/usercenter/internal/server"
	"store/app/usercenter/internal/service"

	"store/pkg/confcenter"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(
	confcenter.Config[configs.BizConfig],
	confcenter.Meta,
	confcenter.Server,
	log.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

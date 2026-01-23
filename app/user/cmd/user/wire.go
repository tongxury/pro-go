//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"store/app/user/internal/biz"
	"store/app/user/internal/conf"
	"store/app/user/internal/data"
	"store/app/user/internal/server"
	"store/app/user/internal/service"
	"store/pkg/confcenter"
)

// wireApp init kratos application.
func wireApp(
	confcenter.Config[conf.BizConfig],
	confcenter.Meta,
	confcenter.Server,
	log.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

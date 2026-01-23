//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"store/app/payment/internal/biz"
	bizconf "store/app/payment/internal/conf"
	"store/app/payment/internal/data"
	"store/app/payment/internal/server"
	"store/app/payment/internal/service"
	"store/pkg/confcenter"
)

// wireApp init kratos application.
func wireApp(
	confcenter.Config[bizconf.BizConfig],
	confcenter.Meta,
	confcenter.Server,
	log.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

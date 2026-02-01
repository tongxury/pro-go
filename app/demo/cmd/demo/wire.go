//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"store/app/demo/configs"
	"store/app/demo/internal/biz"
	"store/app/demo/internal/data"
	"store/app/demo/internal/server"
	"store/app/demo/internal/service"

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

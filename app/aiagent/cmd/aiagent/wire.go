//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"store/app/aiagent/internal/biz"
	"store/app/aiagent/internal/conf"
	"store/app/aiagent/internal/data"
	"store/app/aiagent/internal/server"
	"store/app/aiagent/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"store/pkg/confcenter"
)

func wireApp(
	confcenter.Config[conf.BizConfig],
	confcenter.Meta,
	confcenter.Server,
	log.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"store/app/voiceagent/configs"
	"store/app/voiceagent/internal/biz"
	"store/app/voiceagent/internal/data"
	"store/app/voiceagent/internal/server"
	"store/app/voiceagent/internal/service"

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

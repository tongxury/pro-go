//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"store/app/proj-admin/configs"
	"store/app/proj-admin/internal/biz"
	"store/app/proj-admin/internal/data"
	"store/app/proj-admin/internal/server"
	"store/app/proj-admin/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"store/pkg/confcenter"
)

func wireApp(
	confcenter.Config[configs.BizConfig],
	confcenter.Meta,
	confcenter.Server,
	log.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

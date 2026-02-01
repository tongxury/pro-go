//go:build wireinject
// +build wireinject

package data

import (
	"store/app/demo/internal/data/repo"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, repo.ProviderSet)

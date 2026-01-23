package krathelper

import (
	"errors"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

func FindEnv(name string) string {

	c := config.New(config.WithSource(env.NewSource()))
	err := c.Load()
	if err != nil {
		return ""
	}

	val, err := c.Value(name).String()
	if errors.Is(err, config.ErrNotFound) {
		return ""
	}

	return val
}

func ScanFile[T any](filepath string, target *T) error {
	c := config.New(config.WithSource(file.NewSource(filepath)))
	defer c.Close()

	if err := c.Load(); err != nil {
		return err
	}

	if err := c.Scan(target); err != nil {
		return err
	}

	return nil
}

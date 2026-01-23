package conf

import "store/pkg/clients"

type BizConfig struct {
	Sync   clients.CanalConfig
	Domain string
}

package bitquery

type DEXPool struct {
	Block       Block       `json:"Block"`
	Pool        Pool        `json:"Pool"`
	Transaction Transaction `json:"Transaction"`
}

type DEXPools []*DEXPool

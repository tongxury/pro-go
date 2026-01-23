package confcenter

type IClient interface {
	LoadDatabase(options ...Options) (Database, error)
}

type Options struct {
	Listen   bool
	OnChange func(newValue string)
}

package clients

import (
	"github.com/go-mysql-org/go-mysql/canal"
)

type CanalClient struct {
	*canal.Canal
}

type EventHandler struct {
	canal.DummyEventHandler
	OnEvent func(e *canal.RowsEvent)
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
	h.OnEvent(e)
	return nil
}

func (h *EventHandler) String() string {
	return ""
}

func NewEventHandler(onEvent func(e *canal.RowsEvent)) *EventHandler {
	return &EventHandler{
		OnEvent: onEvent,
	}
}

type CanalConfig struct {
	Addr      string
	User      string
	Password  string
	Databases []string
	Tables    []string
}

func NewCanalClient(conf CanalConfig) *CanalClient {

	cfg := canal.NewDefaultConfig()
	cfg.Addr = conf.Addr
	cfg.User = conf.User
	cfg.Password = conf.Password
	// We only care table canal_test in test db
	cfg.Dump.Databases = conf.Databases
	cfg.Dump.Tables = conf.Tables
	cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)
	if err != nil {
		panic(err)
	}

	// Register a handler to handle RowsEvent
	//c.SetEventHandler()

	// Start canal
	//c.Run()

	return &CanalClient{
		Canal: c,
	}
}

func (t *CanalClient) Run() error {

	set, err := t.GetMasterPos()
	if err != nil {
		return err
	}

	return t.RunFrom(set)
}

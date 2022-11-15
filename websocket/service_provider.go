package websocket

import (
	"github.com/laravel-go-version/v2/pkg/Illuminate/Contracts/IConfig"
	"github.com/laravel-go-version/v2/pkg/Illuminate/Contracts/IFoundation"
	"github.com/laravel-go-version/v2/pkg/Illuminate/Contracts/IWebsocket"
	"sync"
)

type ServiceProvider struct {
}

func (s ServiceProvider) Register(application IFoundation.IApplication) {
	application.NamedSingleton("websocket", func(config IConfig.Config) IWebsocket.WebSocket {
		var wsConfig = config.Get("websocket").(Config)

		upgrader = wsConfig.Upgrader

		return &WebSocket{
			connMutex:   sync.Mutex{},
			fdMutex:     sync.Mutex{},
			connections: map[uint64]IWebsocket.WebSocketConnection{},
			count:       0,
		}
	})
}

func (s ServiceProvider) Start() error {
	return nil
}

func (s ServiceProvider) Stop() {
}

package entity

import (
	"net/url"
	"sync"
)

// Backend информация о бекэнде
type Backend struct {
	URL   *url.URL
	Alive bool
	Mux   sync.RWMutex
}

// SetAlive установить флаг жив ли бекэнд
func (b *Backend) SetAlive(alive bool) {
	b.Mux.Lock()
	b.Alive = alive
	b.Mux.Unlock()
}

// IsAlive проверка жив ли бекэнд
func (b *Backend) IsAlive() bool {
	b.Mux.RLock()
	alive := b.Alive
	b.Mux.RUnlock()

	return alive
}

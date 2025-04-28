package roundrobin

import (
	"log"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/KrllF/Cloud/entity"
	"github.com/KrllF/Cloud/internal/config"
)

// ServerPool информация о бекенде
type ServerPool struct {
	conf     config.Config
	backends []*entity.Backend
	current  uint64
	mu       sync.RWMutex
}

// NewServerPool конструктор ServerPool
func NewServerPool(conf config.Config) *ServerPool {
	serverPool := ServerPool{
		conf:     conf,
		backends: make([]*entity.Backend, 0, 3),
		mu:       sync.RWMutex{},
	}
	tokens := strings.Split(conf.BACKEND_SERVERS, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		serverPool.AddBackend(&entity.Backend{
			URL:   serverUrl,
			Alive: true,
		})
		log.Printf("сервер добавлен: %s\n", serverUrl)
	}

	return &serverPool
}

// AddBackend добавить сервер в пул серверов
func (s *ServerPool) AddBackend(backend *entity.Backend) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backends = append(s.backends, backend)
}

// NextIndex изменить индекс
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// MarkBackendStatus изменить статус бекэнда
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer вернуть следующий активный пир
func (s *ServerPool) GetNextPeer() *entity.Backend {
	next := s.NextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

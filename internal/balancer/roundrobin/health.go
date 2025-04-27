package roundrobin

import (
	"log"
	"net"
	"net/url"
	"time"
)

// healthCheck проверка статуса бекэнда
func (s *ServerPool) HealthCheck() {
	t := time.NewTicker(s.conf.CheckHealth)
	for {
		<-t.C
		log.Println("проверка роботоспособности серверов...")
		s.healthCheck()
		log.Println("проверка закончена")

	}
}

// healthCheck запрос к бекенду, чтобы обновить статус
func (s *ServerPool) healthCheck() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

// isBackendAlive проверка, работает ли бекенд
func isBackendAlive(u *url.URL) bool {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}

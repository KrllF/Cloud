package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KrllF/Cloud/internal/config"
	"github.com/KrllF/Cloud/internal/db"
	handler "github.com/KrllF/Cloud/internal/handler/http"
	"github.com/KrllF/Cloud/internal/handler/http/middleware"
	"github.com/KrllF/Cloud/internal/repository"
	"github.com/KrllF/Cloud/internal/server"
	"github.com/KrllF/Cloud/internal/service/balancer/roundrobin"
	"github.com/KrllF/Cloud/internal/service/ratelimiter/tokenbucket"
)

const (
	sizeCloser = 5
)

type (
	// ServerPool интерфейс пула серверов бекенда
	ServerPool interface {
		HealthCheck()
	}
	// Handler интерфейс хендлера
	Handler interface {
		Init() http.HandlerFunc
	}
	// Server интерфейс сервера
	Server interface {
		Run() error
	}
	// App DI
	App struct {
		Server     Server
		ServerPool ServerPool
		closer     []Closer
	}
	// Closer определяет методы для закрытия
	Closer interface {
		Close()
	}
)

// NewApp конструктор App
func NewApp(ctx context.Context) (*App, error) {
	conf, err := config.NewConfig("config.json", "config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config.NewConfig: %w", err)
	}

	cls := make([]Closer, 0, sizeCloser)
	dbs, err := db.NewDB(ctx, conf)
	if err != nil {
		return &App{}, fmt.Errorf("db.NewDB: %w", err)
	}
	repo, err := repository.NewRepository(dbs, conf)
	if err != nil {
		return &App{}, fmt.Errorf("repository.NewRepository: %w", err)
	}
	bal := roundrobin.NewServerPool(conf)
	hand := handler.NewHandler(bal)
	rate := tokenbucket.NewRateLimiter(conf, repo)
	rateLimiter := middleware.RateLimiter(rate)
	httpServ := server.NewServer(conf, hand.Init(rateLimiter))
	cls = append(cls, httpServ)

	return &App{httpServ, bal, cls}, nil
}

// Run запустить app
func (app *App) Run() error {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT)
	_, cancel := context.WithCancel(context.Background())

	go func() {
		log.Printf("Cервер HTTP прослушивается...")
		if err := app.Server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	go app.ServerPool.HealthCheck()

	sig := <-exit
	log.Printf("Получен сигнал завершения: %v", sig)
	cancel()
	for _, closer := range app.closer {
		closer.Close()
		log.Println("Успешно закрыто")
	}

	log.Println("Закрытие ресурсов прошло успешно")

	return nil
}

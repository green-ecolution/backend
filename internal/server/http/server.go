package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/akatranlp/sentinel/openid"
	"github.com/go-chi/chi/v5"
	"github.com/green-ecolution/backend/internal/worker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/green-ecolution/backend/internal/config"
	"github.com/green-ecolution/backend/internal/service"
)

type HTTPError struct {
	Error  string `json:"error"`
	Code   int    `json:"code"`
	Path   string `json:"path"`
	Method string `json:"method"`
} // @Name HTTPError

type Server struct {
	cfg      *config.Config
	services *service.Services
	ip       *openid.IdentitiyProvider
}

func NewServer(cfg *config.Config, services *service.Services, ip *openid.IdentitiyProvider) *Server {
	return &Server{
		cfg:      cfg,
		services: services,
		ip:       ip,
	}
}

func (s *Server) Run(ctx context.Context) error {
	app := fiber.New(fiber.Config{
		AppName:                  s.cfg.Dashboard.Title,
		ServerHeader:             s.cfg.Dashboard.Title,
		ErrorHandler:             errorHandler,
		EnableSplittingOnParsers: true,
	})

	app.Mount("/", s.middleware())

	router := chi.NewRouter()
	router.Mount("/auth", s.ip.Handler())
	router.Mount("/", adaptor.FiberApp(app))
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler: router,
	}

	go func() {
		slog.Info("starting plugin cleanup service: cleaning up unhealthy plugins")
		s.services.PluginService.StartCleanup(ctx)
	}()

	sensorStatusScheduler := worker.NewScheduler(3*time.Hour, worker.SchedulerFunc(s.services.SensorService.UpdateStatuses))
	go sensorStatusScheduler.Run(ctx)

	var onceADay = 24 * time.Hour

	wateringPlanStatusScheduler := worker.NewScheduler(onceADay, worker.SchedulerFunc(s.services.WateringPlanService.UpdateStatuses))
	go wateringPlanStatusScheduler.Run(ctx)

	clusterWateringStatusScheduler := worker.NewScheduler(onceADay, worker.SchedulerFunc(s.services.TreeClusterService.UpdateWateringStatuses))
	go clusterWateringStatusScheduler.Run(ctx)

	treeWateringStatusScheduler := worker.NewScheduler(onceADay, worker.SchedulerFunc(s.services.TreeService.UpdateWateringStatuses))
	go treeWateringStatusScheduler.Run(ctx)

	go func() {
		<-ctx.Done()
		slog.Info("shutting down http server")
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("error while shutting down http server", "error", err, "service", "fiber")
		}
	}()

	slog.Info("starting server", "url", s.cfg.Server.AppURL, "port", s.cfg.Server.Port, "service", "fiber")
	return server.ListenAndServe()
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(HTTPError{
		Error:  err.Error(),
		Code:   code,
		Path:   c.Path(),
		Method: c.Method(),
	})
}

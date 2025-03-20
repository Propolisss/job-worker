package app

import (
  "context"
  "os"
  "os/signal"
  "syscall"

  "flussonic_tz/internal/app/router"
  "flussonic_tz/internal/app/server"
  errs "flussonic_tz/internal/errors"

  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"

  "flussonic_tz/config"
  delivery "flussonic_tz/internal/delivery/http"
  "flussonic_tz/internal/repository/redis"
  "flussonic_tz/internal/service"
  "flussonic_tz/workerpool"

  "github.com/go-redis/redis/v8"
)

type App struct {
  cfg *config.Config
  srv *server.Server
  mx  *router.Router
}

func New() (*App, error) {
  cfg, err := config.New()
  if err != nil {
    return nil, err
  }

  return &App{
    cfg: cfg,
  }, nil
}

func (a *App) Run() {
  redisClient := redis.NewClient(&redis.Options{
    Addr:            a.cfg.Redis.Address,
    DialTimeout:     a.cfg.Redis.DialTimeout,
    ReadTimeout:     a.cfg.Redis.ReadTimeout,
    WriteTimeout:    a.cfg.Redis.WriteTimeout,
    PoolSize:        a.cfg.Redis.PoolSize,
    MinIdleConns:    a.cfg.Redis.MinIdleConns,
    PoolTimeout:     a.cfg.Redis.PoolTimeout,
    IdleTimeout:     a.cfg.Redis.IdleTimeout,
    MaxRetries:      a.cfg.Redis.MaxRetries,
    MinRetryBackoff: a.cfg.Redis.MinRetryBackoff,
    MaxRetryBackoff: a.cfg.Redis.MaxRetryBackoff,
  })
  repo := repository.NewRedisRepository(redisClient, a.cfg.Redis.QueueName)
  jobSvc := service.NewJobService(repo)
  delJob := delivery.NewJobHandler(jobSvc)

  workerPool := workerpool.NewWorkerPool(config.WrapWorkerPoolContext(context.Background(), &a.cfg.WorkerPool), repo)
  delWp := delivery.NewWorkerPoolHandler(workerPool)
  go workerPool.Start(context.Background())

  mx := router.New()
  mx.SetupMiddlewares()
  mx.SetupJob(delJob)
  mx.SetupWorkerPool(delWp)
  a.mx = mx

  srv := server.New(config.WrapServerContext(context.Background(), &a.cfg.Server), a.mx.Mux())
  a.srv = srv
  a.srv.Run()

  stop := make(chan os.Signal, 1)
  signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

  <-stop
  log.Info().Msg("shutting down server")

  ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
  defer cancel()

  if err := a.srv.Shutdown(ctx); err != nil {
    wrapped := errors.Wrap(err, errs.ErrShutdownServer)
    log.Fatal().Err(wrapped).Msg(wrapped.Error())
  }

  workerPool.Stop()
  err := redisClient.Close()
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrCloseRedis)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return
  }

  log.Info().Msg("Server is shut down gracefully")
}

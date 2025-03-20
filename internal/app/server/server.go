package server

import (
  "context"
  "fmt"
  "net/http"

  "flussonic_tz/config"
  errs "flussonic_tz/internal/errors"

  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"
)

type Server struct {
  srv *http.Server
  cfg *config.Server
}

func New(ctx context.Context, mx http.Handler) *Server {
  cfg := config.FromServerContext(ctx)
  return &Server{
    cfg: cfg,
    srv: &http.Server{
      Addr:         fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
      ReadTimeout:  cfg.ReadTimeout,
      WriteTimeout: cfg.WriteTimeout,
      IdleTimeout:  cfg.IdleTimeout,
      Handler:      mx,
    },
  }
}

func (s *Server) Run() {
  go func() {
    if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
      wrapped := errors.Wrap(err, errs.ErrStartServer)
      log.Fatal().Err(wrapped).Msg(wrapped.Error())
    }
  }()
}

func (s *Server) Shutdown(ctx context.Context) error {
  return s.srv.Shutdown(ctx)
}

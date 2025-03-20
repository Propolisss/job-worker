package router

import (
  delivery "flussonic_tz/internal/delivery/http"

  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
)

const (
  CompressLevel = 5
)

type Router struct {
  mx *chi.Mux
}

func New() *Router {
  return &Router{
    mx: chi.NewRouter(),
  }
}

func (r *Router) Mux() *chi.Mux {
  return r.mx
}

func (r *Router) SetupMiddlewares() {
  r.mx.Use(middleware.Recoverer)
  r.mx.Use(middleware.Compress(CompressLevel))
}

func (r *Router) SetupJob(handler *delivery.JobHandler) {
  r.mx.Post("/jobs", handler.CreateJob)
  r.mx.Get("/jobs/{job_id}", handler.GetJobStatus)
}

func (r *Router) SetupWorkerPool(handler *delivery.WorkerPoolHandler) {
  r.mx.Post("/pause", handler.Pause)
  r.mx.Post("/unpause", handler.Unpause)
}

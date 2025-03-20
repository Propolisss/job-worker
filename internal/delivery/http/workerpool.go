package http

import (
  "net/http"

  "github.com/rs/zerolog/log"
)

type WorkerPoolService interface {
  Pause()
  Unpause()
}

type WorkerPoolHandler struct {
  wpSvc WorkerPoolService
}

func NewWorkerPoolHandler(wpSvc WorkerPoolService) *WorkerPoolHandler {
  return &WorkerPoolHandler{
    wpSvc: wpSvc,
  }
}

func (h *WorkerPoolHandler) Pause(w http.ResponseWriter, r *http.Request) {
  h.wpSvc.Pause()

  _, err := w.Write([]byte("paused"))
  if err != nil {
    log.Error().Err(err).Msg("failed to write paused")
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func (h *WorkerPoolHandler) Unpause(w http.ResponseWriter, r *http.Request) {
  h.wpSvc.Unpause()

  _, err := w.Write([]byte("unpaused"))
  if err != nil {
    log.Error().Err(err).Msg("failed to write paused")
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

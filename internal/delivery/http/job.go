package http

import (
  "context"
  "encoding/json"
  "net/http"

  "flussonic_tz/internal/datastructures"
  errs "flussonic_tz/internal/errors"
  "flussonic_tz/models"

  "github.com/go-chi/chi"
  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"
)

type JobService interface {
  CreateJob(ctx context.Context, name string, score float64) (string, error)
  GetJob(ctx context.Context) (*models.Job, error)
  GetJobStatus(ctx context.Context, jobID string) (string, error)
}

type JobHandler struct {
  jobSvc JobService
}

func NewJobHandler(jobSvc JobService) *JobHandler {
  return &JobHandler{
    jobSvc: jobSvc,
  }
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
  defer func() {
    err := r.Body.Close()
    if err != nil {
      wrapped := errors.Wrap(err, errs.ErrCloseBody)
      log.Error().Err(wrapped).Msg(wrapped.Error())
    }
  }()

  var req models.JobRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    wrapped := errors.Wrap(err, errs.ErrDecodeBody)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    http.Error(w, wrapped.Error(), http.StatusBadRequest)
    return
  }

  id, err := h.jobSvc.CreateJob(r.Context(), req.Name, req.Score)
  if err != nil {
    log.Error().Err(err).Msg(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusAccepted)
  err = json.NewEncoder(w).Encode(datastructures.CreateJobResponse{Status: "created", ID: id})
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrEncodeResp)
    log.Error().Err(wrapped).Msg(err.Error())
    http.Error(w, wrapped.Error(), http.StatusInternalServerError)
  }
}

func (h *JobHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
  jobID := chi.URLParam(r, "job_id")

  status, err := h.jobSvc.GetJobStatus(r.Context(), jobID)
  if err != nil {
    log.Error().Err(err).Msg(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json")
  _, err = w.Write([]byte(status))
  if err != nil {
    wrapped := errors.Wrap(err, errs.ErrWriteStatus)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    http.Error(w, wrapped.Error(), http.StatusInternalServerError)
  }
}

package models

import "time"

type Job struct {
  ID         string    `json:"id"`
  Name       string    `json:"name"`
  Score      float64   `json:"score"`
  Status     string    `json:"status"`
  CreatedAt  time.Time `json:"created_at"`
  StartedAt  time.Time `json:"started_at"`
  FinishedAt time.Time `json:"finished_at"`
}

type JobRequest struct {
  Name  string  `json:"name" validate:"required"`
  Score float64 `json:"score" validate:"required"`
}

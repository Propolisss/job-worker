package datastructures

type CreateJobResponse struct {
  Status string `json:"status"`
  ID     string `json:"id"`
}

type GetStatusResponse struct {
  Status string `json:"status"`
}

package errs

// config
const (
  ErrInitializeConfig = "Error initializing config"
  ErrUnmarshalConfig  = "Error unmarshalling config"
  ErrReadConfig       = "Error reading config"
)

// repository/redis
const (
  ErrMarshalJob         = "Error marshalling job"
  ErrUnmarshalJob       = "Error unmarshalling job"
  ErrAddJob             = "Error adding job"
  ErrGetJob             = "Error getting job"
  ErrCastError          = "Error casting type"
  ErrUpdateJob          = "Error updating job"
  ErrGetJobStatus       = "Error getting job status"
  ErrUnmarshalJobStatus = "Error unmarshalling job status"
)

// pkg/generator
const (
  ErrGenerateID = "Error generating ID"
)

// delivery/http/job
const (
  ErrCloseBody   = "Error closing body"
  ErrDecodeBody  = "Error decoding body"
  ErrEncodeResp  = "Error encoding response"
  ErrWriteStatus = "Error writing status"
)

// workerpool
const (
  ErrCompleteJob = "Error change job status to complete"
  ErrFailJob     = "Error change job status to fail"
)

// internal/app/server
const (
  ErrStartServer = "Error starting server"
)

// internal/app/
const (
  ErrShutdownServer = "Error shutting down server"
  ErrCloseRedis     = "Error closing redis"
)

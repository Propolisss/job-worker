package config

import (
  "time"

  errs "flussonic_tz/internal/errors"

  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"
  "github.com/spf13/viper"
)

// worker pool
const (
  Workers          = 5
  JobLimit         = 100
  JobInterval      = 1 * time.Minute
  MaxRetries       = 3
  Timeout          = 3 * time.Second
  ErrorProbability = 0.1
)

// redis
const (
  RedisAddress         = "redis:6379"
  RedisQueueName       = "jobs"
  RedisPoolSize        = 20
  RedisMinIdleConns    = 5
  RedisMaxRetries      = 5
  RedisDialTimeout     = 10 * time.Second
  RedisReadTimeout     = 5 * time.Second
  RedisWriteTimeout    = 5 * time.Second
  RedisPoolTimeout     = 10 * time.Second
  RedisIdleTimeout     = 10 * time.Minute
  RedisMinRetryBackoff = 100 * time.Millisecond
  RedisMaxRetryBackoff = 1 * time.Second
)

// http server
const (
  Address         = "app"
  Port            = 8080
  ReadTimeout     = time.Second * 5
  WriteTimeout    = time.Second * 5
  ShutdownTimeout = time.Second * 30
  IdleTimeout     = time.Second * 60
)

type Config struct {
  WorkerPool WorkerPool `yaml:"workerpool" mapstructure:"workerpool"`
  Redis      Redis      `yaml:"redis" mapstructure:"redis"`
  Server     Server     `yaml:"server" mapstructure:"server"`
}

type WorkerPool struct {
  Workers          int           `yaml:"workers" mapstructure:"workers"`
  JobLimit         int           `yaml:"job_limit" mapstructure:"job_limit"`
  JobInterval      time.Duration `yaml:"job_interval" mapstructure:"job_interval"`
  MaxRetries       int           `yaml:"max_retries" mapstructure:"max_retries"`
  Timeout          time.Duration `yaml:"timeout" mapstructure:"timeout"`
  ErrorProbability float64       `yaml:"error_probability" mapstructure:"error_probability"`
}

type Redis struct {
  Address         string        `yaml:"address" mapstructure:"address"`
  DialTimeout     time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout"`
  ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
  WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
  PoolSize        int           `yaml:"pool_size" mapstructure:"pool_size"`
  MinIdleConns    int           `yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
  PoolTimeout     time.Duration `yaml:"pool_timeout" mapstructure:"pool_timeout"`
  IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
  MaxRetries      int           `yaml:"max_retries" mapstructure:"max_retries"`
  MinRetryBackoff time.Duration `yaml:"min_retry_backoff" mapstructure:"min_retry_backoff"`
  MaxRetryBackoff time.Duration `yaml:"max_retry_backoff" mapstructure:"max_retry_backoff"`
  QueueName       string        `yaml:"queue_name" mapstructure:"queue_name"`
}

type Server struct {
  Address         string        `yaml:"address" mapstructure:"address"`
  Port            int           `yaml:"port" mapstructure:"port"`
  ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
  WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
  ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
  IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

func New() (*Config, error) {
  log.Info().Msg("Initializing config")

  if err := setupViper(); err != nil {
    wrapped := errors.Wrap(err, errs.ErrInitializeConfig)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return nil, wrapped
  }

  var config Config
  if err := viper.Unmarshal(&config); err != nil {
    wrapped := errors.Wrap(err, errs.ErrUnmarshalConfig)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return nil, wrapped
  }

  log.Info().Msg("Config initialized")
  return &config, nil
}

func setupWorkerPool() {
  viper.SetDefault("workerpool.workers", Workers)
  viper.SetDefault("workerpool.job_limit", JobLimit)
  viper.SetDefault("workerpool.job_interval", JobInterval)
  viper.SetDefault("workerpool.max_retries", MaxRetries)
  viper.SetDefault("workerpool.timeout", Timeout)
  viper.SetDefault("workerpool.error_probability", ErrorProbability)
}

func setupRedis() {
  viper.SetDefault("redis.address", RedisAddress)
  viper.SetDefault("redis.dial_timeout", RedisDialTimeout)
  viper.SetDefault("redis.read_timeout", RedisReadTimeout)
  viper.SetDefault("redis.write_timeout", RedisWriteTimeout)
  viper.SetDefault("redis.pool_size", RedisPoolSize)
  viper.SetDefault("redis.min_idle_conns", RedisMinIdleConns)
  viper.SetDefault("redis.pool_timeout", RedisPoolTimeout)
  viper.SetDefault("redis.idle_timeout", RedisIdleTimeout)
  viper.SetDefault("redis.max_retries", RedisMaxRetries)
  viper.SetDefault("redis.min_retry_backoff", RedisMinRetryBackoff)
  viper.SetDefault("redis.max_retry_backoff", RedisMaxRetryBackoff)
  viper.SetDefault("redis.queue_name", RedisQueueName)
}

func setupServer() {
  viper.SetDefault("server.address", Address)
  viper.SetDefault("server.port", Port)
  viper.SetDefault("server.read_timeout", ReadTimeout)
  viper.SetDefault("server.write_timeout", WriteTimeout)
  viper.SetDefault("server.shutdown_timeout", ShutdownTimeout)
  viper.SetDefault("server.idle_timeout", IdleTimeout)
}

func setupViper() error {
  log.Info().Msg("Initializing viper")

  viper.SetConfigName("config")
  viper.SetConfigType("yml")
  viper.AddConfigPath(".")

  setupWorkerPool()
  setupRedis()
  setupServer()

  if err := viper.ReadInConfig(); err != nil {
    wrapped := errors.Wrap(err, errs.ErrReadConfig)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return wrapped
  }

  log.Info().Msg("Viper initialized")
  return nil
}

package config

import (
  "context"
)

type ContextWorkerPoolKey struct{}
type ContextRedisKey struct{}
type ContextServerKey struct{}

func WrapWorkerPoolContext(ctx context.Context, data interface{}) context.Context {
  return context.WithValue(ctx, ContextWorkerPoolKey{}, data)
}

func FromWorkerPoolContext(ctx context.Context) *WorkerPool {
  srv, ok := ctx.Value(ContextWorkerPoolKey{}).(*WorkerPool)
  if !ok {
    return nil
  }
  return srv
}

func WrapRedisContext(ctx context.Context, data interface{}) context.Context {
  return context.WithValue(ctx, ContextRedisKey{}, data)
}

func FromRedisContext(ctx context.Context) *Redis {
  redis, ok := ctx.Value(ContextRedisKey{}).(*Redis)
  if !ok {
    return nil
  }
  return redis
}

func WrapServerContext(ctx context.Context, data interface{}) context.Context {
  return context.WithValue(ctx, ContextServerKey{}, data)
}

func FromServerContext(ctx context.Context) *Server {
  srv, ok := ctx.Value(ContextServerKey{}).(*Server)
  if !ok {
    return nil
  }
  return srv
}

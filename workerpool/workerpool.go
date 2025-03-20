package workerpool

import (
  "context"
  "fmt"
  "math/rand"
  "sync"
  "sync/atomic"
  "time"

  errs "flussonic_tz/internal/errors"

  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"

  "flussonic_tz/config"
  "flussonic_tz/internal/service"

  "github.com/avast/retry-go"
)

var count uint64

type WorkerPool struct {
  cfg       *config.WorkerPool
  repo      service.JobRepository
  wg        *sync.WaitGroup
  semaphore chan struct{}
  done      chan struct{}
  doneJob   chan struct{}
  ticker    *time.Ticker
  pause     bool
  cond      *sync.Cond
}

func NewWorkerPool(ctx context.Context, repo service.JobRepository) *WorkerPool {
  cfg := config.FromWorkerPoolContext(ctx)
  return &WorkerPool{
    cfg:       cfg,
    repo:      repo,
    done:      make(chan struct{}),
    semaphore: make(chan struct{}, cfg.JobLimit),
    doneJob:   make(chan struct{}, cfg.JobLimit),
    wg:        &sync.WaitGroup{},
    ticker:    time.NewTicker(cfg.JobInterval),
    pause:     false,
    cond:      sync.NewCond(&sync.Mutex{}),
  }
}

func (wp *WorkerPool) Pause() {
  wp.cond.L.Lock()
  wp.pause = true
  log.Info().Msg("paused")
  wp.cond.L.Unlock()
}

func (wp *WorkerPool) Unpause() {
  wp.cond.L.Lock()
  wp.pause = false
  log.Info().Msg("unpaused")
  wp.cond.Broadcast()
  wp.cond.L.Unlock()
}

func (wp *WorkerPool) StartTicker() {
  defer wp.wg.Done()

  for {
    select {
    case <-wp.ticker.C:
      log.Info().Msg("Starting ticker...")
      // очищаем из семафора только те джобы, которые завершились, поэтому есть отдельный канал doneJob, в который
      // пишется после окончания функции
      length := len(wp.doneJob)
      atomic.StoreUint64(&count, 0)
      for length > 0 {
        <-wp.semaphore
        <-wp.doneJob
        length--
      }
    case <-wp.done:
      wp.ticker.Stop()
      return
    }
  }
}

func (wp *WorkerPool) Start(ctx context.Context) {
  wp.wg.Add(1)
  go wp.StartTicker()

  for range wp.cfg.Workers {
    wp.wg.Add(1)
    go wp.worker(ctx)
  }
}

func (wp *WorkerPool) wait() {
  wp.cond.L.Lock()
  for wp.pause {
    wp.cond.Wait()
  }
  wp.cond.L.Unlock()
}

func (wp *WorkerPool) worker(ctx context.Context) {
  defer wp.wg.Done()

  for {
    select {
    case <-wp.done:
      return
    default:
      // если paused, то будем ждать
      wp.wait()
      job, err := wp.repo.GetJob(ctx)
      if err != nil {
        continue
      }
      // для первого выполнения функции сразу пишем в семафор, чтобы не случилось ситуации, когда воркеры забрали
      // все задачи из очереди и запустили на каждую из них по горутине, при этом ни одна из горутин не успела
      // выполниться, а значит семафор пустой. эта запись защитит от этого и заблокируется, когда переполнится канал
      wp.semaphore <- struct{}{}
      go func() {
        retriesCount := 0
        err = retry.Do(
          func() error {
            // если ретрай, то тоже добавляем в семафор, иначе превысим лимит запросов(помимо основных вызовов будут
            // выполняться ретраи)
            if retriesCount != 0 {
              wp.semaphore <- struct{}{}
            }
            // пока ждали очередь на задачу могли поставить на паузу, поэтому если paused, то ждём
            wp.wait()
            ctxTime, cancel := context.WithTimeout(ctx, wp.cfg.Timeout)
            defer cancel()

            errChan := make(chan error, 1)
            go func() {
              errChan <- wp.PerformJob(fmt.Sprintf("%s:%s", job.Name, job.ID), []byte(job.Name))
            }()

            select {
            case <-ctxTime.Done():
              log.Info().Msg("timeout")
              return ctxTime.Err()
            case err = <-errChan:
              return err
            }
          },
          retry.Attempts(uint(wp.cfg.MaxRetries)),
          retry.Delay(1*time.Millisecond),
          retry.OnRetry(func(_ uint, _ error) {
            retriesCount++
            // так как функция выполнилась, то пишем в соответствующий канал
            wp.doneJob <- struct{}{}
          }),
        )
        // если функция в итоге выполнилась успешно, то последнее успешное выполнение надо записать в канал
        if err == nil {
          err = wp.repo.CompleteJob(ctx, job.ID)
          if err != nil {
            wrapped := errors.Wrap(err, errs.ErrCompleteJob)
            log.Info().Err(wrapped).Msg(wrapped.Error())
            return
          }

          wp.doneJob <- struct{}{}
        } else {
          // если все MaxRetries раз упала, то в OnRetry уже были записаны все выполнения
          err = wp.repo.FailJob(ctx, job.ID)
          if err != nil {
            wrapped := errors.Wrap(err, errs.ErrFailJob)
            log.Error().Err(wrapped).Msg(wrapped.Error())
            return
          }
        }
      }()
    }
  }
}

func (wp *WorkerPool) Stop() {
  close(wp.done)
  close(wp.semaphore)
  close(wp.doneJob)
  wp.wg.Wait()
}

func (wp *WorkerPool) PerformJob(name string, jobData []byte) error {
  atomic.AddUint64(&count, 1)
  fmt.Printf("Started job %s at %d task: %d\n", name, time.Now().UnixMilli(), atomic.LoadUint64(&count))

  // добавил случайную возможность вернуть ошибку чтобы работали ретраи
  if rand.Float64() < wp.cfg.ErrorProbability {
    return errors.New("Error happend")
  }

  // do some work. Random sleep time; max = 3s
  <-time.After(time.Duration(rand.Int63n(int64(3 * time.Second))))

  fmt.Printf("Finished job %s at %d\n", name, time.Now().UnixMilli())
  return nil
}

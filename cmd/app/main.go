package main

import (
  "log"

  "flussonic_tz/internal/app"
)

func main() {
  a, err := app.New()
  if err != nil {
    log.Fatal(err)
  }

  a.Run()
}

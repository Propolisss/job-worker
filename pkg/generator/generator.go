package generator

import (
  "crypto/rand"
  "encoding/hex"

  errs "flussonic_tz/internal/errors"

  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"
)

func GenerateID(length int) (string, error) {
  bytes := make([]byte, length)
  if _, err := rand.Read(bytes); err != nil {
    wrapped := errors.Wrap(err, errs.ErrGenerateID)
    log.Error().Err(wrapped).Msg(wrapped.Error())
    return "", wrapped
  }

  return hex.EncodeToString(bytes), nil
}

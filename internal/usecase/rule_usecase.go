package usecase

import (
  "context"
  "github.com/kai-zenn/backend_iot_absensi/internal/dto"
)

type IRuleUsecase interface {
  List(ctx context.Context) (dto.ListRuleResponse, error)
  Upsert(ctx context.Context, actorUserID uint64, req dto.UpsertRuleRequest) error
}

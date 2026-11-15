package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

type IRuleUsecase interface {
  List(ctx context.Context) (*dto.ListRuleResponse, error)
  Upsert(ctx context.Context, actorUserID uuid.UUID, req dto.UpsertRuleRequest) error
}

type ruleUsecase struct {
  ruleRepo repository.IRuleRepository
}

func NewRuleUsecase(rr repository.IRuleRepository) IRuleUsecase {
  return &ruleUsecase{
    ruleRepo: rr,
  }
}

func (ru *ruleUsecase) List(ctx context.Context) (*dto.ListRuleResponse, error){
  rules, err := ru.ruleRepo.List(ctx)
  if err != nil {
    return nil, err
  }

  data := make([]dto.RuleItem, 0, len(rules))  

  for _, r := range rules {
    createdBy := r.CreatedBy.String()
    data = append(data, dto.RuleItem{
      ID: r.ID.String(),
      EffectiveFrom: r.EffectiveFrom.Format("2006-01-02"),
      OnTimeUntil: r.OnTimeUntil.Format("15:04:05"),
      CreatedBy: &createdBy,
      CreatedAt: r.CreatedAt,
    })
  }

  return &dto.ListRuleResponse{
    Data: data,
  }, nil
}

func (ru *ruleUsecase) Upsert(ctx context.Context, actorUserID uuid.UUID, req dto.UpsertRuleRequest) error {
  effectiveFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
  if err != nil {
    return err
  }

  onTimeUntil, err := time.Parse("15:04:05", req.OnTimeUntil)
  if err != nil {
    return err
  }

  rule := entity.AttendanceRule{
    EffectiveFrom: effectiveFrom,
    OnTimeUntil: onTimeUntil,
    CreatedBy: &actorUserID,
    CreatedAt: time.Now(),
  }

  if err := ru.ruleRepo.Upsert(ctx, &rule); err != nil {
    return nil
  }

  return nil
}

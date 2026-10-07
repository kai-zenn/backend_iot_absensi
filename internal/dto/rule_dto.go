package dto

import "time"

type UpsertRuleRequest struct {
    EffectiveFrom string `json:"effectiveFrom" binding:"required,datetime=2006-01-02"` 
    OnTimeUntil   string `json:"onTimeUntil" binding:"required,datetime=15:04"`      
}

type RuleItem struct {
  ID            string    `json:"id"` 
  EffectiveFrom string    `json:"effectiveFrom"`
  OnTimeUntil   string    `json:"onTimeUntil"` 
  CreatedBy     *string   `json:"createdBy,omitempty"` 
  CreatedAt     time.Time `json:"createdAt"`
}

type ListRuleResponse struct {
  Data []RuleItem `json:"data"`
}

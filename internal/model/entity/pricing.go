package entity

import "time"

// Pricing is the singleton fee-sheet image row.
type Pricing struct {
	Id          uint64    `json:"id"          orm:"id,primary"`
	OriginalUrl string    `json:"originalUrl" orm:"original_url"`
	ThumbUrl    string    `json:"thumbUrl"    orm:"thumb_url"`
	UpdatedAt   time.Time `json:"updatedAt"   orm:"updated_at"`
}

// TableName returns the table name for GoFrame ORM.
func (Pricing) TableName() string { return "pricing" }

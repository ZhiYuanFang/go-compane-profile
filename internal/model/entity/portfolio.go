package entity

import "time"

// Portfolio is a project case study.
type Portfolio struct {
	Id               uint64    `json:"id"               orm:"id,primary"`
	Slug             string    `json:"slug"             orm:"slug"`
	SortOrder        int       `json:"sortOrder"        orm:"sort_order"`
	Address          string    `json:"address"          orm:"address"`
	Area             string    `json:"area"             orm:"area"`
	Style            string    `json:"style"            orm:"style"`
	HeartFlow        string    `json:"heartFlow"        orm:"heart_flow"`
	CoverOriginalUrl string    `json:"coverOriginalUrl" orm:"cover_original_url"`
	CoverThumbUrl    string    `json:"coverThumbUrl"    orm:"cover_thumb_url"`
	CreatedAt        time.Time `json:"createdAt"        orm:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt"        orm:"updated_at"`
}

// TableName returns the table name for GoFrame ORM.
func (Portfolio) TableName() string { return "portfolio" }

package entity

import "time"

// Activity is one promotional activity in the collection.
type Activity struct {
	Id               uint64    `json:"id"               orm:"id,primary"`
	Title            string    `json:"title"            orm:"title"`
	BodyHtml         string    `json:"bodyHtml"         orm:"body_html"`
	ImageOriginalUrl string    `json:"imageOriginalUrl" orm:"image_original_url"`
	ImageThumbUrl    string    `json:"imageThumbUrl"    orm:"image_thumb_url"`
	SortOrder        int       `json:"sortOrder"        orm:"sort_order"`
	CreatedAt        time.Time `json:"createdAt"        orm:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt"        orm:"updated_at"`
}

// TableName returns the table name for GoFrame ORM.
func (Activity) TableName() string { return "activity" }

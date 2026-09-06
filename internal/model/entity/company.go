package entity

import "time"

// Company is the singleton company profile row.
type Company struct {
	Id                 uint64    `json:"id"                 orm:"id,primary"`
	IntroTitle         string    `json:"introTitle"         orm:"intro_title"`
	IntroBody          string    `json:"introBody"          orm:"intro_body"`
	Values             string    `json:"values"             orm:"design_values"`
	YearsLabel         string    `json:"yearsLabel"         orm:"years_label"`
	Address            string    `json:"address"            orm:"address"`
	LogoOriginalUrl    string    `json:"logoOriginalUrl"    orm:"logo_original_url"`
	LogoThumbUrl       string    `json:"logoThumbUrl"       orm:"logo_thumb_url"`
	LogoHorOriginalUrl string    `json:"logoHorOriginalUrl" orm:"logo_hor_original_url"`
	LogoHorThumbUrl    string    `json:"logoHorThumbUrl"    orm:"logo_hor_thumb_url"`
	UpdatedAt          time.Time `json:"updatedAt"          orm:"updated_at"`
}

// TableName returns the table name for GoFrame ORM.
func (Company) TableName() string { return "company" }

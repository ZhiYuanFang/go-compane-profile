package service

import (
	"context"

	"go-compane-profile/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// DualURL is a thumb + original CDN pair.
type DualURL struct {
	Thumb    string `json:"thumb"`
	Original string `json:"original"`
}

// CompanyView is the public/admin company payload.
type CompanyView struct {
	IntroTitle string  `json:"introTitle"`
	IntroBody  string  `json:"introBody"`
	Values     string  `json:"values"`
	YearsLabel string  `json:"yearsLabel"`
	Address    string  `json:"address"`
	Logo       DualURL `json:"logo"`
	LogoHor    DualURL `json:"logoHor"`
}

// GetCompany returns the singleton company profile.
func GetCompany(ctx context.Context) (*CompanyView, error) {
	var row entity.Company
	err := g.DB().Model("company").Ctx(ctx).OrderAsc("id").Limit(1).Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return &CompanyView{
			Logo:    DualURL{},
			LogoHor: DualURL{},
		}, nil
	}
	return companyToView(&row), nil
}

// UpdateCompanyInput is admin PUT body for company.
type UpdateCompanyInput struct {
	IntroTitle string  `json:"introTitle"`
	IntroBody  string  `json:"introBody"`
	Values     string  `json:"values"`
	YearsLabel string  `json:"yearsLabel"`
	Address    string  `json:"address"`
	Logo       DualURL `json:"logo"`
	LogoHor    DualURL `json:"logoHor"`
}

// UpdateCompany upserts the singleton company row.
func UpdateCompany(ctx context.Context, in UpdateCompanyInput) (*CompanyView, error) {
	data := g.Map{
		"intro_title":           in.IntroTitle,
		"intro_body":            in.IntroBody,
		"design_values":         in.Values,
		"years_label":           in.YearsLabel,
		"address":               in.Address,
		"logo_original_url":     in.Logo.Original,
		"logo_thumb_url":        in.Logo.Thumb,
		"logo_hor_original_url": in.LogoHor.Original,
		"logo_hor_thumb_url":    in.LogoHor.Thumb,
	}
	var existing entity.Company
	_ = g.DB().Model("company").Ctx(ctx).OrderAsc("id").Limit(1).Scan(&existing)
	if existing.Id == 0 {
		if _, err := g.DB().Model("company").Ctx(ctx).Data(data).InsertAndGetId(); err != nil {
			return nil, err
		}
	} else {
		if _, err := g.DB().Model("company").Ctx(ctx).Where("id", existing.Id).Data(data).Update(); err != nil {
			return nil, err
		}
	}
	return GetCompany(ctx)
}

func companyToView(row *entity.Company) *CompanyView {
	return &CompanyView{
		IntroTitle: row.IntroTitle,
		IntroBody:  row.IntroBody,
		Values:     row.Values,
		YearsLabel: row.YearsLabel,
		Address:    row.Address,
		Logo: DualURL{
			Thumb:    row.LogoThumbUrl,
			Original: row.LogoOriginalUrl,
		},
		LogoHor: DualURL{
			Thumb:    row.LogoHorThumbUrl,
			Original: row.LogoHorOriginalUrl,
		},
	}
}

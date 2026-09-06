package service

import (
	"context"

	"go-compane-profile/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// GetPricing returns the singleton pricing dual URLs.
func GetPricing(ctx context.Context) (*DualURL, error) {
	var row entity.Pricing
	err := g.DB().Model("pricing").Ctx(ctx).OrderAsc("id").Limit(1).Scan(&row)
	if err != nil {
		return nil, err
	}
	return &DualURL{
		Thumb:    row.ThumbUrl,
		Original: row.OriginalUrl,
	}, nil
}

// UpdatePricing upserts the singleton pricing URLs.
func UpdatePricing(ctx context.Context, in DualURL) (*DualURL, error) {
	data := g.Map{
		"original_url": in.Original,
		"thumb_url":    in.Thumb,
	}
	var existing entity.Pricing
	_ = g.DB().Model("pricing").Ctx(ctx).OrderAsc("id").Limit(1).Scan(&existing)
	if existing.Id == 0 {
		if _, err := g.DB().Model("pricing").Ctx(ctx).Data(data).Insert(); err != nil {
			return nil, err
		}
	} else {
		if _, err := g.DB().Model("pricing").Ctx(ctx).Where("id", existing.Id).Data(data).Update(); err != nil {
			return nil, err
		}
	}
	return GetPricing(ctx)
}

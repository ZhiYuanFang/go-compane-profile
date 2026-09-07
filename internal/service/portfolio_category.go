package service

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Portfolio category enum values.
const (
	CategoryResidential  = "residential"
	CategoryCommercial   = "commercial"
	CategoryOffice       = "office"
	CategoryInstallation = "installation"
)

var validPortfolioCategories = map[string]struct{}{
	CategoryResidential:  {},
	CategoryCommercial:   {},
	CategoryOffice:       {},
	CategoryInstallation: {},
}

// NormalizePortfolioCategory trims and validates category. Empty is invalid unless allowEmpty.
func NormalizePortfolioCategory(category string, allowEmpty bool) (string, error) {
	c := strings.TrimSpace(category)
	if c == "" {
		if allowEmpty {
			return "", nil
		}
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "作品类别不能为空")
	}
	if _, ok := validPortfolioCategories[c]; !ok {
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "无效的作品类别")
	}
	return c, nil
}

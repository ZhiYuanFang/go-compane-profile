package v1

import (
	"go-compane-profile/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type GetCompanyReq struct {
	g.Meta `path:"/company" tags:"Public" method:"get" summary:"Get company profile"`
}

type GetCompanyRes struct {
	*service.CompanyView
}

type ListPortfoliosReq struct {
	g.Meta `path:"/portfolios" tags:"Public" method:"get" summary:"List portfolios"`
}

type ListPortfoliosRes struct {
	List []service.PortfolioListItem `json:"list"`
}

type GetPortfolioReq struct {
	g.Meta `path:"/portfolios/{id}" tags:"Public" method:"get" summary:"Get portfolio detail"`
	Id     string `json:"id" v:"required#作品 id 不能为空" in:"path"`
}

type GetPortfolioRes struct {
	*service.PortfolioDetail
}

type GetPricingReq struct {
	g.Meta `path:"/pricing" tags:"Public" method:"get" summary:"Get pricing image"`
}

type GetPricingRes struct {
	*service.DualURL
}

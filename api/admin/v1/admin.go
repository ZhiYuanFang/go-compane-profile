package v1

import (
	"go-compane-profile/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/login" tags:"AdminAuth" method:"post" summary:"Admin login"`
	Password string `json:"password" v:"required#密码不能为空"`
}

type LoginRes struct {
	Ok bool `json:"ok"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" tags:"AdminAuth" method:"post" summary:"Admin logout"`
}

type LogoutRes struct {
	Ok bool `json:"ok"`
}

type GetCompanyReq struct {
	g.Meta `path:"/company" tags:"AdminCompany" method:"get" summary:"Get company (admin)"`
}

type GetCompanyRes struct {
	*service.CompanyView
}

type UpdateCompanyReq struct {
	g.Meta `path:"/company" tags:"AdminCompany" method:"put" summary:"Update company"`
	service.UpdateCompanyInput
}

type UpdateCompanyRes struct {
	*service.CompanyView
}

type ListPortfoliosReq struct {
	g.Meta `path:"/portfolios" tags:"AdminPortfolio" method:"get" summary:"List portfolios (admin)"`
}

type ListPortfoliosRes struct {
	List []service.AdminPortfolio `json:"list"`
}

type GetPortfolioReq struct {
	g.Meta `path:"/portfolios/{id}" tags:"AdminPortfolio" method:"get" summary:"Get portfolio (admin)"`
	Id     string `json:"id" v:"required" in:"path"`
}

type GetPortfolioRes struct {
	*service.AdminPortfolio
}

type CreatePortfolioReq struct {
	g.Meta `path:"/portfolios" tags:"AdminPortfolio" method:"post" summary:"Create portfolio"`
	service.PortfolioInput
}

type CreatePortfolioRes struct {
	*service.AdminPortfolio
}

type UpdatePortfolioReq struct {
	g.Meta `path:"/portfolios/{id}" tags:"AdminPortfolio" method:"put" summary:"Update portfolio"`
	Id     string `json:"id" v:"required" in:"path"`
	service.PortfolioInput
}

type UpdatePortfolioRes struct {
	*service.AdminPortfolio
}

type DeletePortfolioReq struct {
	g.Meta `path:"/portfolios/{id}" tags:"AdminPortfolio" method:"delete" summary:"Delete portfolio"`
	Id     string `json:"id" v:"required" in:"path"`
}

type DeletePortfolioRes struct {
	Ok bool `json:"ok"`
}

type ReorderPortfoliosReq struct {
	g.Meta `path:"/portfolios/reorder" tags:"AdminPortfolio" method:"post" summary:"Reorder portfolios"`
	Ids    []string `json:"ids" v:"required#ids 不能为空"`
}

type ReorderPortfoliosRes struct {
	Ok bool `json:"ok"`
}

type SaveGalleryReq struct {
	g.Meta  `path:"/portfolios/{id}/gallery" tags:"AdminPortfolio" method:"put" summary:"Save portfolio gallery"`
	Id      string           `json:"id" v:"required" in:"path"`
	Renders []service.DualURL `json:"renders"`
	Reals   []service.DualURL `json:"reals"`
}

type SaveGalleryRes struct {
	*service.AdminPortfolio
}

type GetPricingReq struct {
	g.Meta `path:"/pricing" tags:"AdminPricing" method:"get" summary:"Get pricing (admin)"`
}

type GetPricingRes struct {
	*service.DualURL
}

type UpdatePricingReq struct {
	g.Meta `path:"/pricing" tags:"AdminPricing" method:"put" summary:"Update pricing"`
	service.DualURL
}

type UpdatePricingRes struct {
	*service.DualURL
}

type UploadReq struct {
	g.Meta `path:"/upload" tags:"AdminUpload" method:"post" mime:"multipart/form-data" summary:"Upload dual image pair"`
}

type UploadRes struct {
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
}

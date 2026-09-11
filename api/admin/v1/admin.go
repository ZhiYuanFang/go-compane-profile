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
	g.Meta   `path:"/portfolios" tags:"AdminPortfolio" method:"get" summary:"List portfolios (admin)"`
	Category string `json:"category" in:"query" v:"required#类别不能为空"`
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
	g.Meta   `path:"/portfolios/reorder" tags:"AdminPortfolio" method:"post" summary:"Reorder portfolios"`
	Category string   `json:"category" v:"required#类别不能为空"`
	Ids      []string `json:"ids" v:"required#ids 不能为空"`
}

type ReorderPortfoliosRes struct {
	Ok bool `json:"ok"`
}

type SaveGalleryReq struct {
	g.Meta `path:"/portfolios/{id}/gallery" tags:"AdminPortfolio" method:"put" summary:"Save portfolio gallery"`
	Id     string            `json:"id" v:"required" in:"path"`
	Images []service.DualURL `json:"images"`
}

type SaveGalleryRes struct {
	*service.AdminPortfolio
}

type ListActivitiesReq struct {
	g.Meta `path:"/activities" tags:"AdminActivity" method:"get" summary:"List activities (admin)"`
}

type ListActivitiesRes struct {
	List []service.ActivityListItem `json:"list"`
}

type GetActivityReq struct {
	g.Meta `path:"/activities/{id}" tags:"AdminActivity" method:"get" summary:"Get activity (admin)"`
	Id     string `json:"id" v:"required" in:"path"`
}

type GetActivityRes struct {
	*service.ActivityDetail
}

type CreateActivityReq struct {
	g.Meta `path:"/activities" tags:"AdminActivity" method:"post" summary:"Create activity"`
	service.ActivityInput
}

type CreateActivityRes struct {
	*service.ActivityDetail
}

type UpdateActivityReq struct {
	g.Meta `path:"/activities/{id}" tags:"AdminActivity" method:"put" summary:"Update activity"`
	Id     string `json:"id" v:"required" in:"path"`
	service.ActivityInput
}

type UpdateActivityRes struct {
	*service.ActivityDetail
}

type DeleteActivityReq struct {
	g.Meta `path:"/activities/{id}" tags:"AdminActivity" method:"delete" summary:"Delete activity"`
	Id     string `json:"id" v:"required" in:"path"`
}

type DeleteActivityRes struct {
	Ok bool `json:"ok"`
}

type ReorderActivitiesReq struct {
	g.Meta `path:"/activities/reorder" tags:"AdminActivity" method:"post" summary:"Reorder activities"`
	Ids    []string `json:"ids" v:"required#ids 不能为空"`
}

type ReorderActivitiesRes struct {
	Ok bool `json:"ok"`
}

type UploadReq struct {
	g.Meta `path:"/upload" tags:"AdminUpload" method:"post" mime:"multipart/form-data" summary:"Upload dual image pair"`
}

type UploadRes struct {
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
}

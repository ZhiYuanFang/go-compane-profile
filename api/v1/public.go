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
	g.Meta   `path:"/portfolios" tags:"Public" method:"get" summary:"List portfolios"`
	Page     int    `json:"page" in:"query" d:"1" dc:"Page number (1-based)"`
	PageSize int    `json:"pageSize" in:"query" d:"10" dc:"Items per page (max 50)"`
	Category string `json:"category" in:"query" dc:"Optional category filter; omit for all"`
}

type ListPortfoliosRes struct {
	List     []service.PortfolioListItem `json:"list"`
	Total    int                         `json:"total"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"pageSize"`
}

type GetPortfolioReq struct {
	g.Meta `path:"/portfolios/{id}" tags:"Public" method:"get" summary:"Get portfolio detail"`
	Id     string `json:"id" v:"required#作品 id 不能为空" in:"path"`
}

type GetPortfolioRes struct {
	*service.PortfolioDetail
}

type ListActivitiesReq struct {
	g.Meta `path:"/activities" tags:"Public" method:"get" summary:"List activities"`
}

type ListActivitiesRes struct {
	List []service.ActivityListItem `json:"list"`
}

type GetActivityReq struct {
	g.Meta `path:"/activities/{id}" tags:"Public" method:"get" summary:"Get activity detail"`
	Id     string `json:"id" v:"required#活动 id 不能为空" in:"path"`
}

type GetActivityRes struct {
	*service.ActivityDetail
}

type ViewPortfolioReq struct {
	g.Meta `path:"/portfolios/{id}/view" tags:"Public" method:"post" summary:"Increment portfolio view count"`
	Id     string `json:"id" v:"required#作品 id 不能为空" in:"path"`
}

type ViewPortfolioRes struct{}

type ViewActivityReq struct {
	g.Meta `path:"/activities/{id}/view" tags:"Public" method:"post" summary:"Increment activity view count"`
	Id     string `json:"id" v:"required#活动 id 不能为空" in:"path"`
}

type ViewActivityRes struct{}

type ViewAboutReq struct {
	g.Meta `path:"/company/about/view" tags:"Public" method:"post" summary:"Increment about-us view count"`
}

type ViewAboutRes struct{}

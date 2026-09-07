package public

import (
	"context"

	v1 "go-compane-profile/api/v1"
	"go-compane-profile/internal/service"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) GetCompany(ctx context.Context, req *v1.GetCompanyReq) (res *v1.GetCompanyRes, err error) {
	view, err := service.GetCompany(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetCompanyRes{CompanyView: view}, nil
}

func (c *Controller) ListPortfolios(ctx context.Context, req *v1.ListPortfoliosReq) (res *v1.ListPortfoliosRes, err error) {
	paged, err := service.ListPortfoliosPublic(ctx, req.Page, req.PageSize, req.Category)
	if err != nil {
		return nil, err
	}
	list := paged.List
	if list == nil {
		list = []service.PortfolioListItem{}
	}
	return &v1.ListPortfoliosRes{
		List:     list,
		Total:    paged.Total,
		Page:     paged.Page,
		PageSize: paged.PageSize,
	}, nil
}

func (c *Controller) GetPortfolio(ctx context.Context, req *v1.GetPortfolioReq) (res *v1.GetPortfolioRes, err error) {
	detail, err := service.GetPortfolioPublic(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetPortfolioRes{PortfolioDetail: detail}, nil
}

func (c *Controller) ListActivities(ctx context.Context, req *v1.ListActivitiesReq) (res *v1.ListActivitiesRes, err error) {
	list, err := service.ListActivities(ctx)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []service.ActivityListItem{}
	}
	return &v1.ListActivitiesRes{List: list}, nil
}

func (c *Controller) GetActivity(ctx context.Context, req *v1.GetActivityReq) (res *v1.GetActivityRes, err error) {
	detail, err := service.GetActivity(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetActivityRes{ActivityDetail: detail}, nil
}

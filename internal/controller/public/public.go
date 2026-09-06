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
	list, err := service.ListPortfoliosPublic(ctx)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []service.PortfolioListItem{}
	}
	return &v1.ListPortfoliosRes{List: list}, nil
}

func (c *Controller) GetPortfolio(ctx context.Context, req *v1.GetPortfolioReq) (res *v1.GetPortfolioRes, err error) {
	detail, err := service.GetPortfolioPublic(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetPortfolioRes{PortfolioDetail: detail}, nil
}

func (c *Controller) GetPricing(ctx context.Context, req *v1.GetPricingReq) (res *v1.GetPricingRes, err error) {
	pricing, err := service.GetPricing(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetPricingRes{DualURL: pricing}, nil
}

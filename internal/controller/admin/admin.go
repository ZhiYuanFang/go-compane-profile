package admin

import (
	"context"
	"io"
	"net/http"
	"time"

	v1 "go-compane-profile/api/admin/v1"
	"go-compane-profile/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	if err = service.CheckAdminPassword(req.Password); err != nil {
		return nil, err
	}
	token, maxAge, err := service.IssueAdminSession()
	if err != nil {
		return nil, err
	}
	r := g.RequestFromCtx(ctx)
	r.Cookie.SetCookie(service.AdminSessionCookie, token, "", "/", time.Duration(maxAge)*time.Second, ghttp.CookieOptions{
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		HttpOnly: true,
	})
	return &v1.LoginRes{Ok: true}, nil
}

func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Cookie.SetCookie(service.AdminSessionCookie, "", "", "/", -1*time.Second, ghttp.CookieOptions{
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		HttpOnly: true,
	})
	return &v1.LogoutRes{Ok: true}, nil
}

func (c *Controller) GetCompany(ctx context.Context, req *v1.GetCompanyReq) (res *v1.GetCompanyRes, err error) {
	view, err := service.GetCompany(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetCompanyRes{CompanyView: view}, nil
}

func (c *Controller) UpdateCompany(ctx context.Context, req *v1.UpdateCompanyReq) (res *v1.UpdateCompanyRes, err error) {
	view, err := service.UpdateCompany(ctx, req.UpdateCompanyInput)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateCompanyRes{CompanyView: view}, nil
}

func (c *Controller) ListPortfolios(ctx context.Context, req *v1.ListPortfoliosReq) (res *v1.ListPortfoliosRes, err error) {
	list, err := service.ListPortfoliosAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []service.AdminPortfolio{}
	}
	return &v1.ListPortfoliosRes{List: list}, nil
}

func (c *Controller) GetPortfolio(ctx context.Context, req *v1.GetPortfolioReq) (res *v1.GetPortfolioRes, err error) {
	item, err := service.GetPortfolioAdmin(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetPortfolioRes{AdminPortfolio: item}, nil
}

func (c *Controller) CreatePortfolio(ctx context.Context, req *v1.CreatePortfolioReq) (res *v1.CreatePortfolioRes, err error) {
	item, err := service.CreatePortfolio(ctx, req.PortfolioInput)
	if err != nil {
		return nil, err
	}
	return &v1.CreatePortfolioRes{AdminPortfolio: item}, nil
}

func (c *Controller) UpdatePortfolio(ctx context.Context, req *v1.UpdatePortfolioReq) (res *v1.UpdatePortfolioRes, err error) {
	item, err := service.UpdatePortfolio(ctx, req.Id, req.PortfolioInput)
	if err != nil {
		return nil, err
	}
	return &v1.UpdatePortfolioRes{AdminPortfolio: item}, nil
}

func (c *Controller) DeletePortfolio(ctx context.Context, req *v1.DeletePortfolioReq) (res *v1.DeletePortfolioRes, err error) {
	if err = service.DeletePortfolio(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeletePortfolioRes{Ok: true}, nil
}

func (c *Controller) ReorderPortfolios(ctx context.Context, req *v1.ReorderPortfoliosReq) (res *v1.ReorderPortfoliosRes, err error) {
	if err = service.ReorderPortfolios(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.ReorderPortfoliosRes{Ok: true}, nil
}

func (c *Controller) SaveGallery(ctx context.Context, req *v1.SaveGalleryReq) (res *v1.SaveGalleryRes, err error) {
	renders := req.Renders
	reals := req.Reals
	if renders == nil {
		renders = []service.DualURL{}
	}
	if reals == nil {
		reals = []service.DualURL{}
	}
	item, err := service.SavePortfolioGallery(ctx, req.Id, renders, reals)
	if err != nil {
		return nil, err
	}
	return &v1.SaveGalleryRes{AdminPortfolio: item}, nil
}

func (c *Controller) GetPricing(ctx context.Context, req *v1.GetPricingReq) (res *v1.GetPricingRes, err error) {
	pricing, err := service.GetPricing(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetPricingRes{DualURL: pricing}, nil
}

func (c *Controller) UpdatePricing(ctx context.Context, req *v1.UpdatePricingReq) (res *v1.UpdatePricingRes, err error) {
	pricing, err := service.UpdatePricing(ctx, req.DualURL)
	if err != nil {
		return nil, err
	}
	return &v1.UpdatePricingRes{DualURL: pricing}, nil
}

func (c *Controller) Upload(ctx context.Context, req *v1.UploadReq) (res *v1.UploadRes, err error) {
	r := g.RequestFromCtx(ctx)
	origFile := r.GetUploadFile("original")
	thumbFile := r.GetUploadFile("thumb")
	if origFile == nil || thumbFile == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "需要 multipart 字段 original 与 thumb")
	}
	category := r.Get("category", "misc").String()

	origFH, err := origFile.Open()
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "读取原图失败")
	}
	defer origFH.Close()
	thumbFH, err := thumbFile.Open()
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "读取缩略图失败")
	}
	defer thumbFH.Close()

	origData, err := io.ReadAll(io.LimitReader(origFH, service.MaxOriginalBytes+1))
	if err != nil {
		return nil, err
	}
	thumbData, err := io.ReadAll(io.LimitReader(thumbFH, service.MaxThumbBytes+1))
	if err != nil {
		return nil, err
	}
	result, err := service.UploadDualImageBytes(origData, thumbData, category, origFile.Filename, thumbFile.Filename)
	if err != nil {
		return nil, err
	}
	return &v1.UploadRes{Original: result.Original, Thumb: result.Thumb}, nil
}

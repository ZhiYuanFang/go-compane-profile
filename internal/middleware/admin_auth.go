package middleware

import (
	"go-compane-profile/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

// AdminAuth protects /admin/api/* except login.
func AdminAuth(r *ghttp.Request) {
	if r.Method == "OPTIONS" {
		r.Middleware.Next()
		return
	}
	path := r.URL.Path
	if path == "/admin/api/login" || path == "/admin/api/login/" {
		r.Middleware.Next()
		return
	}
	token := r.Cookie.Get(service.AdminSessionCookie).String()
	if token == "" {
		token = r.Header.Get("X-Admin-Session")
	}
	if err := service.ValidateAdminSession(token); err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	r.Middleware.Next()
}

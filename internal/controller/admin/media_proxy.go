package admin

import (
	"net/http"
	"strings"

	"go-compane-profile/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

// HandleMediaProxy streams an OSS image back through the admin origin so the
// browser can read existing CDN images without cross-origin restrictions.
// Registered without MiddlewareHandlerResponse so the body stays binary.
func HandleMediaProxy(r *ghttp.Request) {
	rawURL := strings.TrimSpace(r.Get("url").String())
	if rawURL == "" {
		writeProxyError(r, http.StatusBadRequest, "缺少 url 参数")
		return
	}

	data, contentType, err := service.ReadObjectBytes(r.Context(), rawURL)
	if err != nil {
		writeProxyError(r, http.StatusBadRequest, err.Error())
		return
	}

	r.Response.Header().Set("Content-Type", contentType)
	r.Response.Header().Set("Cache-Control", "no-store")
	r.Response.Write(data)
}

func writeProxyError(r *ghttp.Request, status int, message string) {
	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.Response.WriteStatus(status, message)
}

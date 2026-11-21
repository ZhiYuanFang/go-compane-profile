package cmd

import (
	"context"
	"strings"

	"go-compane-profile/internal/config"
	"go-compane-profile/internal/controller/admin"
	"go-compane-profile/internal/controller/public"
	"go-compane-profile/internal/middleware"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start company profile http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			cfg := config.LoadAppConfig()
			if cfg.MySQLLink != "" {
				gdb.SetConfigGroup("default", gdb.ConfigGroup{
					gdb.ConfigNode{Link: cfg.MySQLLink},
				})
			}
			s := g.Server()
			if cfg.ServerAddress != "" {
				s.SetAddr(cfg.ServerAddress)
			}

			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Middleware(corsMiddleware)
				group.Bind(public.New())
			})

			s.Group("/admin/api", func(group *ghttp.RouterGroup) {
				group.Middleware(corsMiddleware)
				group.Middleware(middleware.AdminAuth)
				// NDJSON stream must not use MiddlewareHandlerResponse wrapping.
				group.POST("/upload/stream", admin.HandleUploadStream)
				group.Group("/", func(api *ghttp.RouterGroup) {
					api.Middleware(ghttp.MiddlewareHandlerResponse)
					api.Bind(admin.New())
				})
			})

			adminDir := "resource/public/admin"
			if !gfile.Exists(gfile.Join(adminDir, "index.html")) && gfile.Exists("admin/dist/index.html") {
				adminDir = "admin/dist"
			}
			if gfile.Exists(gfile.Join(adminDir, "index.html")) {
				serveAdmin := func(r *ghttp.Request) {
					rel := strings.TrimPrefix(r.URL.Path, "/admin")
					rel = strings.TrimPrefix(rel, "/")
					if rel == "" {
						r.Response.ServeFile(gfile.Join(adminDir, "index.html"))
						return
					}
					candidate := gfile.Join(adminDir, rel)
					if gfile.Exists(candidate) && !gfile.IsDir(candidate) {
						r.Response.ServeFile(candidate)
						return
					}
					r.Response.ServeFile(gfile.Join(adminDir, "index.html"))
				}
				s.BindHandler("GET:/admin", serveAdmin)
				s.BindHandler("GET:/admin/{*path}", serveAdmin)
			}

			g.Log().Infof(ctx, "listening on %s (db=%s)", cfg.ServerAddress, maskLink(cfg.MySQLLink))
			s.Run()
			return nil
		},
	}
)

func corsMiddleware(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func maskLink(link string) string {
	if link == "" {
		return "(empty)"
	}
	return "mysql:***@" + trimAfterAt(link)
}

func trimAfterAt(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '@' {
			return s[i+1:]
		}
	}
	return "configured"
}

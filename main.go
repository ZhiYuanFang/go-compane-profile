package main

import (
	_ "go-compane-profile/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"go-compane-profile/internal/cmd"
	"go-compane-profile/internal/config"
)

func main() {
	// Load env (MYSQL_LINK required) before GoFrame / DB init.
	config.LoadAppConfig()
	cmd.Main.Run(gctx.GetInitCtx())
}

package config

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

// AppConfig holds runtime settings loaded from environment.
type AppConfig struct {
	ServerAddress      string
	MySQLLink          string
	AdminPassword      string
	AdminSessionSecret string
	OSSAccessKeyID     string
	OSSAccessKeySecret string
	OSSBucket          string
	OSSEndpoint        string
	OSSPrefix          string
	CDNBase            string
}

var (
	appOnce sync.Once
	appCfg  AppConfig
)

// LoadDotEnvFiles loads KEY=VALUE pairs from the first existing file among paths
// without overriding already-set environment variables.
func LoadDotEnvFiles(paths ...string) {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			idx := strings.IndexByte(line, '=')
			if idx <= 0 {
				continue
			}
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			if len(val) >= 2 {
				if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
					val = val[1 : len(val)-1]
				}
			}
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
		_ = f.Close()
		return
	}
}

// LoadAppConfig loads (once) application config from the environment.
// MYSQL_LINK is required and must be set in env / .env.prod.
func LoadAppConfig() AppConfig {
	appOnce.Do(func() {
		LoadDotEnvFiles(
			"manifest/docker/env/.env.prod",
			"manifest/docker/env/.env.local",
			".env",
		)
		link := strings.TrimSpace(os.Getenv("MYSQL_LINK"))
		if link == "" {
			panic("config: MYSQL_LINK is required")
		}
		_ = os.Setenv("GF_DATABASE_DEFAULT_LINK", link)
		if addr := strings.TrimSpace(os.Getenv("SERVER_ADDRESS")); addr != "" {
			_ = os.Setenv("GF_SERVER_ADDRESS", addr)
		}
		prefix := envOr("OSS_PREFIX", "muhou/")
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		cdn := strings.TrimRight(envOr("CDN_BASE", "https://resorce.cuplay.top"), "/")
		sessionSecret := strings.TrimSpace(os.Getenv("ADMIN_SESSION_SECRET"))
		adminPass := envOr("ADMIN_PASSWORD", "abby")
		if sessionSecret == "" {
			sessionSecret = adminPass + "|muhou-admin-session-salt"
		}
		appCfg = AppConfig{
			ServerAddress:      envOr("SERVER_ADDRESS", ":9100"),
			MySQLLink:          link,
			AdminPassword:      adminPass,
			AdminSessionSecret: sessionSecret,
			OSSAccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
			OSSAccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
			OSSBucket:          envOr("OSS_BUCKET", "pang-bao"),
			OSSEndpoint:        envOr("OSS_ENDPOINT", "oss-cn-beijing.aliyuncs.com"),
			OSSPrefix:          prefix,
			CDNBase:            cdn,
		}
	})
	return appCfg
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-compane-profile/internal/config"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	AdminSessionCookie = "admin_session"
	adminSessionTTL    = 7 * 24 * time.Hour
)

// IssueAdminSession returns a signed session token valid for adminSessionTTL.
func IssueAdminSession() (token string, maxAge int, err error) {
	cfg := config.LoadAppConfig()
	exp := time.Now().Add(adminSessionTTL).Unix()
	payload := strconv.FormatInt(exp, 10)
	sig := sign(payload, cfg.AdminSessionSecret)
	token = payload + "." + sig
	maxAge = int(adminSessionTTL.Seconds())
	return token, maxAge, nil
}

// ValidateAdminSession verifies the HMAC session cookie value.
func ValidateAdminSession(token string) error {
	cfg := config.LoadAppConfig()
	token = strings.TrimSpace(token)
	if token == "" {
		return gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return gerror.NewCode(gcode.CodeNotAuthorized, "会话无效")
	}
	payload, sig := parts[0], parts[1]
	if !hmac.Equal([]byte(sig), []byte(sign(payload, cfg.AdminSessionSecret))) {
		return gerror.NewCode(gcode.CodeNotAuthorized, "会话无效")
	}
	exp, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return gerror.NewCode(gcode.CodeNotAuthorized, "会话无效")
	}
	if time.Now().Unix() > exp {
		return gerror.NewCode(gcode.CodeNotAuthorized, "会话已过期")
	}
	return nil
}

// CheckAdminPassword compares the submitted password with ADMIN_PASSWORD.
func CheckAdminPassword(password string) error {
	cfg := config.LoadAppConfig()
	if password == "" || password != cfg.AdminPassword {
		return gerror.NewCode(gcode.CodeValidationFailed, "密码错误")
	}
	return nil
}

func sign(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// OpaqueToken helpers kept for potential bearer use.
func EncodeOpaque(raw string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func DecodeOpaque(enc string) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(enc)
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	return string(b), nil
}

package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
	"time"

	"go-compane-profile/internal/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/guid"
)

const (
	MaxOriginalBytes = 3 * 1024 * 1024 // 3MiB
	MaxThumbBytes    = 512 * 1024      // 0.5MiB
)

// DualUploadResult is a CDN URL pair for original + thumb.
type DualUploadResult struct {
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
	OrigKey  string `json:"-"`
	ThumbKey string `json:"-"`
}

// ProgressFunc reports overall weighted OSS upload percent and phase.
type ProgressFunc func(pct int, phase string)

var (
	ossOnce   sync.Once
	ossClient *oss.Client
	ossBucket *oss.Bucket
	ossErr    error
)

func getOSSBucket() (*oss.Bucket, error) {
	ossOnce.Do(func() {
		cfg := config.LoadAppConfig()
		if cfg.OSSAccessKeyID == "" || cfg.OSSAccessKeySecret == "" {
			ossErr = gerror.NewCode(gcode.CodeInvalidConfiguration, "OSS 凭证未配置")
			return
		}
		endpoint := cfg.OSSEndpoint
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = "https://" + endpoint
		}
		ossClient, ossErr = oss.New(endpoint, cfg.OSSAccessKeyID, cfg.OSSAccessKeySecret)
		if ossErr != nil {
			return
		}
		ossBucket, ossErr = ossClient.Bucket(cfg.OSSBucket)
	})
	if ossErr != nil {
		return nil, ossErr
	}
	return ossBucket, nil
}

type countingReader struct {
	ctx    context.Context
	r      io.Reader
	n      *int64
	report func()
}

func (c *countingReader) Read(p []byte) (int, error) {
	if c.ctx != nil {
		if err := c.ctx.Err(); err != nil {
			return 0, err
		}
	}
	n, err := c.r.Read(p)
	if n > 0 {
		*c.n += int64(n)
		if c.report != nil {
			c.report()
		}
	}
	return n, err
}

// UploadDualImage uploads original + thumb with distinct keys under muhou/.
func UploadDualImage(original, thumb io.Reader, originalSize, thumbSize int64, category, originalName, thumbName string) (*DualUploadResult, error) {
	origData, err := io.ReadAll(io.LimitReader(original, originalSize))
	if err != nil {
		return nil, err
	}
	thumbData, err := io.ReadAll(io.LimitReader(thumb, thumbSize))
	if err != nil {
		return nil, err
	}
	return UploadDualImageWithProgress(context.Background(), origData, thumbData, category, originalName, thumbName, nil)
}

// UploadDualImageBytes is a convenience wrapper around UploadDualImage.
func UploadDualImageBytes(original, thumb []byte, category, originalName, thumbName string) (*DualUploadResult, error) {
	return UploadDualImageWithProgress(context.Background(), original, thumb, category, originalName, thumbName, nil)
}

// UploadDualImageWithProgress uploads dual images and reports weighted OSS put progress.
func UploadDualImageWithProgress(
	ctx context.Context,
	original, thumb []byte,
	category, originalName, thumbName string,
	onProgress ProgressFunc,
) (*DualUploadResult, error) {
	originalSize := int64(len(original))
	thumbSize := int64(len(thumb))
	if originalSize <= 0 || originalSize > MaxOriginalBytes {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "原图须 ≤3MB")
	}
	if thumbSize <= 0 || thumbSize > MaxThumbBytes {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "缩略图须 ≤0.5MB")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cfg := config.LoadAppConfig()
	cat := sanitizeCategory(category)
	id := guid.S()
	origExt := extFromName(originalName, ".jpg")
	thumbExt := extFromName(thumbName, ".jpg")
	origKey := cfg.OSSPrefix + cat + "/" + id + "_o" + origExt
	thumbKey := cfg.OSSPrefix + cat + "/" + id + "_t" + thumbExt
	if err := assertMuHouKey(origKey); err != nil {
		return nil, err
	}
	if err := assertMuHouKey(thumbKey); err != nil {
		return nil, err
	}
	bucket, err := getOSSBucket()
	if err != nil {
		return nil, err
	}

	var origRead, thumbRead int64
	total := originalSize + thumbSize
	lastPct := -1
	lastEmit := time.Time{}
	report := func(phase string) {
		if onProgress == nil || total <= 0 {
			return
		}
		pct := int((origRead + thumbRead) * 100 / total)
		if pct > 100 {
			pct = 100
		}
		now := time.Now()
		if pct == lastPct && now.Sub(lastEmit) < 150*time.Millisecond {
			return
		}
		lastPct = pct
		lastEmit = now
		onProgress(pct, phase)
	}

	origReader := &countingReader{
		ctx:    ctx,
		r:      bytes.NewReader(original),
		n:      &origRead,
		report: func() { report("oss_original") },
	}
	if err := bucket.PutObject(origKey, io.LimitReader(origReader, originalSize)); err != nil {
		_ = bucket.DeleteObject(origKey)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "上传原图失败")
	}
	origRead = originalSize
	report("oss_original")

	if err := ctx.Err(); err != nil {
		_ = bucket.DeleteObject(origKey)
		return nil, err
	}

	thumbReader := &countingReader{
		ctx:    ctx,
		r:      bytes.NewReader(thumb),
		n:      &thumbRead,
		report: func() { report("oss_thumb") },
	}
	if err := bucket.PutObject(thumbKey, io.LimitReader(thumbReader, thumbSize)); err != nil {
		_ = bucket.DeleteObject(origKey)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "上传缩略图失败")
	}
	thumbRead = thumbSize
	report("oss_thumb")

	return &DualUploadResult{
		Original: cdnURL(cfg.CDNBase, origKey),
		Thumb:    cdnURL(cfg.CDNBase, thumbKey),
		OrigKey:  origKey,
		ThumbKey: thumbKey,
	}, nil
}

// DeleteObject deletes an object only when its key is under the muhou/ prefix.
func DeleteObject(objectKeyOrURL string) error {
	key, err := resolveObjectKey(objectKeyOrURL)
	if err != nil {
		return err
	}
	if err := assertMuHouKey(key); err != nil {
		return err
	}
	bucket, err := getOSSBucket()
	if err != nil {
		return err
	}
	if err := bucket.DeleteObject(key); err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, "删除对象失败")
	}
	return nil
}

func assertMuHouKey(key string) error {
	cfg := config.LoadAppConfig()
	key = strings.TrimPrefix(key, "/")
	if !strings.HasPrefix(key, cfg.OSSPrefix) {
		return gerror.NewCode(gcode.CodeInvalidOperation, fmt.Sprintf("仅允许删除 %s 前缀对象", cfg.OSSPrefix))
	}
	if strings.Contains(key, "..") {
		return gerror.NewCode(gcode.CodeInvalidParameter, "非法 object key")
	}
	return nil
}

func resolveObjectKey(objectKeyOrURL string) (string, error) {
	cfg := config.LoadAppConfig()
	s := strings.TrimSpace(objectKeyOrURL)
	if s == "" {
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "空 object key")
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		base := cfg.CDNBase + "/"
		if strings.HasPrefix(s, base) {
			return strings.TrimPrefix(s, base), nil
		}
		// fallback: take path after host
		if i := strings.Index(s, "://"); i >= 0 {
			rest := s[i+3:]
			if j := strings.Index(rest, "/"); j >= 0 {
				return rest[j+1:], nil
			}
		}
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "无法解析 CDN URL")
	}
	return strings.TrimPrefix(s, "/"), nil
}

func cdnURL(cdnBase, key string) string {
	return strings.TrimRight(cdnBase, "/") + "/" + strings.TrimPrefix(key, "/")
}

func sanitizeCategory(category string) string {
	c := strings.Trim(strings.ToLower(strings.TrimSpace(category)), "/")
	if c == "" {
		return "misc"
	}
	var b strings.Builder
	for _, r := range c {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '/' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "misc"
	}
	return out
}

func extFromName(name, fallback string) string {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		if ext == ".jpeg" {
			return ".jpg"
		}
		return ext
	default:
		return fallback
	}
}

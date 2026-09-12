package admin

import (
	"encoding/json"
	"io"

	"go-compane-profile/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

// HandleUploadStream streams NDJSON progress while uploading dual images to OSS.
// Registered without MiddlewareHandlerResponse so chunks are not JSON-wrapped.
func HandleUploadStream(r *ghttp.Request) {
	r.Response.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
	r.Response.Header().Set("Cache-Control", "no-cache, no-transform")
	r.Response.Header().Set("X-Accel-Buffering", "no")

	writeEvent := func(v any) {
		b, err := json.Marshal(v)
		if err != nil {
			return
		}
		r.Response.Write(b)
		r.Response.Write([]byte("\n"))
		r.Response.Flush()
	}

	origFile := r.GetUploadFile("original")
	thumbFile := r.GetUploadFile("thumb")
	if origFile == nil || thumbFile == nil {
		r.Response.WriteStatus(400)
		writeEvent(map[string]any{"type": "error", "message": "需要 multipart 字段 original 与 thumb"})
		return
	}
	category := r.Get("category", "misc").String()

	origFH, err := origFile.Open()
	if err != nil {
		r.Response.WriteStatus(400)
		writeEvent(map[string]any{"type": "error", "message": "读取原图失败"})
		return
	}
	defer origFH.Close()
	thumbFH, err := thumbFile.Open()
	if err != nil {
		r.Response.WriteStatus(400)
		writeEvent(map[string]any{"type": "error", "message": "读取缩略图失败"})
		return
	}
	defer thumbFH.Close()

	origData, err := io.ReadAll(io.LimitReader(origFH, service.MaxOriginalBytes+1))
	if err != nil {
		r.Response.WriteStatus(400)
		writeEvent(map[string]any{"type": "error", "message": "读取原图失败"})
		return
	}
	thumbData, err := io.ReadAll(io.LimitReader(thumbFH, service.MaxThumbBytes+1))
	if err != nil {
		r.Response.WriteStatus(400)
		writeEvent(map[string]any{"type": "error", "message": "读取缩略图失败"})
		return
	}

	writeEvent(map[string]any{"type": "progress", "pct": 0, "phase": "oss_original"})

	ctx := r.Context()
	result, err := service.UploadDualImageWithProgress(
		ctx,
		origData,
		thumbData,
		category,
		origFile.Filename,
		thumbFile.Filename,
		func(pct int, phase string) {
			writeEvent(map[string]any{
				"type":  "progress",
				"pct":   pct,
				"phase": phase,
			})
		},
	)
	if err != nil {
		msg := err.Error()
		if ctx.Err() != nil {
			msg = "上传已取消"
		}
		writeEvent(map[string]any{"type": "error", "message": msg})
		return
	}

	writeEvent(map[string]any{
		"type":     "done",
		"original": result.Original,
		"thumb":    result.Thumb,
	})
}

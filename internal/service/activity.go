package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-compane-profile/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ActivityListItem is a public/admin list card.
type ActivityListItem struct {
	Id    uint64  `json:"id"`
	Title string  `json:"title"`
	Image DualURL `json:"image"`
}

// ActivityDetail is public/admin detail payload.
type ActivityDetail struct {
	Id       uint64  `json:"id"`
	Title    string  `json:"title"`
	BodyHtml string  `json:"bodyHtml"`
	Image    DualURL `json:"image"`
	SortOrder int    `json:"sortOrder"`
}

// ActivityInput is create/update body.
type ActivityInput struct {
	Title    string  `json:"title" v:"required#活动标题不能为空"`
	BodyHtml string  `json:"bodyHtml"`
	Image    DualURL `json:"image"`
}

func dualImageEmpty(u DualURL) bool {
	return strings.TrimSpace(u.Thumb) == "" && strings.TrimSpace(u.Original) == ""
}

func requireActivityImage(u DualURL) error {
	if dualImageEmpty(u) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "活动图片不能为空")
	}
	return nil
}

// ListActivities returns all activities in sort order.
func ListActivities(ctx context.Context) ([]ActivityListItem, error) {
	var rows []entity.Activity
	if err := g.DB().Model("activity").Ctx(ctx).OrderAsc("sort_order").OrderAsc("id").Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]ActivityListItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, ActivityListItem{
			Id:    r.Id,
			Title: r.Title,
			Image: DualURL{Thumb: r.ImageThumbUrl, Original: r.ImageOriginalUrl},
		})
	}
	return out, nil
}

// GetActivity returns one activity by numeric id.
func GetActivity(ctx context.Context, id string) (*ActivityDetail, error) {
	row, err := findActivityByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return activityToDetail(row), nil
}

// CreateActivity inserts an activity at the end of the list.
func CreateActivity(ctx context.Context, in ActivityInput) (*ActivityDetail, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "活动标题不能为空")
	}
	if err := requireActivityImage(in.Image); err != nil {
		return nil, err
	}
	body := SanitizeActivityHTML(in.BodyHtml)
	var created *ActivityDetail
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		val, err := tx.Model("activity").Ctx(ctx).Max("sort_order")
		if err != nil {
			return err
		}
		sortOrder := int(val) + 1
		id, err := tx.Model("activity").Ctx(ctx).Data(g.Map{
			"title":              title,
			"body_html":          body,
			"image_original_url": in.Image.Original,
			"image_thumb_url":    in.Image.Thumb,
			"sort_order":         sortOrder,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		row := entity.Activity{}
		if err := tx.Model("activity").Ctx(ctx).Where("id", id).Scan(&row); err != nil {
			return err
		}
		created = activityToDetail(&row)
		return nil
	})
	return created, err
}

// UpdateActivity updates fields; image remains required.
func UpdateActivity(ctx context.Context, id string, in ActivityInput) (*ActivityDetail, error) {
	row, err := findActivityByID(ctx, id)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "活动标题不能为空")
	}
	if err := requireActivityImage(in.Image); err != nil {
		return nil, err
	}
	body := SanitizeActivityHTML(in.BodyHtml)
	oldImage := DualURL{Thumb: row.ImageThumbUrl, Original: row.ImageOriginalUrl}
	_, err = g.DB().Model("activity").Ctx(ctx).Where("id", row.Id).Data(g.Map{
		"title":              title,
		"body_html":          body,
		"image_original_url": in.Image.Original,
		"image_thumb_url":    in.Image.Thumb,
	}).Update()
	if err != nil {
		return nil, err
	}
	deleteReplacedDualsBestEffort(ctx, []DualURL{oldImage}, []DualURL{in.Image})
	return GetActivity(ctx, fmt.Sprintf("%d", row.Id))
}

// DeleteActivity removes one activity and its OSS image best-effort.
func DeleteActivity(ctx context.Context, id string) error {
	row, err := findActivityByID(ctx, id)
	if err != nil {
		return err
	}
	urls := collectDualURLs(nil, DualURL{Thumb: row.ImageThumbUrl, Original: row.ImageOriginalUrl})
	if _, err = g.DB().Model("activity").Ctx(ctx).Where("id", row.Id).Delete(); err != nil {
		return err
	}
	deleteOSSURLsBestEffort(ctx, urls)
	return nil
}

// ReorderActivities sets sort_order from ordered id list.
func ReorderActivities(ctx context.Context, ids []string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for i, id := range ids {
			row, err := findActivityByIDTx(ctx, tx, id)
			if err != nil {
				return err
			}
			if _, err := tx.Model("activity").Ctx(ctx).Where("id", row.Id).Data(g.Map{"sort_order": i}).Update(); err != nil {
				return err
			}
		}
		return nil
	})
}

func findActivityByID(ctx context.Context, id string) (*entity.Activity, error) {
	return findActivityByIDTx(ctx, g.DB(), id)
}

func findActivityByIDTx(ctx context.Context, db dbQuerier, id string) (*entity.Activity, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, gerror.NewCode(gcode.CodeNotFound, "活动不存在")
	}
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "活动不存在")
	}
	var row entity.Activity
	err := db.Model("activity").Ctx(ctx).Where("id", id).Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "活动不存在")
	}
	return &row, nil
}

func activityToDetail(row *entity.Activity) *ActivityDetail {
	return &ActivityDetail{
		Id:        row.Id,
		Title:     row.Title,
		BodyHtml:  row.BodyHtml,
		Image:     DualURL{Thumb: row.ImageThumbUrl, Original: row.ImageOriginalUrl},
		SortOrder: row.SortOrder,
	}
}

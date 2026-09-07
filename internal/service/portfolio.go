package service

import (
	"context"
	"fmt"
	"strings"

	"go-compane-profile/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

// PortfolioListItem is a public list card.
type PortfolioListItem struct {
	Id            string `json:"id"`
	Category      string `json:"category"`
	Cover         string `json:"cover"`
	CoverOriginal string `json:"coverOriginal"`
	Address       string `json:"address"`
	Area          string `json:"area"`
	Style         string `json:"style"`
}

// PortfolioDetail is public detail payload.
type PortfolioDetail struct {
	Id        string    `json:"id"`
	Category  string    `json:"category"`
	Address   string    `json:"address"`
	Area      string    `json:"area"`
	Style     string    `json:"style"`
	HeartFlow string    `json:"heartFlow"`
	Cover     DualURL   `json:"cover"`
	Renders   []DualURL `json:"renders"`
	Reals     []DualURL `json:"reals"`
}

// AdminPortfolio is admin list/detail with numeric id + slug.
type AdminPortfolio struct {
	Id        uint64    `json:"id"`
	Slug      string    `json:"slug"`
	Category  string    `json:"category"`
	SortOrder int       `json:"sortOrder"`
	ViewCount int       `json:"viewCount"`
	Address   string    `json:"address"`
	Area      string    `json:"area"`
	Style     string    `json:"style"`
	HeartFlow string    `json:"heartFlow"`
	Cover     DualURL   `json:"cover"`
	Renders   []DualURL `json:"renders"`
	Reals     []DualURL `json:"reals"`
}

// PortfolioInput is create/update body.
type PortfolioInput struct {
	Slug      string    `json:"slug"`
	Category  string    `json:"category"`
	SortOrder *int      `json:"sortOrder"`
	Address   string    `json:"address"`
	Area      string    `json:"area"`
	Style     string    `json:"style"`
	HeartFlow string    `json:"heartFlow"`
	Cover     DualURL   `json:"cover"`
	Renders   []DualURL `json:"renders"`
	Reals     []DualURL `json:"reals"`
}

// PortfolioListPage is a paginated public portfolio list.
type PortfolioListPage struct {
	List     []PortfolioListItem
	Total    int
	Page     int
	PageSize int
}

const (
	defaultPortfolioPageSize = 10
	maxPortfolioPageSize     = 50
)

// NormalizePortfolioPage clamps page (≥1) and pageSize (default 10, max 50).
func NormalizePortfolioPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPortfolioPageSize
	}
	if pageSize > maxPortfolioPageSize {
		pageSize = maxPortfolioPageSize
	}
	return page, pageSize
}

// ListPortfoliosPublic returns a page of portfolios ordered by sort_order.
// category may be empty (all) or a valid category filter.
func ListPortfoliosPublic(ctx context.Context, page, pageSize int, category string) (*PortfolioListPage, error) {
	page, pageSize = NormalizePortfolioPage(page, pageSize)
	cat, err := NormalizePortfolioCategory(category, true)
	if err != nil {
		return nil, err
	}
	m := g.DB().Model("portfolio").Ctx(ctx)
	if cat != "" {
		m = m.Where("category", cat)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	q := g.DB().Model("portfolio").Ctx(ctx)
	if cat != "" {
		q = q.Where("category", cat)
	}
	var rows []entity.Portfolio
	if err := q.OrderAsc("sort_order").OrderAsc("id").Page(page, pageSize).Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]PortfolioListItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, PortfolioListItem{
			Id:            r.Slug,
			Category:      r.Category,
			Cover:         r.CoverThumbUrl,
			CoverOriginal: r.CoverOriginalUrl,
			Address:       r.Address,
			Area:          r.Area,
			Style:         r.Style,
		})
	}
	return &PortfolioListPage{
		List:     out,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetPortfolioPublic returns detail by slug (public id).
func GetPortfolioPublic(ctx context.Context, id string) (*PortfolioDetail, error) {
	row, err := findPortfolioBySlugOrID(ctx, id)
	if err != nil {
		return nil, err
	}
	renders, reals, err := loadGallery(ctx, row.Id)
	if err != nil {
		return nil, err
	}
	return &PortfolioDetail{
		Id:        row.Slug,
		Category:  row.Category,
		Address:   row.Address,
		Area:      row.Area,
		Style:     row.Style,
		HeartFlow: row.HeartFlow,
		Cover: DualURL{
			Thumb:    row.CoverThumbUrl,
			Original: row.CoverOriginalUrl,
		},
		Renders: renders,
		Reals:   reals,
	}, nil
}

// ListPortfoliosAdmin returns portfolios in a category with galleries.
func ListPortfoliosAdmin(ctx context.Context, category string) ([]AdminPortfolio, error) {
	cat, err := NormalizePortfolioCategory(category, false)
	if err != nil {
		return nil, err
	}
	var rows []entity.Portfolio
	if err := g.DB().Model("portfolio").Ctx(ctx).Where("category", cat).OrderAsc("sort_order").OrderAsc("id").Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]AdminPortfolio, 0, len(rows))
	for _, r := range rows {
		item, err := toAdminPortfolio(ctx, &r)
		if err != nil {
			return nil, err
		}
		out = append(out, *item)
	}
	return out, nil
}

// GetPortfolioAdmin returns one portfolio by numeric id or slug.
func GetPortfolioAdmin(ctx context.Context, id string) (*AdminPortfolio, error) {
	row, err := findPortfolioBySlugOrID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toAdminPortfolio(ctx, row)
}

// CreatePortfolio inserts a portfolio; auto-generates slug when empty.
func CreatePortfolio(ctx context.Context, in PortfolioInput) (*AdminPortfolio, error) {
	cat, err := NormalizePortfolioCategory(in.Category, false)
	if err != nil {
		return nil, err
	}
	var created *AdminPortfolio
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		sortOrder := 0
		if in.SortOrder != nil {
			sortOrder = *in.SortOrder
		} else {
			val, err := tx.Model("portfolio").Ctx(ctx).Where("category", cat).Max("sort_order")
			if err != nil {
				return err
			}
			sortOrder = int(val) + 1
		}
		data := g.Map{
			"category":           cat,
			"sort_order":         sortOrder,
			"address":            in.Address,
			"area":               in.Area,
			"style":              in.Style,
			"heart_flow":         in.HeartFlow,
			"cover_original_url": in.Cover.Original,
			"cover_thumb_url":    in.Cover.Thumb,
		}
		if slug := strings.TrimSpace(in.Slug); slug != "" {
			data["slug"] = slug
		} else {
			data["slug"] = "tmp-" + guid.S()[:12]
		}
		id, err := tx.Model("portfolio").Ctx(ctx).Data(data).InsertAndGetId()
		if err != nil {
			return err
		}
		if strings.TrimSpace(in.Slug) == "" {
			finalSlug := fmt.Sprintf("p-%d", id)
			if _, err := tx.Model("portfolio").Ctx(ctx).Where("id", id).Data(g.Map{"slug": finalSlug}).Update(); err != nil {
				return err
			}
		}
		if err := replaceGalleryTx(ctx, tx, uint64(id), in.Renders, in.Reals); err != nil {
			return err
		}
		row := entity.Portfolio{}
		if err := tx.Model("portfolio").Ctx(ctx).Where("id", id).Scan(&row); err != nil {
			return err
		}
		item, err := toAdminPortfolioTx(ctx, tx, &row)
		if err != nil {
			return err
		}
		created = item
		return nil
	})
	return created, err
}

// UpdatePortfolio updates fields and optional galleries.
func UpdatePortfolio(ctx context.Context, id string, in PortfolioInput) (*AdminPortfolio, error) {
	row, err := findPortfolioBySlugOrID(ctx, id)
	if err != nil {
		return nil, err
	}
	newCat := row.Category
	if strings.TrimSpace(in.Category) != "" {
		newCat, err = NormalizePortfolioCategory(in.Category, false)
		if err != nil {
			return nil, err
		}
	}
	oldCover := DualURL{Thumb: row.CoverThumbUrl, Original: row.CoverOriginalUrl}
	oldRenders, oldReals, err := loadGallery(ctx, row.Id)
	if err != nil {
		return nil, err
	}
	touchGallery := in.Renders != nil || in.Reals != nil
	categoryChanged := newCat != row.Category
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		data := g.Map{
			"address":            in.Address,
			"area":               in.Area,
			"style":              in.Style,
			"heart_flow":         in.HeartFlow,
			"cover_original_url": in.Cover.Original,
			"cover_thumb_url":    in.Cover.Thumb,
		}
		if slug := strings.TrimSpace(in.Slug); slug != "" {
			data["slug"] = slug
		}
		if categoryChanged {
			data["category"] = newCat
			val, err := tx.Model("portfolio").Ctx(ctx).Where("category", newCat).Max("sort_order")
			if err != nil {
				return err
			}
			data["sort_order"] = int(val) + 1
		} else if in.SortOrder != nil {
			data["sort_order"] = *in.SortOrder
		}
		if _, err := tx.Model("portfolio").Ctx(ctx).Where("id", row.Id).Data(data).Update(); err != nil {
			return err
		}
		if categoryChanged {
			if err := resequenceCategoryTx(ctx, tx, row.Category); err != nil {
				return err
			}
		}
		if touchGallery {
			renders := in.Renders
			reals := in.Reals
			if renders == nil {
				renders, _, _ = loadGalleryTx(ctx, tx, row.Id)
			}
			if reals == nil {
				_, reals, _ = loadGalleryTx(ctx, tx, row.Id)
			}
			if err := replaceGalleryTx(ctx, tx, row.Id, renders, reals); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	deleteReplacedDualsBestEffort(ctx, []DualURL{oldCover}, []DualURL{in.Cover})
	if touchGallery {
		renders := in.Renders
		reals := in.Reals
		if renders == nil {
			renders = oldRenders
		}
		if reals == nil {
			reals = oldReals
		}
		oldURLs := collectDualURLLists(oldRenders, oldReals)
		newURLs := collectDualURLLists(renders, reals)
		deleteOSSURLsBestEffort(ctx, urlsNotIn(oldURLs, newURLs))
	}
	return GetPortfolioAdmin(ctx, fmt.Sprintf("%d", row.Id))
}

// DeletePortfolio removes a portfolio and cascaded images.
func DeletePortfolio(ctx context.Context, id string) error {
	row, err := findPortfolioBySlugOrID(ctx, id)
	if err != nil {
		return err
	}
	renders, reals, err := loadGallery(ctx, row.Id)
	if err != nil {
		return err
	}
	urls := collectDualURLs(nil, DualURL{Thumb: row.CoverThumbUrl, Original: row.CoverOriginalUrl})
	urls = append(urls, collectDualURLLists(renders, reals)...)
	if _, err = g.DB().Model("portfolio").Ctx(ctx).Where("id", row.Id).Delete(); err != nil {
		return err
	}
	deleteOSSURLsBestEffort(ctx, urls)
	return nil
}

// ReorderPortfolios sets sort_order by ordered slug/id list within a category.
func ReorderPortfolios(ctx context.Context, category string, ids []string) error {
	cat, err := NormalizePortfolioCategory(category, false)
	if err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for i, id := range ids {
			row, err := findPortfolioBySlugOrIDTx(ctx, tx, id)
			if err != nil {
				return err
			}
			if row.Category != cat {
				return gerror.NewCode(gcode.CodeInvalidParameter, "作品不属于该类")
			}
			if _, err := tx.Model("portfolio").Ctx(ctx).Where("id", row.Id).Data(g.Map{"sort_order": i}).Update(); err != nil {
				return err
			}
		}
		return nil
	})
}

func resequenceCategoryTx(ctx context.Context, tx gdb.TX, category string) error {
	var rows []entity.Portfolio
	if err := tx.Model("portfolio").Ctx(ctx).Where("category", category).OrderAsc("sort_order").OrderAsc("id").Scan(&rows); err != nil {
		return err
	}
	for i, r := range rows {
		if _, err := tx.Model("portfolio").Ctx(ctx).Where("id", r.Id).Data(g.Map{"sort_order": i}).Update(); err != nil {
			return err
		}
	}
	return nil
}

// SavePortfolioGallery replaces renders/reals for a portfolio.
func SavePortfolioGallery(ctx context.Context, id string, renders, reals []DualURL) (*AdminPortfolio, error) {
	row, err := findPortfolioBySlugOrID(ctx, id)
	if err != nil {
		return nil, err
	}
	oldRenders, oldReals, err := loadGallery(ctx, row.Id)
	if err != nil {
		return nil, err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return replaceGalleryTx(ctx, tx, row.Id, renders, reals)
	})
	if err != nil {
		return nil, err
	}
	oldURLs := collectDualURLLists(oldRenders, oldReals)
	newURLs := collectDualURLLists(renders, reals)
	deleteOSSURLsBestEffort(ctx, urlsNotIn(oldURLs, newURLs))
	return GetPortfolioAdmin(ctx, fmt.Sprintf("%d", row.Id))
}

// dbQuerier is satisfied by gdb.DB and gdb.TX.
type dbQuerier interface {
	Model(tableNameOrStruct ...interface{}) *gdb.Model
}

func findPortfolioBySlugOrID(ctx context.Context, id string) (*entity.Portfolio, error) {
	return findPortfolioBySlugOrIDTx(ctx, g.DB(), id)
}

// IncrementPortfolioView atomically increments view_count for a portfolio.
func IncrementPortfolioView(ctx context.Context, id string) error {
	row, err := findPortfolioBySlugOrID(ctx, id)
	if err != nil {
		return err
	}
	_, err = g.DB().Model("portfolio").Ctx(ctx).Where("id", row.Id).Data("view_count=view_count+1").Update()
	return err
}

func findPortfolioBySlugOrIDTx(ctx context.Context, db dbQuerier, id string) (*entity.Portfolio, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, gerror.NewCode(gcode.CodeNotFound, "作品不存在")
	}
	var row entity.Portfolio
	err := db.Model("portfolio").Ctx(ctx).Where("slug", id).WhereOr("id", id).Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "作品不存在")
	}
	return &row, nil
}

func toAdminPortfolio(ctx context.Context, row *entity.Portfolio) (*AdminPortfolio, error) {
	return toAdminPortfolioTx(ctx, g.DB(), row)
}

func toAdminPortfolioTx(ctx context.Context, db dbQuerier, row *entity.Portfolio) (*AdminPortfolio, error) {
	renders, reals, err := loadGalleryTx(ctx, db, row.Id)
	if err != nil {
		return nil, err
	}
	return &AdminPortfolio{
		Id:        row.Id,
		Slug:      row.Slug,
		Category:  row.Category,
		SortOrder: row.SortOrder,
		ViewCount: row.ViewCount,
		Address:   row.Address,
		Area:      row.Area,
		Style:     row.Style,
		HeartFlow: row.HeartFlow,
		Cover: DualURL{
			Thumb:    row.CoverThumbUrl,
			Original: row.CoverOriginalUrl,
		},
		Renders: renders,
		Reals:   reals,
	}, nil
}

func loadGallery(ctx context.Context, portfolioID uint64) (renders, reals []DualURL, err error) {
	return loadGalleryTx(ctx, g.DB(), portfolioID)
}

func loadGalleryTx(ctx context.Context, db dbQuerier, portfolioID uint64) (renders, reals []DualURL, err error) {
	var imgs []entity.PortfolioImage
	if err = db.Model("portfolio_image").Ctx(ctx).
		Where("portfolio_id", portfolioID).
		OrderAsc("kind").OrderAsc("sort_order").OrderAsc("id").
		Scan(&imgs); err != nil {
		return nil, nil, err
	}
	renders = make([]DualURL, 0)
	reals = make([]DualURL, 0)
	for _, img := range imgs {
		pair := DualURL{Thumb: img.ThumbUrl, Original: img.OriginalUrl}
		switch img.Kind {
		case entity.ImageKindRender:
			renders = append(renders, pair)
		case entity.ImageKindReal:
			reals = append(reals, pair)
		}
	}
	return renders, reals, nil
}

func replaceGalleryTx(ctx context.Context, tx dbQuerier, portfolioID uint64, renders, reals []DualURL) error {
	if _, err := tx.Model("portfolio_image").Ctx(ctx).Where("portfolio_id", portfolioID).Delete(); err != nil {
		return err
	}
	rows := make([]g.Map, 0, len(renders)+len(reals))
	for i, u := range renders {
		rows = append(rows, g.Map{
			"portfolio_id": portfolioID,
			"kind":         entity.ImageKindRender,
			"sort_order":   i,
			"original_url": u.Original,
			"thumb_url":    u.Thumb,
		})
	}
	for i, u := range reals {
		rows = append(rows, g.Map{
			"portfolio_id": portfolioID,
			"kind":         entity.ImageKindReal,
			"sort_order":   i,
			"original_url": u.Original,
			"thumb_url":    u.Thumb,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err := tx.Model("portfolio_image").Ctx(ctx).Data(rows).Insert()
	return err
}

package entity

// PortfolioImage is a render or real photo for a portfolio.
type PortfolioImage struct {
	Id          uint64 `json:"id"          orm:"id,primary"`
	PortfolioId uint64 `json:"portfolioId" orm:"portfolio_id"`
	Kind        string `json:"kind"        orm:"kind"` // render | real
	SortOrder   int    `json:"sortOrder"   orm:"sort_order"`
	OriginalUrl string `json:"originalUrl" orm:"original_url"`
	ThumbUrl    string `json:"thumbUrl"    orm:"thumb_url"`
}

// TableName returns the table name for GoFrame ORM.
func (PortfolioImage) TableName() string { return "portfolio_image" }

const (
	ImageKindRender = "render"
	ImageKindReal   = "real"
)

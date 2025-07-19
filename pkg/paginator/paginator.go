package paginator

import (
	"math"
	"reapp/pkg/filterscopes"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit     int         `json:"limit,omitempty" form:"limit"`
	Page      int         `json:"page,omitempty" form:"page"`
	SortBy    string      `json:"sort_by,omitempty" form:"sort_by"`
	SortDir   string      `json:"sort_dir,omitempty" form:"sort_dir"`
	Total     int         `json:"total"`
	TotalPage int         `json:"total_page"`
	Rows      interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 1000 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetSort() string {
	sortField := p.SortBy
	if sortField == "" {
		sortField = "created_at"
	}
	direction := "asc"
	if p.SortDir == "desc" {
		direction = "desc"
	}
	return sortField + " " + direction
}

func (p *Pagination) SetRows(rows interface{}) {
	p.Rows = rows
}

// Paginate returns a GORM scope that applies filtering, pagination, and sorting.
//
// Parameters:
//   - db: the base GORM DB instance.
//   - dbModel: the model type for counting total records (e.g., &User{}).
//   - pg: pointer to a Pagination struct containing page, limit, and sort information.
//   - filters: slice of QueryFilter values to apply as WHERE conditions.
//
// Returns:
//   - a GORM scope function that applies OFFSET, LIMIT, ORDER, and WHERE clauses.
func Paginate(db *gorm.DB, dbModel interface{}, pg *Pagination, filters []filterscopes.QueryFilter) func(db *gorm.DB) *gorm.DB {
	var total int64

	filteredDB := filterscopes.QueryFilterScopes(db.Model(dbModel), filters)
	filteredDB.Count(&total)

	if pg.Page < 1 {
		pg.Page = 1
	}
	pg.Total = int(total)
	pg.TotalPage = int(math.Ceil(float64(total) / float64(pg.GetLimit())))

	return func(db *gorm.DB) *gorm.DB {
		offset := pg.GetOffset()
		return filterscopes.QueryFilterScopes(db, filters).Offset(offset).Limit(pg.GetLimit()).Order(pg.GetSort())
	}
}

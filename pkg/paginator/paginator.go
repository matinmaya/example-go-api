package paginator

import (
	"encoding/json"
	"math"

	"gorm.io/gorm"

	"reapp/pkg/crypto"
	"reapp/pkg/queryfilter"
	"reapp/pkg/services/rediservice"
)

type Pagination struct {
	Limit         int                       `json:"limit,omitempty" form:"limit"`
	Page          int                       `json:"page,omitempty" form:"page"`
	SortBy        string                    `json:"sort_by,omitempty" form:"sort_by"`
	SortDir       string                    `json:"sort_dir,omitempty" form:"sort_dir"`
	Total         int                       `json:"total"`
	TotalPage     int                       `json:"total_page"`
	Filters       []queryfilter.QueryFilter `json:"-"`
	ListCacheKey  string                    `json:"-"`
	CountCacheKey string                    `json:"-"`
	Rows          interface{}               `json:"data"`
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

func Paginate(db *gorm.DB, repositoryNamespace string, dbModel interface{}, pg *Pagination, filters []queryfilter.QueryFilter) func(db *gorm.DB) *gorm.DB {
	var total int64
	if pg.Page < 1 {
		pg.Page = 1
	}

	collectionKey := "count"
	if err := rediservice.CacheOfRepository(repositoryNamespace, collectionKey, pg.GetCountCacheKey(), &total); err != nil {
		filteredDB := queryfilter.QueryFilterScopes(db.Model(dbModel), filters)
		filteredDB.Count(&total)

		rediservice.SetCacheOfRepository(repositoryNamespace, collectionKey, pg.GetCountCacheKey(), total)
	}

	pg.Total = int(total)
	pg.TotalPage = int(math.Ceil(float64(total) / float64(pg.GetLimit())))

	return func(db *gorm.DB) *gorm.DB {
		offset := pg.GetOffset()
		return queryfilter.QueryFilterScopes(db, filters).Offset(offset).Limit(pg.GetLimit()).Order(pg.GetSort())
	}
}

func (p *Pagination) GenerateListKey() (string, error) {
	bFilter, err := json.Marshal(p.Filters)
	if err != nil {
		return "", err
	}

	fields := map[string]interface{}{
		"Limit":   p.GetLimit(),
		"Page":    p.GetPage(),
		"Offset":  p.GetOffset(),
		"Order":   p.GetSort(),
		"Filters": string(bFilter),
	}

	keys, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}

	return crypto.CacheKey(string(keys)), nil
}

func (p *Pagination) GenerateCountKey() (string, error) {
	bFilter, err := json.Marshal(p.Filters)
	if err != nil {
		return "", err
	}

	return crypto.CacheKey(string(bFilter)), nil
}

func (p *Pagination) GetListCacheKey() string {
	if p.ListCacheKey == "" {
		key, err := p.GenerateListKey()
		if err == nil {
			p.ListCacheKey = key
		}
	}

	return p.ListCacheKey
}

func (p *Pagination) GetCountCacheKey() string {
	if p.CountCacheKey == "" {
		key, err := p.GenerateCountKey()
		if err == nil {
			p.CountCacheKey = key
		}
	}

	return p.CountCacheKey
}

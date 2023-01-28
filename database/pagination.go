package database

import (
	"math"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Pagination struct {
	Limit         int    `json:"limit,omitempty"`
	Page          int    `json:"page,omitempty"`
	SortByColumn  string `json:"sort_by_column,omitempty"`
	SortDirection string `json:"sort_direction,omitempty"`
	TotalRows     int64  `json:"total_rows"`
	TotalPages    int    `json:"total_pages"`
	Rows          any    `json:"rows"`
}

func (p *Pagination) getOffset() int {
	return (p.getPage() - 1) * p.getLimit()
}

func (p *Pagination) getLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) getPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Page > p.TotalPages {
		p.Page = p.TotalPages
	}
	return p.Page
}

func (p *Pagination) getSort(dest any) string {
	sortByColumn := "id"
	if p.SortByColumn != "" {
		p.SortByColumn = strings.ToLower(p.SortByColumn)
		s, _ := schema.Parse(dest, &sync.Map{}, schema.NamingStrategy{})
		for _, field := range s.Fields {
			if p.SortByColumn == field.DBName {
				sortByColumn = p.SortByColumn
				break
			}
		}
	}
	p.SortByColumn = sortByColumn

	sortDirection := "desc"
	if p.SortDirection != "" {
		p.SortDirection = strings.ToLower(p.SortDirection)
		switch p.SortDirection {
		case "desc":
			sortDirection = p.SortDirection
		case "asc":
			sortDirection = p.SortDirection
		}
	}
	p.SortDirection = sortDirection

	return p.SortByColumn + " " + p.SortDirection
}

func Paginate(value any, p *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	p.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.Limit)))
	p.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.getOffset()).Limit(p.getLimit()).Order(p.getSort(value))
	}
}

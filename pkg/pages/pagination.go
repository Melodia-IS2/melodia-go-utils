package pages

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type PaginationResult struct {
	Page         uint `json:"page"`
	PageSize     uint `json:"page_size"`
	TotalRecords uint `json:"total_records"`
	TotalPages   uint `json:"total_pages"`
}

type Pagination struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
}

func (p *Pagination) GetURLQuery() string {
	return fmt.Sprintf("page=%d&pagesize=%d", p.Page, p.PageSize)
}

func FromURLQuery(r *http.Request, defaultPageSize uint64, defaultPage uint64) Pagination {
	page, _ := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if page == 0 {
		page = defaultPage
	}
	pageSize, _ := strconv.ParseUint(r.URL.Query().Get("page_size"), 10, 64)
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	return Pagination{
		Page:     uint(page),
		PageSize: uint(pageSize),
	}
}

func (p *Pagination) Validate(maxPageSize uint) error {
	if p.PageSize > maxPageSize {
		return errors.New("page size must be less than or equal to max page size")
	}
	return nil
}

package pages

import "fmt"

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

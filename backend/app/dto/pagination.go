package dto

type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

type PaginationQuery struct {
	Page    int `form:"page"`
	PerPage int `form:"per_page"`
}

func (p *PaginationQuery) Normalize(defaultPerPage int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = defaultPerPage
	}
}

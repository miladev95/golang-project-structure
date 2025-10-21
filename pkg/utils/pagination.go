package utils

import (
	"math"
)

// Pagination represents pagination parameters
type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// NewPagination creates a new pagination object
func NewPagination(page, pageSize int, total int64) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}

// GetOffset calculates offset for database query
// Example: page=2, pageSize=10 -> offset=10
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size as limit
func (p Pagination) GetLimit() int {
	return p.PageSize
}

// IsValid checks if pagination is valid
func (p Pagination) IsValid() bool {
	return p.Page >= 1 && p.PageSize >= 1 && p.PageSize <= 100
}

// HasPreviousPage checks if there's a previous page
func (p Pagination) HasPreviousPage() bool {
	return p.Page > 1
}

// HasNextPage checks if there's a next page
func (p Pagination) HasNextPage() bool {
	return p.Page < p.TotalPage
}

// GetPreviousPage returns the previous page number
func (p Pagination) GetPreviousPage() int {
	if !p.HasPreviousPage() {
		return p.Page
	}
	return p.Page - 1
}

// GetNextPage returns the next page number
func (p Pagination) GetNextPage() int {
	if !p.HasNextPage() {
		return p.Page
	}
	return p.Page + 1
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, pagination Pagination) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}
}
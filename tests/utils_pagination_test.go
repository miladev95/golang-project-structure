package tests

import (
	"testing"

	"github.com/yourusername/yourproject/pkg/utils"
)

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		pageSize      int
		total         int64
		expectedPage  int
		expectedSize  int
		expectedTotal int
	}{
		{
			name:          "valid pagination",
			page:          1,
			pageSize:      10,
			total:         100,
			expectedPage:  1,
			expectedSize:  10,
			expectedTotal: 10,
		},
		{
			name:          "page zero - default to 1",
			page:          0,
			pageSize:      10,
			total:         100,
			expectedPage:  1,
			expectedSize:  10,
			expectedTotal: 10,
		},
		{
			name:          "negative page - default to 1",
			page:          -1,
			pageSize:      10,
			total:         100,
			expectedPage:  1,
			expectedSize:  10,
			expectedTotal: 10,
		},
		{
			name:          "page size zero - default to 10",
			page:          1,
			pageSize:      0,
			total:         100,
			expectedPage:  1,
			expectedSize:  10,
			expectedTotal: 10,
		},
		{
			name:          "page size > 100 - capped at 100",
			page:          1,
			pageSize:      200,
			total:         1000,
			expectedPage:  1,
			expectedSize:  100,
			expectedTotal: 10,
		},
		{
			name:          "page 2",
			page:          2,
			pageSize:      10,
			total:         100,
			expectedPage:  2,
			expectedSize:  10,
			expectedTotal: 10,
		},
		{
			name:          "partial last page",
			page:          1,
			pageSize:      10,
			total:         25,
			expectedPage:  1,
			expectedSize:  10,
			expectedTotal: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.NewPagination(tt.page, tt.pageSize, tt.total)

			if p.Page != tt.expectedPage {
				t.Errorf("Page: got %d, want %d", p.Page, tt.expectedPage)
			}

			if p.PageSize != tt.expectedSize {
				t.Errorf("PageSize: got %d, want %d", p.PageSize, tt.expectedSize)
			}

			if p.TotalPage != tt.expectedTotal {
				t.Errorf("TotalPage: got %d, want %d", p.TotalPage, tt.expectedTotal)
			}
		})
	}
}

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "page 1",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "page 2",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "page 3",
			page:     3,
			pageSize: 20,
			expected: 40,
		},
		{
			name:     "page 5 size 5",
			page:     5,
			pageSize: 5,
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.Pagination{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}

			result := p.GetOffset()
			if result != tt.expected {
				t.Errorf("GetOffset(): got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestGetLimit(t *testing.T) {
	p := utils.Pagination{
		Page:     1,
		PageSize: 25,
	}

	if p.GetLimit() != 25 {
		t.Errorf("GetLimit(): got %d, want 25", p.GetLimit())
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected bool
	}{
		{
			name:     "valid",
			page:     1,
			pageSize: 10,
			expected: true,
		},
		{
			name:     "valid page 5",
			page:     5,
			pageSize: 50,
			expected: true,
		},
		{
			name:     "invalid - page 0",
			page:     0,
			pageSize: 10,
			expected: false,
		},
		{
			name:     "invalid - negative page",
			page:     -1,
			pageSize: 10,
			expected: false,
		},
		{
			name:     "invalid - page size 0",
			page:     1,
			pageSize: 0,
			expected: false,
		},
		{
			name:     "invalid - page size > 100",
			page:     1,
			pageSize: 101,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.Pagination{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}

			result := p.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid(): got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasPreviousPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected bool
	}{
		{
			name:     "page 1 - no previous",
			page:     1,
			expected: false,
		},
		{
			name:     "page 2 - has previous",
			page:     2,
			expected: true,
		},
		{
			name:     "page 5 - has previous",
			page:     5,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.Pagination{Page: tt.page}
			result := p.HasPreviousPage()
			if result != tt.expected {
				t.Errorf("HasPreviousPage(): got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		totalPage   int
		expected    bool
	}{
		{
			name:        "page 1 of 5 - has next",
			page:        1,
			totalPage:   5,
			expected:    true,
		},
		{
			name:        "page 5 of 5 - no next",
			page:        5,
			totalPage:   5,
			expected:    false,
		},
		{
			name:        "page 3 of 10 - has next",
			page:        3,
			totalPage:   10,
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.Pagination{
				Page:      tt.page,
				TotalPage: tt.totalPage,
			}
			result := p.HasNextPage()
			if result != tt.expected {
				t.Errorf("HasNextPage(): got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPreviousPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected int
	}{
		{
			name:     "page 1",
			page:     1,
			expected: 1,
		},
		{
			name:     "page 2",
			page:     2,
			expected: 1,
		},
		{
			name:     "page 5",
			page:     5,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.Pagination{Page: tt.page}
			result := p.GetPreviousPage()
			if result != tt.expected {
				t.Errorf("GetPreviousPage(): got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestGetNextPage(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		totalPage   int
		expected    int
	}{
		{
			name:        "page 1 of 5",
			page:        1,
			totalPage:   5,
			expected:    2,
		},
		{
			name:        "page 5 of 5",
			page:        5,
			totalPage:   5,
			expected:    5,
		},
		{
			name:        "page 3 of 10",
			page:        3,
			totalPage:   10,
			expected:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := utils.Pagination{
				Page:      tt.page,
				TotalPage: tt.totalPage,
			}
			result := p.GetNextPage()
			if result != tt.expected {
				t.Errorf("GetNextPage(): got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	data := []string{"item1", "item2"}
	p := utils.Pagination{
		Page:      1,
		PageSize:  10,
		Total:     2,
		TotalPage: 1,
	}

	response := utils.NewPaginatedResponse(data, p)

	// Check pagination fields
	if response.Pagination.Page != p.Page {
		t.Error("Pagination mismatch")
	}

	if response.Pagination.PageSize != p.PageSize {
		t.Error("Pagination size mismatch")
	}

	// Check data is not nil
	if response.Data == nil {
		t.Error("Data is nil")
	}
}
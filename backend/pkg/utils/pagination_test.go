package utils

import "testing"

func TestCalcTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		pageSize int
		want     int64
	}{
		{"exact division", 100, 20, 5},
		{"with remainder", 101, 20, 6},
		{"single item", 1, 20, 1},
		{"zero items", 0, 20, 0},
		{"page size 1", 5, 1, 5},
		{"large dataset", 10000, 25, 400},
		{"large with remainder", 10001, 25, 401},
		{"zero page size", 100, 0, 0},
		{"negative page size", 100, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcTotalPages(tt.total, tt.pageSize)
			if got != tt.want {
				t.Errorf("CalcTotalPages(%d, %d) = %d, want %d", tt.total, tt.pageSize, got, tt.want)
			}
		})
	}
}

func TestPaginationParamsSkip(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     int64
	}{
		{"first page", 1, 20, 0},
		{"second page", 2, 20, 20},
		{"third page", 3, 10, 20},
		{"large page", 100, 50, 4950},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationParams{Page: tt.page, PageSize: tt.pageSize}
			got := p.Skip()
			if got != tt.want {
				t.Errorf("PaginationParams{Page: %d, PageSize: %d}.Skip() = %d, want %d",
					tt.page, tt.pageSize, got, tt.want)
			}
		})
	}
}

func TestPaginationParamsLimit(t *testing.T) {
	tests := []struct {
		name     string
		pageSize int
		want     int64
	}{
		{"default size", 20, 20},
		{"small size", 5, 5},
		{"max size", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationParams{PageSize: tt.pageSize}
			got := p.Limit()
			if got != tt.want {
				t.Errorf("PaginationParams{PageSize: %d}.Limit() = %d, want %d",
					tt.pageSize, got, tt.want)
			}
		})
	}
}

package service

import "testing"

func TestNormalizePortfolioPage(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{"defaults when zero", 0, 0, 1, 10},
		{"defaults when negative", -3, -1, 1, 10},
		{"keeps valid", 2, 10, 2, 10},
		{"clamps oversized pageSize", 1, 999, 1, 50},
		{"keeps max pageSize", 3, 50, 3, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotSize := NormalizePortfolioPage(tt.page, tt.pageSize)
			if gotPage != tt.wantPage || gotSize != tt.wantPageSize {
				t.Fatalf("NormalizePortfolioPage(%d,%d)=(%d,%d); want (%d,%d)",
					tt.page, tt.pageSize, gotPage, gotSize, tt.wantPage, tt.wantPageSize)
			}
		})
	}
}

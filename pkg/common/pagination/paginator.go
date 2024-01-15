// Package pagination provides pagination functions
package pagination

import (
	"math"
	"strconv"
)

// Paginator struct
type Paginator struct {
	// Limit item per page
	Limit int
	// Offset skip item based on page
	Offset int
	// Total page
	Total int
	// CurretPage
	CurrentPage int
	// PathUrl
	PathUrl string
	// Count item
	Count int64
}

// GetLinks function will map the meta of pagination
func (p *Paginator) GetLinks() map[string]any {
	// Set detault currentpage so its not greater than total
	if p.CurrentPage > p.Total {
		p.CurrentPage = p.Total
	}

	// Set default currentpage so its not less than 1
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}

	// Create next url based on PathUrl
	nextUrl := p.PathUrl + "?page=" + strconv.Itoa(p.CurrentPage+1)
	if p.CurrentPage >= p.Total {
		nextUrl = ""
	}

	// Create previous url based on PathUrl
	prevUrl := p.PathUrl + "?page=" + strconv.Itoa(p.CurrentPage-1)
	if p.CurrentPage <= 1 {
		prevUrl = ""
	}

	// Return links
	return map[string]any{
		"current_page": p.CurrentPage,
		"total_page":   p.Total,
		"limit":        p.Limit,
		"next_url":     nextUrl,
		"previous_url": prevUrl,
	}
}

// Crusor function will map the meta of pagination only previous and next page
func (p *Paginator) Cursor() map[string]any {
	// Set detault currentpage so its not greater than total
	if p.CurrentPage > p.Total {
		p.CurrentPage = p.Total
	}

	// Set default currentpage so its not less than 1
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}

	// Create next url based on PathUrl
	nextUrl := p.PathUrl + "?page=" + strconv.Itoa(p.CurrentPage+1)
	if p.CurrentPage >= p.Total {
		nextUrl = ""
	}

	// Create previous url based on PathUrl
	prevUrl := p.PathUrl + "?page=" + strconv.Itoa(p.CurrentPage-1)
	if p.CurrentPage <= 1 {
		prevUrl = ""
	}

	// Return links
	return map[string]any{
		"next_url":     nextUrl,
		"previous_url": prevUrl,
	}
}

// Create new paginator
func (p *Paginator) Create(limit, page int, path string) {
	// Default limit
	p.Limit = 10
	if limit != 0 {
		p.Limit = limit
	}

	// Default current page
	p.CurrentPage = 1
	if page != 0 {
		p.CurrentPage = page
	}

	// Set default current page is 1 if current page is less than 1
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}

	// Create pagination offset based on current page and limit
	p.Offset = (p.CurrentPage - 1) * p.Limit

	// Set the path url
	p.PathUrl = path
}

// Set total page
func (p *Paginator) SetTotal(total int64) {
	if total == 0 {
		return
	}
	p.Count = total
	p.Total = int(math.Ceil(float64(total) / float64(p.Limit)))
}

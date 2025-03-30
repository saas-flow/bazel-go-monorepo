package pagination

import (
	"encoding/base64"
	"encoding/json"

	"gorm.io/gorm"
)

type PaginationRequest struct {
	Cursor    string              `form:"cursor"`
	Limit     int                 `form:"limit,default=10" validate:"gte=1,lte=250"` // Min 1, Max 250
	SortBy    string              `form:"sort_by"`
	OrderBy   string              `form:"order_by"`
	StartDate string              `form:"start_date"`
	EndDate   string              `form:"end_date"`
	Ranges    map[string][2]int64 `form:"ranges"`
}

type Cursor struct {
	NextCursor     string `json:"next_cursor"`
	PreviousCursor string `json:"previous_cursor"`
}

func EncodeCursor(data Cursor) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeCursor(data string) (*Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var cursor Cursor
	if err := json.Unmarshal(b, &cursor); err != nil {
		return nil, err
	}

	return &cursor, nil
}

// Function apply filter ke query
func ApplyPaginationAndFilter(p *PaginationRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		// Apply filter date range
		if p.StartDate != "" {
			db = db.Where("created_at >= ?", p.StartDate)
		}

		if p.EndDate != "" {
			db = db.Where("created_at <= ?", p.EndDate)
		}

		// Apply range filters secara otomatis
		for field, rangeVal := range p.Ranges {
			db = db.Where(field+" BETWEEN ? AND ?", rangeVal[0], rangeVal[1])
		}

		// Sorting
		if p.SortBy != "" {
			order := "ASC"
			if p.OrderBy == "desc" {
				order = "DESC"
			}
			db = db.Order(p.SortBy + " " + order)
		}

		// Pagination
		if p.Limit > 0 {
			db = db.Limit(p.Limit + 1)
		}

		// Cursor
		if p.Cursor != "" {
			db = db.Where("created_at < ?", p.Cursor)
		}

		return db
	}
}

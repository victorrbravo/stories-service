package types

// PageBase the base page object
type PageBase struct {
	ID         Int `json:"id"`
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

// Page a full page object with knowledge about next/previous pages and Data
type Page struct {
	PageBase

	PageCount    int           `json:"page_count"`
	PreviousPage Int           `json:"previous_page"`
	NextPage     Int           `json:"next_page"`
	Data         []interface{} `json:"data"`
}

// EmptyPage a special empty page
var EmptyPage = Page{}

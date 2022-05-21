package dns

type PagedRequest struct {
	Page    int
	PerPage int
}

type PageMetadata struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	LastPage     int `json:"last_page"`
	TotalEntries int `json:"total_entries"`
}

type PagedMetadata struct {
	Pagination PageMetadata `json:"pagination"`
}

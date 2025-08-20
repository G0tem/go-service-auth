package types

type PaginationRequest struct {
	Page     int    `query:"page" json:"page"`
	PageSize int    `query:"page_size" json:"page_size"`
	SortType string `query:"sort_type" json:"sort_type" enums:"normal,recommended"`
}

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

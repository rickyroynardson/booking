package entity

type PaginationMetadata struct {
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

type PaginationRequest struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type PaginationResponse struct {
	Data     any                `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

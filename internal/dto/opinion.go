package dto

type CreateOpinionRequest struct {
	Content string `json:"content" binding:"required"`
}

// data transfer object

package domain

type PaginationRequest struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1,max=100"`
}

type CreateAuthorRequest struct {
	Name        string `json:"name" binding:"required"`
	Bio         string `json:"bio" binding:"required"`
	CreatedByID uint   `json:"createdById" binding:"required"`
}

type UpdateAuthorRequest struct {
	ID          uint    `json:"id" binding:"required"`
	Name        *string `json:"name" binding:"required"`
	Bio         *string `json:"bio" binding:"omitempty"`
	UpdatedByID uint    `json:"updatedById" binding:"required"`
}

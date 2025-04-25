package domain

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	ISBN        string `json:"isbn" binding:"required"`
	PublishYear int    `json:"publishYear" binding:"required"`
	Stock       int    `json:"stock" binding:"required"`
	CreatedByID uint   `json:"createdById" binding:"required"`
	AuthorIDs   []uint `json:"authorIds" binding:"required"`
	CategoryIDs []uint `json:"categoryIds" binding:"required"`
}

type UpdateBookRequest struct {
	ID          uint    `json:"id" binding:"required"`
	Title       *string `json:"title" binding:"omitempty"`
	ISBN        *string `json:"isbn" binding:"omitempty"`
	PublishYear *int    `json:"publishYear" binding:"omitempty"`
	Stock       *int    `json:"stock" binding:"omitempty"`
	UpdatedByID uint    `json:"updatedById" binding:"required"`
	AuthorIDs   *[]uint `json:"authorIds" binding:"omitempty"`
	CategoryIDs *[]uint `json:"categoryIds" binding:"omitempty"`
}

type IncreaseStockRequest struct {
	ID     uint `json:"id" binding:"required"`
	Amount int  `json:"amount" binding:"required,gt=0"`
}

type DecreaseStockRequest struct {
	ID     uint `json:"id" binding:"required"`
	Amount int  `json:"amount" binding:"required,gt=0"`
}

type PaginationRequest struct {
	Page   int    `form:"page" binding:"required,min=1"`
	Limit  int    `form:"limit" binding:"required,min=1,max=100"`
	Search string `form:"search"`
}

type DeleteCategoryRequest struct {
	ID uint `json:"id" binding:"required"`
}

type BorrowBookRequest struct {
	UserID      uint `json:"userId" binding:"required"`
	BookID      uint `json:"bookId" binding:"required"`
	CreatedByID uint `json:"createdById" binding:"required"`
}

type ReturnBookRequest struct {
	ID          uint `json:"id" binding:"required"`
	UserID      uint `json:"userId" binding:"required"`
	BookID      uint `json:"bookId" binding:"required"`
	UpdatedByID uint `json:"updatedById" binding:"required"`
}

type CreateRecommendationRequest struct {
	UserID uint   `json:"userId" binding:"required"`
	BookID uint   `json:"bookId" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

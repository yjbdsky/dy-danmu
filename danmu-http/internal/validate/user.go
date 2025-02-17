package validate

type UserPageRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=500"`
}

type UserSearchRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=500"`
	Keyword  string `form:"keyword" binding:"omitempty"`
}

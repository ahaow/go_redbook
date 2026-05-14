package req

type CreateArticleRequest struct {
	Title    string   `json:"title" binding:"required,min=2,max=100"`
	Content  string   `json:"content" binding:"required,min=10"`
	Location string   `json:"location" binding:"max=255"`
	IsPublic bool     `json:"is_public"`
	Images   []string `json:"images" binding:"required,min=1,max=9"`
	Topics   []string `json:"topics" binding:"max=5"`
}

type UpdateArticleRequest struct {
	Title    *string   `json:"title" binding:"omitempty,min=2,max=100"`
	Content  *string   `json:"content" binding:"omitempty,min=10"`
	Location *string   `json:"location" binding:"omitempty,max=255"`
	IsPublic *bool     `json:"is_public"`
	Images   *[]string `json:"images" binding:"omitempty,max=9"`
	Topics   *[]string `json:"topics" binding:"omitempty,max=5"`
}

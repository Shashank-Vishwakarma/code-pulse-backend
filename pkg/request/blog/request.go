package blog

type CommentRequest struct {
	Body string `json:"body" validate:"required"`
}
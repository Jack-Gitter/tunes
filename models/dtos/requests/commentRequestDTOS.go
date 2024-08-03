package requests

type CreateCommentDTO struct {
    CommentText string
}

type CommentIDPathParams struct {
    CommentID int `uri:"commentID" binding:"required,numeric"`
}

type UpdateCommentDTO struct {
    CommentText *string
}


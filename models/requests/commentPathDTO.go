package requests

type CommentIDPathParams struct {
    CommentID int `uri:"commentID" binding:"required,numeric"`
}


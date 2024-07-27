package requests

type UpdatePostRequestDTO struct {
    Rating *int `binding:"gte=0,lte=5"`
	Text   *string  
}

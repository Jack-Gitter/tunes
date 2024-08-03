package requests

type UpdatePostRequestDTO struct {
    Rating *int 
	Review   *string  
}

type CreatePostDTO struct {
	SongID *string
	Rating *int
	Text   *string
}

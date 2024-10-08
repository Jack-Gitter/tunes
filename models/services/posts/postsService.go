package posts

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"time"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/Jack-Gitter/tunes/models/daos"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/Jack-Gitter/tunes/models/services/rabbitmqservice"
	"github.com/Jack-Gitter/tunes/models/services/spotify"
	"github.com/gin-gonic/gin"
)

type PostsService struct {
    DB *sql.DB
    PostsDAO daos.IPostsDAO
    UsersDAO daos.IUsersDAO
    CommentsDAO daos.CommentsDAO
    SpotifyService spotify.ISpotifyService
    RabbitMQService rabbitmqservice.IRabbitMQService
}

type IPostsService interface {
    CreatePostForCurrentUser(c *gin.Context) 
    LikePost(c *gin.Context)
    DislikePost(c *gin.Context)
    GetAllPostsForUserByID(c *gin.Context) 
    GetAllPostsForCurrentUser(c *gin.Context) 
    GetPostBySpotifyIDAndSongID(c *gin.Context)
    GetPostCurrentUserBySongID(c *gin.Context) 
    DeletePostBySpotifyIDAndSongID(c *gin.Context)  
    DeletePostForCurrentUserBySongID(c *gin.Context) 
    UpdateCurrentUserPost(c *gin.Context) // yes
    RemovePostVote(c *gin.Context) 
    GetPostCommentsPaginated(c *gin.Context) 
    GetCurrentUserFeed(c *gin.Context) 
}

// @Summary Creates a post for the current user
// @Description Creates a post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param createPostDTO body requests.CreatePostDTO true "Information required to create a post"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 409 {string} string
// @Failure 500 {string} string
// @Router /posts/ [post]
// @Security Bearer
func(p *PostsService) CreatePostForCurrentUser(c *gin.Context) {

	spotifyID, spotifyIDExists := c.Get("spotifyID")
	spotifyUsername, spotifyUsernameExists := c.Get("spotifyUsername")
	spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")

	if !spotifyIDExists || !spotifyAccessTokenExists || !spotifyUsernameExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}

	createPostDTO := &requests.CreatePostDTO{}
	c.ShouldBindBodyWithJSON(createPostDTO)

	spotifySongResponse, err := p.SpotifyService.GetSongDetailsFromSpotify(*createPostDTO.SongID, spotifyAccessToken.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var albumImage string = ""
	if len(spotifySongResponse.Album.Images) > 0 {
		albumImage = spotifySongResponse.Album.Images[0].Url
	}

    if createPostDTO.Rating == nil {
        rating := 0
        createPostDTO.Rating = &rating
    }
    if createPostDTO.Text == nil {
        text := ""
        createPostDTO.Text = &text
    }

	resp, err := p.PostsDAO.CreatePost(
        p.DB,
		spotifyID.(string),
		*createPostDTO.SongID,
		spotifySongResponse.Name,
		spotifySongResponse.Album.Id,
		spotifySongResponse.Album.Name,
		albumImage,
		*createPostDTO.Rating,
		*createPostDTO.Text,
		time.Now().UTC(),
		spotifyUsername.(string),
	)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

    rabbitMQMessage := rabbitmqservice.RabbitMQPostMessage{
        Type: rabbitmqservice.POST,
        Poster: spotifyID.(string),
    }

    err = p.RabbitMQService.Enqueue(rabbitMQMessage)

	resp.Likes = []responses.UserIdentifer{}
	resp.Dislikes = []responses.UserIdentifer{}

	c.JSON(http.StatusOK, resp)

}

// @Summary Likes a post for the current user
// @Description Likes a post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "Song ID of the post to like"
// @Param songID path string true "Spotify ID of the user who posted the song"
// @Success 204
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 409 {string} string 
// @Failure 500 {string} string 
// @Router /posts/likes/{spotifyID}/{songID} [post]
// @Security Bearer
func(p *PostsService) LikePost(c *gin.Context) {
	currentUserSpotifyID, found := c.Get("spotifyID")
	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()
        
        likes, _, err := p.PostsDAO.GetPostVotes(tx, songID, spotifyID)

        if err != nil {
            return err
        }

        for _, userIdentifier := range likes {
            if userIdentifier.SpotifyID == spotifyID {
                return &customerrors.CustomError{StatusCode: http.StatusConflict, Msg: "cannot like a post twice"}
            }
        }

        err = p.PostsDAO.LikePost(tx, currentUserSpotifyID.(string), spotifyID, songID)

        if err != nil {
            return err
        }

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil

    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

	c.Status(http.StatusNoContent)
}

// @Summary Dislikes a post for the current user
// @Description Dislikes a post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "Song ID of the post to dislike"
// @Param songID path string true "Spotify ID of the user who posted the song"
// @Success 204
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 409 {string} string 
// @Failure 500 {string} string 
// @Router /posts/dislikes/{spotifyID}/{songID} [post]
// @Security Bearer
func(p *PostsService) DislikePost(c *gin.Context) {

	currentUserSpotifyID, found := c.Get("spotifyID")
	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}
    
    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()
        
        _, dislikes, err := p.PostsDAO.GetPostVotes(tx, songID, spotifyID)

        if err != nil {
            return err
        }

        for _, userIdentifier := range dislikes {
            if userIdentifier.SpotifyID == spotifyID {
                return &customerrors.CustomError{StatusCode: http.StatusConflict, Msg: "cannot like a post twice"}
            }
        }

        err = p.PostsDAO.DislikePost(tx, currentUserSpotifyID.(string), spotifyID, songID)

        if err != nil {
            return err
        }

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil

    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

	c.Status(http.StatusNoContent)
}

// @Summary Get all of a users post previews
// @Description Get all of a users post previews
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The user whos posts are recieved. Value is a spotify ID"
// @Param createdAt query string false "Pagination Key. Format is UTC timestamp"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/previews/users/{spotifyID} [get]
// @Security Bearer
func(p *PostsService) GetAllPostsForUserByID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")
	createdAt := c.Query("createdAt")

	var t time.Time = time.Now().UTC()
    var err error

	if createdAt != "" {
        t, err = time.Parse(time.RFC3339, createdAt)

        if err != nil {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid time format"})
            c.Abort()
            return
        }

	}

    paginationResponse := responses.PaginationResponse[[]responses.PostPreview, time.Time]{PaginationKey: time.Now().UTC()}

    tx, err := p.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    _, err = p.UsersDAO.GetUser(tx, spotifyID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    posts, err := p.PostsDAO.GetUserPostsProperties(tx, spotifyID, t)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    for i := 0; i < len(posts); i++ {
        likes, dislikes, err := p.PostsDAO.GetPostVotes(tx, posts[i].SongID, spotifyID)
        if err != nil {
            c.Error(err)
            c.Abort()
            return
        }
        posts[i].Likes = likes
        posts[i].Dislikes = dislikes
    }

    paginationResponse.DataResponse = posts
    if len(posts) > 0 {
        paginationResponse.PaginationKey = posts[len(posts)-1].CreatedAt
    }

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Get all of a users post previews
// @Description Get all of a users post previews
// @Tags Posts
// @Accept json
// @Produce json
// @Param createdAt query string false "Pagination Key. Format is UTC timestamp"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 500 {string} string 
// @Router /posts/previews/users/current [get]
// @Security Bearer
func(p *PostsService) GetAllPostsForCurrentUser(c *gin.Context) {
	spotifyID, spotifyIDExists := c.Get("spotifyID")
	createdAt := c.Query("createdAt")

	if !spotifyIDExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt lookup"})
		c.Abort()
		return
	}

	var t time.Time = time.Now().UTC()
    var err error

	if createdAt != "" {
        t, err = time.Parse(time.RFC3339, createdAt)

        if err != nil {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid time format"})
            c.Abort()
            return
        }

	}

    paginationResponse := responses.PaginationResponse[[]responses.PostPreview, time.Time]{PaginationKey: time.Now().UTC()}

    tx, err := p.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }
    _, err = p.UsersDAO.GetUser(tx, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    posts, err := p.PostsDAO.GetUserPostsProperties(tx, spotifyID.(string), t)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    for i := 0; i < len(posts); i++ {
        likes, dislikes, err := p.PostsDAO.GetPostVotes(tx, posts[i].SongID, spotifyID.(string))
        if err != nil {
            c.Error(err)
            c.Abort()
            return
        }
        posts[i].Likes = likes
        posts[i].Dislikes = dislikes
    }

    paginationResponse.DataResponse = posts
    if len(posts) > 0 {
        paginationResponse.PaginationKey = posts[len(posts)-1].CreatedAt
    }

    err = tx.Commit()

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, paginationResponse)
}
// @Summary Get apath specific post
// @Description Get a specific post
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The user who posted the song"
// @Param songID path string true "The songID of the posted song"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/{spotifyID}/{songID} [get]
// @Security Bearer
func(p *PostsService) GetPostBySpotifyIDAndSongID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

    tx, err := p.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    post, err := p.PostsDAO.GetPostProperties(tx, songID, spotifyID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    likes, dislikes, err := p.PostsDAO.GetPostVotes(tx, songID, spotifyID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    post.Likes = likes
    post.Dislikes = dislikes

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }


	c.JSON(http.StatusOK, post)
}

// @Summary Get a specific post for the current user
// @Description Get a specific post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Success 200 {object} responses.PostPreview
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [get]
// @Security Bearer
func(p *PostsService) GetPostCurrentUserBySongID(c *gin.Context) {

	currentUserSpotifyID, found := c.Get("spotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuckin jwt"})
		c.Abort()
		return
	}

    tx, err := p.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    post, err := p.PostsDAO.GetPostProperties(tx, songID, currentUserSpotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    likes, dislikes, err := p.PostsDAO.GetPostVotes(tx, songID, currentUserSpotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    post.Likes = likes
    post.Dislikes = dislikes

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }


	c.JSON(http.StatusOK, post)
}

// @Summary Deletes a specific post. Only accessible to admins
// @Description Deletes a specific post. Only accessible to admins
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The spotify ID of the user who posted the song"
// @Param songID path string true "The songID of the posted song"
// @Success 204
// @Failure 400 {string} string 
// @Failure 403 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/admin/{spotifyID}/{songID} [delete]
// @Security Bearer
func(p *PostsService) DeletePostBySpotifyIDAndSongID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	err := p.PostsDAO.DeletePost(p.DB, songID, spotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Deletes a post made by the current user
// @Description Deletes a post made by the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Success 204
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [delete]
// @Security Bearer
func(p *PostsService) DeletePostForCurrentUserBySongID(c *gin.Context) {

	requestorSpotifyID, found := c.Get("spotifyID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}
	songID := c.Param("songID")

	err := p.PostsDAO.DeletePost(p.DB, songID, requestorSpotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Updates a post made by the current user
// @Description Updates a post made by the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Param UpdatePostDTO body requests.UpdatePostRequestDTO true "The fields to update"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [patch]
// @Security Bearer
func(p *PostsService) UpdateCurrentUserPost(c *gin.Context) {

	spotifyID, exists := c.Get("spotifyID")
	spotifyUsername, uexists := c.Get("spotifyUsername")
	songID := c.Param("songID")
	updatePostReq := &requests.UpdatePostRequestDTO{}

	c.ShouldBindBodyWithJSON(updatePostReq)

	if !exists || !uexists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "no jwt"})
		c.Abort()
		return
	}

    post := &responses.PostPreview{}

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

        if err != nil {
            return err
        }

        post, err = p.PostsDAO.UpdatePost(tx, spotifyID.(string), songID, updatePostReq, spotifyUsername.(string))

        if err != nil {
            return err
        }

        likes, dislikes, err := p.PostsDAO.GetPostVotes(tx, post.SongID, spotifyID.(string))

        if err != nil {
            return err
        }

        post.Likes = likes
        post.Dislikes = dislikes

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil

    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

	c.JSON(http.StatusOK, post)
}

// @Summary Removes a vote for the current user on a post
// @Description Removes a vote for the current user on a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Param posterSpotifyID path string true "The user who posted the post spotify ID"
// @Success 204
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/votes/current/{posterSpotifyID}/{songID} [delete]
// @Security Bearer
func(p *PostsService) RemovePostVote(c *gin.Context) {
	voterSpotifyID, found := c.Get("spotifyID")
	posterSpotifyID := c.Param("posterSpotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "forgot to set JWT"})
	}

	err := p.PostsDAO.RemovePostVote(p.DB, voterSpotifyID.(string), posterSpotifyID, songID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Gets the comments of a post
// @Description Gets the comments of a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Param spotifyID path string true "The user who posted the post spotify ID"
// @Param createdAt query string false "Pagination Key. In the form of UTC timestamp"
// @Success 200 {object} responses.PaginationResponse[[]responses.Comment, time.Time]
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/comments/{spotifyID}/{songID} [get]
// @Security Bearer
func(p *PostsService) GetPostCommentsPaginated(c *gin.Context) {
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")
    createdAt := c.Query("createdAt")

    var t time.Time = time.Now().UTC()
    var err error

	if createdAt != "" {
		t, err = time.Parse(time.RFC3339, createdAt)

        if err != nil {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid time format"})
            c.Abort()
            return
        }
	}

    tx, err := p.DB.BeginTx(context.Background(), nil)

    paginatedComments := responses.PaginationResponse[[]responses.Comment, time.Time]{PaginationKey: time.Now().UTC()}

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    _, err = p.PostsDAO.GetPostProperties(tx, songID, spotifyID)
    
    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    comments, err := p.PostsDAO.GetPostComments(tx, spotifyID, songID, t)

    for i := 0; i < len(comments); i++ {
        likes, dislikes, err := p.CommentsDAO.GetCommentVotes(tx, fmt.Sprint(comments[i].CommentID))

        if err != nil {
            c.Error(err)
            c.Abort()
            return
        }

        comments[i].Likes = len(likes)
        comments[i].Dislikes = len(dislikes)

        if i == len(comments)-1 {
            paginatedComments.PaginationKey = comments[i].CreatedAt
        }

    }

    paginatedComments.DataResponse = comments

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, paginatedComments)
}


// @Summary Gets the comments of a post
// @Description Gets the comments of a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param createdAt query string false "Pagination Key. In the form of UTC timestamp"
// @Success 200 {object} responses.PaginationResponse[[]responses.PostPreview, time.Time]
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/feed [get]
// @Security Bearer
func(p *PostsService) GetCurrentUserFeed(c *gin.Context) {
    spotifyID, exists := c.Get("spotifyID")
    createdAt := c.Query("createdAt")

	var t time.Time = time.Now().UTC()
    var err error

	if createdAt != "" {
        t, err = time.Parse(time.RFC3339, createdAt)

        if err != nil {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid time format"})
            c.Abort()
            return
        }
	}

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "failed to set JWT"})
        c.Abort()
        return
    }
    
    tx, err := p.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    following, err := p.UsersDAO.GetAllUserFollowing(tx, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    posts := []responses.PostPreview{}

    for _, user_followed := range following {

        user_posts, err := p.PostsDAO.GetUserPostsProperties(tx, user_followed.SpotifyID, t)

        if err != nil {
            c.Error(err)
            c.Abort()
            return
        }

        posts = append(posts, user_posts...)

    }

    slices.SortFunc[[]responses.PostPreview](posts, func(p1, p2 responses.PostPreview) int {
        if p1.CreatedAt.After(p2.CreatedAt) {
            return -1
        }
        return 1
    })

    feedLength := 25

    if len(posts) < feedLength {
        feedLength = len(posts)
    }

    posts = posts[:feedLength]

    for i := 0; i < len(posts); i++ {
        likes, dislikes, err := p.PostsDAO.GetPostVotes(tx, posts[i].SongID, posts[i].SpotifyID)

        if err != nil {
            c.Error(err)
            c.Abort()
            return
        }

        posts[i].Likes = likes
        posts[i].Dislikes = dislikes

    }

    paginationResponse := responses.PaginationResponse[[]responses.PostPreview, time.Time]{DataResponse: posts}
    paginationKey := time.Now().UTC()

    if feedLength > 0 {
        paginationKey = posts[feedLength-1].CreatedAt
    }

    paginationResponse.PaginationKey = paginationKey

    c.JSON(http.StatusOK, paginationResponse)

}

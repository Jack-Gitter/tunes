package comments

import "github.com/gin-gonic/gin"

func CreateComment(c *gin.Context) {

    c.Get("spotifyID")
    c.Param("spotifyID")
    c.Param("songID")




}

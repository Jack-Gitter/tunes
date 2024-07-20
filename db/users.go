package db

import (
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	//"github.com/mitchellh/mapstructure"
)

/* =================== CREATE ================== */

func InsertUserIntoDBIfNeeded(username string, spotifyID string, role responses.Role) (*responses.User, error) {
    return nil, nil
}

/* =================== READ ================== */

func GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, bool, error) {
    return nil, false, nil
}


/* PROPERTY UPDATES */
func UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, bool, error) { 
    return nil, false, nil
}

func DeleteUserByID(spotifyID string) (bool, error) {
    return false, nil
}

func GetFollowers(spotifyID string, paginationKey string) (*responses.PaginationResponse[[]responses.User, string], bool, error) {
    return nil, false, nil
}

func UnfollowUser(spotifyID string, otherUserSpotifyID string) (bool, error)  {
    return false, nil
}

/* RELATIONAL UDPATES */
func FollowUser(spotifyID string, otherUserSpotifyID string) (bool, error) {
    return false, nil
}




func UnFollowUserBySpotifyID(){}

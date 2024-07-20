package db

import (
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	_ "github.com/lib/pq"
)

/* =================== CREATE ================== */

func UpsertUser(username string, spotifyID string, role responses.Role) (*responses.User, error) {
    query := "INSERT INTO users (spotifyid, username, userRole) values ($1, $2, $3) ON CONFLICT (spotifyID) DO UPDATE SET username=$2, userRole=$3 RETURNING bio"
    row := DB.Driver.QueryRow(query, spotifyID, username, role)

    err := row.Err()
    if err != nil {
        return nil, err
    }

    userResponse := &responses.User{}
    userResponse.Role = role
    userResponse.Username = username
    userResponse.SpotifyID = spotifyID
    row.Scan(&userResponse.Bio)


    return userResponse, nil
}

/* =================== READ ================== */

func GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, bool, error) {
    query := "SELECT spotifyid, userrole, username, bio FROM users WHERE spotifyid = $1"
    row := DB.Driver.QueryRow(query, spotifyID)

    err := row.Err()
    if err != nil {
        return nil, false, err
    }

    userResponse := &responses.User{}
    row.Scan(&userResponse.SpotifyID, &userResponse.Username, &userResponse.Role, &userResponse.Bio)

    return userResponse, true, nil
}


/* PROPERTY UPDATES */
func UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, bool, error) { 
    return nil, false, nil
}

func DeleteUserByID(spotifyID string) (bool, error) {
    query := "DELETE FROM users WHERE spotifyID = $1"
    res, err := DB.Driver.Exec(query, spotifyID)

    if err != nil {
        return false, err
    }

    num, err :=  res.RowsAffected()

    if err != nil {
        return false, err
    }

    if num < 1 {
        return false, nil
    }

    return true, nil
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

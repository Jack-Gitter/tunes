package db

import (
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	_ "github.com/lib/pq"
)

/* =================== CREATE ================== */

func UpsertUser(username string, spotifyID string) (*responses.User, error) {
    query := "INSERT INTO users (spotifyid, username, userrole) values ($1, $2, 'BASIC') ON CONFLICT (spotifyID) DO UPDATE SET username=$2 RETURNING bio, userrole"
    row := DB.Driver.QueryRow(query, spotifyID, username)

    err := row.Err()
    if err != nil {
        return nil, err
    }

    userResponse := &responses.User{}
    userResponse.Username = username
    userResponse.SpotifyID = spotifyID
    row.Scan(&userResponse.Bio, &userResponse.Role)


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
// todo
func UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (bool, error) { 
    query := "UPDATE users SET "
    args := []any{}
    if updatedUser.Bio != nil {

        args = append(args, updatedUser.Bio)
        query += "bio = $1 "

    } 
    if updatedUser.Role != nil {
        args = append(args, updatedUser.Role)
        query += ", userrole = $2 "
    }

    //query += "WHERE spotifyID "


    res, err := DB.Driver.Exec(query, args...)

    if err != nil {
        return false, err
    }

    num, err := res.RowsAffected()

    if err != nil {
        return false, err
    }

    if num < 1 {
        return false, nil
    }

    return true, nil
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



/* RELATIONAL UDPATES */
func UnfollowUser(spotifyID string, otherUserSpotifyID string) (bool, error)  {
    return false, nil
}

func FollowUser(spotifyID string, otherUserSpotifyID string) (bool, error) {
    return false, nil
}


func UnFollowUserBySpotifyID(){}

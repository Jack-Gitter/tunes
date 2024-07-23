package db

import (
	"database/sql"
	"fmt"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	_ "github.com/lib/pq"
)

/* =================== CREATE ================== */

func UpsertUser(username string, spotifyID string) (*responses.User, error) {
    query := "INSERT INTO users (spotifyid, username, userrole) values ($1, $2, 'BASIC') ON CONFLICT (spotifyID) DO UPDATE SET username=$2 RETURNING bio, userrole"
    row := DB.Driver.QueryRow(query, spotifyID, username)

    userResponse := &responses.User{}
    userResponse.Username = username
    userResponse.SpotifyID = spotifyID
    err := row.Scan(&userResponse.Bio, &userResponse.Role)

    if err != nil {
        return nil, err
    }

    return userResponse, nil
}

/* =================== READ ================== */

func GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, bool, error) {
    query := "SELECT spotifyid, userrole, username, bio FROM users WHERE spotifyid = $1"
    row := DB.Driver.QueryRow(query, spotifyID)


    userResponse := &responses.User{}
    err := row.Scan(&userResponse.SpotifyID, &userResponse.Username, &userResponse.Role, &userResponse.Bio)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, nil
        } 
        return nil, false, err
    }

    return userResponse, true, nil
}


/* PROPERTY UPDATES */
func UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, bool, error) { 
    query := "UPDATE users SET "
    args := []any{}
    varNum := 1
    if updatedUser.Bio != nil {
        args = append(args, updatedUser.Bio)
        query += fmt.Sprintf("bio = $%d", varNum)
        varNum += 1 

    } 
    if updatedUser.Role != nil {
        args = append(args, updatedUser.Role)
        if varNum > 1 {
            query += fmt.Sprintf(", userrole = $%d", varNum)
        } else {
            query += fmt.Sprintf("userrole = $%d", varNum)
        }
        varNum += 1
    }

    query += fmt.Sprintf(" WHERE spotifyID = $%d RETURNING bio, userrole, spotifyid, username", varNum)
    args = append(args, spotifyID)


    res := DB.Driver.QueryRow(query, args...)

    userResponse := &responses.User{}
    err := res.Scan(&userResponse.Bio, &userResponse.Role, &userResponse.SpotifyID, &userResponse.Username)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, nil 
        } 
        return nil, false, err
    }

    return userResponse, true, nil

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


/* RELATIONAL UDPATES */
func UnfollowUser(spotifyID string, otherUserSpotifyID string) (bool, error)  {
    query := "DELETE FROM followers WHERE follower = $1 AND userfollowed = $2"

    res, err := DB.Driver.Exec(query, spotifyID, otherUserSpotifyID)

    if err != nil {
        return false, err
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return false, err
    }

    if rows < 1 {
        return false, nil
    }

    return true, nil
}

func FollowUser(spotifyID string, otherUserSpotifyID string) (bool, error) {
    query := "INSERT INTO followers (follower, userFollowed) VALUES ($1, $2)"

    res, err := DB.Driver.Exec(query, spotifyID, otherUserSpotifyID)

    if err != nil {
        return false, err
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return false, err
    }

    if rows < 1 {
        return false, nil
    }

    return true, nil
}

func GetFollowers(spotifyID string, paginationKey string) (*responses.PaginationResponse[[]responses.User, string], bool, error) {
    query := `
            SELECT users.spotifyid, users.username, users.bio, users.userrole 
            FROM followers 
            INNER JOIN  users 
            ON users.spotifyid = followers.userfollowed 
            WHERE followers.userfollowed = $1 AND users.spotifyid < $2 ORDER BY users.spotifyid LIMIT 25 `

    rows, err := DB.Driver.Query(query, spotifyID, paginationKey)

    if err != nil {
        return nil, false, err
    }

    userResponses := []responses.User{}

    for rows.Next() {
        user := responses.User{}
        err := rows.Scan(&user.SpotifyID, &user.Username, &user.Bio, &user.Role)
        if err != nil {
            if err == sql.ErrNoRows {
                return nil, false, nil 
            } 
            return nil, false, err
        }
        userResponses = append(userResponses, user)
    }

    paginationResponse := &responses.PaginationResponse[[]responses.User, string]{}
    paginationResponse.DataResponse = userResponses

    if len(userResponses) > 0 {
        lastUser := userResponses[len(userResponses)-1]
        paginationResponse.PaginationKey = lastUser.SpotifyID
    } else {
        paginationResponse.PaginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
    }

    return paginationResponse, true, nil
}


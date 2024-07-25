package db

import (
	"context"
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
    bio := sql.NullString{}
    err := row.Scan(&bio, &userResponse.Role)
    userResponse.Bio = bio.String

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    return userResponse, nil
}

/* =================== READ ================== */

func GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, error) {
    query := "SELECT spotifyid, userrole, username, bio FROM users WHERE spotifyid = $1"
    row := DB.Driver.QueryRow(query, spotifyID)

    userResponse := &responses.User{}
    bio := sql.NullString{}
    err := row.Scan(&userResponse.SpotifyID, &userResponse.Username, &userResponse.Role, &bio)
    userResponse.Bio = bio.String

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    return userResponse, nil
}


/* PROPERTY UPDATES */
func UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, error) { 
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
    bio := sql.NullString{}
    err := res.Scan(&bio, &userResponse.Role, &userResponse.SpotifyID, &userResponse.Username)
    userResponse.Bio = bio.String

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    return userResponse, nil

}

func DeleteUserByID(spotifyID string) error {
    query := "DELETE FROM users WHERE spotifyID = $1"
    res, err := DB.Driver.Exec(query, spotifyID)

    if err != nil {
        return HandleDatabaseError(err)
    }

    num, err :=  res.RowsAffected()

    if err != nil {
        return HandleDatabaseError(err)
    }

    if num < 1 {
        return HandleDatabaseError(sql.ErrNoRows)
    }

    return nil
}


/* RELATIONAL UDPATES */
func UnfollowUser(spotifyID string, otherUserSpotifyID string) error {
    query := "DELETE FROM followers WHERE follower = $1 AND userfollowed = $2"

    res, err := DB.Driver.Exec(query, spotifyID, otherUserSpotifyID)

    if err != nil {
        return HandleDatabaseError(err)
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return HandleDatabaseError(err)
    }

    if rows < 1 {
        return HandleDatabaseError(sql.ErrNoRows)
    }

    return nil
}

func FollowUser(spotifyID string, otherUserSpotifyID string) error {
    query := "INSERT INTO followers (follower, userFollowed) VALUES ($1, $2)"

    res, err := DB.Driver.Exec(query, spotifyID, otherUserSpotifyID)

    if err != nil {
        return HandleDatabaseError(err)
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return HandleDatabaseError(err)
    }

    if rows < 1 {
        return HandleDatabaseError(sql.ErrNoRows)
    }

    return nil
}

func GetFollowers(spotifyID string, paginationKey string) (*responses.PaginationResponse[[]responses.User, string], error) {
    tx, err := DB.Driver.BeginTx(context.Background(), nil)

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    defer tx.Rollback()
    query := `SELECT spotifyid FROM users WHERE spotifyid = $1`

    row, err := tx.Exec(query, spotifyID)

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    count, err := row.RowsAffected()

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    if count < 1 {
        return nil, HandleDatabaseError(sql.ErrNoRows)
    }

    query = `
            SELECT users.spotifyid, users.username, users.bio, users.userrole 
            FROM followers 
            INNER JOIN  users 
            ON users.spotifyid = followers.userfollowed 
            WHERE followers.userfollowed = $1 AND users.spotifyid < $2 ORDER BY users.spotifyid LIMIT 25 `

    rows, err := tx.Query(query, spotifyID, paginationKey)

    if err != nil {
        return nil, HandleDatabaseError(err)
    }

    userResponses := []responses.User{}

    for rows.Next() {
        user := responses.User{}
        err := rows.Scan(&user.SpotifyID, &user.Username, &user.Bio, &user.Role)
        if err != nil {
            return nil, HandleDatabaseError(err)
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

    err = tx.Commit()

    if err != nil {
        return nil, HandleDatabaseError(err)
    } 

    return paginationResponse, nil
}


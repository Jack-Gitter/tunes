package daos

import (
	"context"
	"database/sql"

	"github.com/Jack-Gitter/tunes/db"
	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

type UsersDAO struct {
    DB *sql.DB
}

type IUsersDAO interface {
    UpsertUserOnLogin(username string, spotifyID string) (*responses.User, error)
    GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, error)
    UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, error)
    DeleteUserByID(spotifyID string) error
    UnfollowUser(spotifyID string, otherUserSpotifyID string) error
    FollowUser(spotifyID string, otherUserSpotifyID string) error 
    GetFollowers(spotifyID string, paginationKey string) (*responses.PaginationResponse[[]responses.User, string], error)
}

func(u *UsersDAO) UpsertUserOnLogin(username string, spotifyID string) (*responses.User, error) {
	query := "INSERT INTO users (spotifyid, username, userrole) values ($1, $2, 'BASIC') ON CONFLICT (spotifyID) DO UPDATE SET username=$2 RETURNING bio, userrole"
	row := u.DB.QueryRow(query, spotifyID, username)

	userResponse := &responses.User{}
	userResponse.Username = username
	userResponse.SpotifyID = spotifyID

	bio := sql.NullString{}
	err := row.Scan(&bio, &userResponse.Role)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	userResponse.Bio = bio.String

	return userResponse, nil
}

func(u *UsersDAO) GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, error) {
	query := "SELECT spotifyid, userrole, username, bio FROM users WHERE spotifyid = $1"
	row := u.DB.QueryRow(query, spotifyID)

	userResponse := &responses.User{}

	bio := sql.NullString{}
	err := row.Scan(&userResponse.SpotifyID, &userResponse.Role, &userResponse.Username, &bio)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	userResponse.Bio = bio.String

	return userResponse, nil
}

func(u *UsersDAO) UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, error) {

    updateUserMap := make(map[string]any)
    mapstructure.Decode(updatedUser, &updateUserMap)

    conditionals := make(map[string]any)
    conditionals["spotifyID"] = spotifyID

    returning := []string{"bio", "userrole", "spotifyid", "username"}

    query, values := db.PatchQueryBuilder("users", updateUserMap, conditionals, returning)

	res := u.DB.QueryRow(query, values...)

	userResponse := &responses.User{}
	bio := sql.NullString{}
	err := res.Scan(&bio, &userResponse.Role, &userResponse.SpotifyID, &userResponse.Username)
	userResponse.Bio = bio.String

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	return userResponse, nil

}

func(u *UsersDAO) DeleteUserByID(spotifyID string) error {
	query := "DELETE FROM users WHERE spotifyID = $1"
	res, err := u.DB.Exec(query, spotifyID)

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	num, err := res.RowsAffected()

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	if num < 1 {
		return customerrors.WrapBasicError(sql.ErrNoRows)
	}

	return nil
}

func(u *UsersDAO) UnfollowUser(spotifyID string, otherUserSpotifyID string) error {
	query := "DELETE FROM followers WHERE follower = $1 AND userfollowed = $2"

	res, err := u.DB.Exec(query, spotifyID, otherUserSpotifyID)

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	if rows < 1 {
		return customerrors.WrapBasicError(sql.ErrNoRows)
	}

	return nil
}

func(u *UsersDAO) FollowUser(spotifyID string, otherUserSpotifyID string) error {
	query := "INSERT INTO followers (follower, userFollowed) VALUES ($1, $2)"

	res, err := u.DB.Exec(query, spotifyID, otherUserSpotifyID)

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	if rows < 1 {
		return customerrors.WrapBasicError(sql.ErrNoRows)
	}

	return nil
}

func(u *UsersDAO) GetFollowers(spotifyID string, paginationKey string) (*responses.PaginationResponse[[]responses.User, string], error) {

    paginationResponse := &responses.PaginationResponse[[]responses.User, string]{}

    tx, err := u.DB.BeginTx(context.Background(), nil)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    exists, err := doesUserExist(tx, spotifyID)

    if err != nil {
        return nil, err
    }

    if !exists {
        return nil, customerrors.WrapBasicError(sql.ErrNoRows)
    }

    followers, err := getUserFollowersPaginated(tx, spotifyID, paginationKey)

    if err != nil {
        return nil, err
    }

    paginationResponse.DataResponse = followers
    paginationResponse.PaginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"

    if len(followers) > 0 {
        paginationResponse.PaginationKey = followers[len(followers)-1].SpotifyID
    } 

    err = tx.Commit()

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    return paginationResponse, nil
}


func doesUserExist(executor db.QueryExecutor, spotifyID string) (bool, error) {
    query := `SELECT COUNT(*) FROM users WHERE spotifyid = $1`

    row := executor.QueryRow(query, spotifyID)

    count := 0
    err := row.Scan(&count)

    if err != nil {
        return false, customerrors.WrapBasicError(err)
    }

    if count >= 1 {
        return true, nil
    }

    return false, nil

}

func getUserFollowersPaginated(executor db.QueryExecutor, spotifyID string, paginationKey string) ([]responses.User, error) {
    query := ` SELECT users.spotifyid, users.username, users.bio, users.userrole 
                FROM followers 
                INNER JOIN  users 
                ON users.spotifyid = followers.follower 
                WHERE followers.userfollowed = $1 AND users.spotifyid < $2 ORDER BY users.spotifyid LIMIT 25 `

        rows, err := executor.Query(query, spotifyID, paginationKey)

        if err != nil {
            return nil, customerrors.WrapBasicError(err)
        }

        userResponses := []responses.User{}
        bio := sql.NullString{}

        for rows.Next() {
            user := responses.User{}
            err := rows.Scan(&user.SpotifyID, &user.Username, &bio, &user.Role)
            if err != nil {
                return nil, customerrors.WrapBasicError(err)
            }
            user.Bio = bio.String
            userResponses = append(userResponses, user)
        }

        return userResponses, nil
}

package daos

import (
	"database/sql"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

type UsersDAO struct { }

type IUsersDAO interface {
    UpsertUser(executor db.QueryExecutor, username string, spotifyID string) (*responses.User, error)
    GetUser(executor db.QueryExecutor, spotifyID string) (*responses.User, error)
    UpdateUser(executor db.QueryExecutor, spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, error)
    DeleteUser(executor db.QueryExecutor, spotifyID string) error
    UnfollowUser(executor db.QueryExecutor, spotifyID string, otherUserSpotifyID string) error
    FollowUser(executor db.QueryExecutor, spotifyID string, otherUserSpotifyID string) error 
    GetUserFollowers(executor db.QueryExecutor, spotifyID string, paginationKey string) ([]responses.User, error)
    GetUserFollowing(executor db.QueryExecutor, spotifyID string, paginationKey string) ([]responses.User, error)
    GetAllUserFollowing(executor db.QueryExecutor, spotifyID string) ([]responses.User, error)
    UpsertUserProfilePicture(executor db.QueryExecutor, spotifyID string) (*responses.ProfileImage, error)
}

func(u *UsersDAO) UpsertUser(executor db.QueryExecutor, username string, spotifyID string) (*responses.User, error) {
	query := "INSERT INTO users (spotifyid, username, userrole) values ($1, $2, 'BASIC') ON CONFLICT (spotifyID) DO UPDATE SET username=$2 RETURNING bio, userrole"
	row := executor.QueryRow(query, spotifyID, username)

	userResponse := &responses.User{}
	userResponse.Username = username
	userResponse.SpotifyID = spotifyID

	bio := sql.NullString{}
	err := row.Scan(&bio, &userResponse.Role)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	userResponse.Bio = bio.String
    userResponse.Email = ""

	return userResponse, nil
}

func(u *UsersDAO) GetUser(executor db.QueryExecutor, spotifyID string) (*responses.User, error) {
	query := "SELECT spotifyid, userrole, username, bio, email FROM users WHERE spotifyid = $1"
	row := executor.QueryRow(query, spotifyID)

	userResponse := &responses.User{}

	bio := sql.NullString{}
    email := sql.NullString{}
	err := row.Scan(&userResponse.SpotifyID, &userResponse.Role, &userResponse.Username, &bio, &email)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	userResponse.Bio = bio.String
    userResponse.Email = email.String

	return userResponse, nil
}

func(u *UsersDAO) UpdateUser(executor db.QueryExecutor, spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, error) {

    updateUserMap := make(map[string]any)
    mapstructure.Decode(updatedUser, &updateUserMap)

    conditionals := make(map[string]any)
    conditionals["spotifyID"] = spotifyID

    returning := []string{"bio", "userrole", "spotifyid", "username", "email"}

    query, values := db.PatchQueryBuilder("users", updateUserMap, conditionals, returning)

	res := executor.QueryRow(query, values...)

	userResponse := &responses.User{}
	bio := sql.NullString{}
    email := sql.NullString{}
	err := res.Scan(&bio, &userResponse.Role, &userResponse.SpotifyID, &userResponse.Username, &email)
	userResponse.Bio = bio.String
    userResponse.Email = email.String

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	return userResponse, nil

}

func(u *UsersDAO) DeleteUser(executor db.QueryExecutor, spotifyID string) error {
	query := "DELETE FROM users WHERE spotifyID = $1"
	res, err := executor.Exec(query, spotifyID)

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

func(u *UsersDAO) UnfollowUser(executor db.QueryExecutor, spotifyID string, otherUserSpotifyID string) error {
	query := "DELETE FROM followers WHERE follower = $1 AND userfollowed = $2"

	res, err := executor.Exec(query, spotifyID, otherUserSpotifyID)

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

func(u *UsersDAO) FollowUser(executor db.QueryExecutor, spotifyID string, otherUserSpotifyID string) error {
	query := "INSERT INTO followers (follower, userFollowed) VALUES ($1, $2)"

	res, err := executor.Exec(query, spotifyID, otherUserSpotifyID)

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

func(u *UsersDAO) GetUserFollowers(executor db.QueryExecutor, spotifyID string, paginationKey string) ([]responses.User, error) {

    query := ` SELECT users.spotifyid, users.username, users.bio, users.userrole 
                FROM followers 
                INNER JOIN  users 
                ON users.spotifyid = followers.follower 
                WHERE followers.userfollowed = $1 AND followers.follower > $2 ORDER BY users.spotifyid LIMIT 25 `

    rows, err := executor.Query(query, spotifyID, paginationKey)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    followers := []responses.User{}
    bio := sql.NullString{}

    for rows.Next() {
        user := responses.User{}
        err := rows.Scan(&user.SpotifyID, &user.Username, &bio, &user.Role)
        if err != nil {
            return nil, customerrors.WrapBasicError(err)
        }
        user.Bio = bio.String
        followers = append(followers, user)
    }

    return followers, nil
}

func(u *UsersDAO) GetUserFollowing(executor db.QueryExecutor, spotifyID string, paginationKey string) ([]responses.User, error) {

    query := ` SELECT users.spotifyid, users.username, users.bio, users.userrole 
                FROM followers 
                INNER JOIN  users 
                ON users.spotifyid = followers.userfollowed 
                WHERE followers.follower = $1 AND users.spotifyid > $2 ORDER BY users.spotifyid LIMIT 25 `

    rows, err := executor.Query(query, spotifyID, paginationKey)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    following := []responses.User{}
    bio := sql.NullString{}

    for rows.Next() {
        user := responses.User{}
        err := rows.Scan(&user.SpotifyID, &user.Username, &bio, &user.Role)
        if err != nil {
            return nil, customerrors.WrapBasicError(err)
        }
        user.Bio = bio.String
        following = append(following, user)
    }

    return following, nil

}

func(u *UsersDAO) GetAllUserFollowing(executor db.QueryExecutor, spotifyID string) ([]responses.User, error) {

    query := ` SELECT users.spotifyid, users.username, users.bio, users.userrole 
                FROM followers 
                INNER JOIN  users 
                ON users.spotifyid = followers.userfollowed 
                WHERE followers.follower = $1`

    rows, err := executor.Query(query, spotifyID)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    following := []responses.User{}
    bio := sql.NullString{}

    for rows.Next() {
        user := responses.User{}
        err := rows.Scan(&user.SpotifyID, &user.Username, &bio, &user.Role)
        if err != nil {
            return nil, customerrors.WrapBasicError(err)
        }
        user.Bio = bio.String
        following = append(following, user)
    }

    return following, nil

}

func(u *UsersDAO) UpsertUserProfilePicture(executor db.QueryExecutor, spotifyID string) (*responses.ProfileImage, error) {
    return nil, nil
}

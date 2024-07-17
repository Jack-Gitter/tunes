package db

import (
	"errors"
	"os"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/* =================== CREATE ================== */

func InsertUserIntoDB(username string, spotifyID string, role string) (*responses.User, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MERGE (u:User {spotifyID: $spotifyID, username: $username, bio: $bio, role: $role}) return properties(u) as User",
        map[string]any{
            "spotifyID": spotifyID,
            "username": username,
            "role": role,
            "bio": "",
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    userResponse := &responses.User{}
    user, _ := resp.Records[0].Get("User")
    mapstructure.Decode(user, userResponse)
    return userResponse, err
}

/* =================== READ ================== */

func GetUserFromDbBySpotifyID(spotifyID string) (*responses.User, bool, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u:User {spotifyID: $spotifyID}) RETURN properties(u) as userProperties",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, false, err
    }

    if len(res.Records) < 1 {
        return nil, false, nil
    }

    userResponse, found := res.Records[0].Get("userProperties")

    if !found {
        return nil, true, errors.New("user within the database has no properties")
    }

    user := &responses.User{}
    mapstructure.Decode(userResponse, user)

    return user, true, nil
}


/* PROPERTY UPDATES */
func UpdateUserPropertiesBySpotifyID(updatedUser *responses.User) (*responses.User, bool, error) { return nil, false, nil}


/* RELATIONAL UDPATES */
func FollowUserBySpotifyID() {}
func UnFollowUserBySpotifyID(){}

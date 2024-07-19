package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/* =================== CREATE ================== */

func InsertUserIntoDBIfNeeded(username string, spotifyID string, role responses.Role) (*responses.User, error) {
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
func UpdateUserPropertiesBySpotifyID(spotifyID string, updatedUser *requests.UpdateUserRequestDTO) (*responses.User, bool, error) { 
    query := "MATCH (u:User {spotifyID: $spotifyID}) SET "
    if updatedUser.Bio != nil {
        query += "u.Bio = $Bio"
    }
    if updatedUser.Role != nil {
        query += ", u.Role = $Role"
    }
    query += " return properties(u) as userProperties"

    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    query,
        map[string]any{
            "spotifyID": spotifyID,
            "Bio": updatedUser.Bio,
            "Role": updatedUser.Role,
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

func DeleteUserByID(spotifyID string) (bool, error) {
    query := `  MATCH (u:User {spotifyID: $spotifyID}) return true UNION
                MATCH (u)-[:Posted]->(p) DETACH DELETE p return true UNION
                MATCH (u2)-[f:Follows]->(u) DETACH DELETE f RETURN true UNION
                MATCH (u:User {spotifyID: $spotifyID}) DELETE u return true`

    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, query,
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return false, err
    }

    if len(resp.Records) < 1 {
        return false, nil
    }

    return true, nil

}

func GetFollowers(spotifyID string, paginationKey string) (*responses.PaginationResponse[[]responses.User, string], bool, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u1:User {spotifyID: $spotifyID}) return properties(u1) as User UNION
     MATCH (u2)-[f:Follows]->(u1) WHERE normalize(u2.spotifyID) < normalize($paginationKey) return properties(u2) as User ORDER BY u2.spotifyID DESC LIMIT 25`,
        map[string]any{
            "spotifyID": spotifyID,
            "paginationKey": paginationKey,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, false, err
    }

    if len(resp.Records) < 1 {
        return nil, false, nil
    }

    fmt.Println(resp.Records[0].Get("User"))
    users := []responses.User{}

    for _, record := range resp.Records {
        userResponse, exists := record.Get("User")
        if !exists { return nil, false, errors.New("post has no properties in database") }
        user := &responses.User{}
        mapstructure.Decode(userResponse, user)
        if user.SpotifyID != spotifyID {
            users = append(users, (*user))
        }
    }
    paginationResponse := &responses.PaginationResponse[[]responses.User, string]{}
    paginationResponse.DataResponse = users
    paginationResponse.PaginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
    if len(users) > 0 {
        paginationResponse.PaginationKey = users[len(users)-1].SpotifyID
    }

    return paginationResponse, true, nil
}

func UnfollowUser(spotifyID string, otherUserSpotifyID string) (bool, error)  {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u1:User {spotifyID: $spotifyID}) MATCH (u2:User {spotifyID: $otherUserSpotifyID}) MATCH (u1)-[f:Follows]->(u2) DELETE f return properties(u1) as user1, properties(u2) as user2",
        map[string]any{
            "spotifyID": spotifyID,
            "otherUserSpotifyID": otherUserSpotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return false, err
    }

    if len(resp.Records) < 1 {
        return false, nil
    }

    _, foundU1 := resp.Records[0].Get("user1")
    _, foundU2 := resp.Records[0].Get("user2")

    if !foundU1 || !foundU2 {
        return false, nil
    }
    
    return true, nil
}

/* RELATIONAL UDPATES */
func FollowUser(spotifyID string, otherUserSpotifyID string) (bool, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u1:User {spotifyID: $spotifyID}) MATCH (u2:User {spotifyID: $otherUserSpotifyID}) MERGE (u1)-[:Follows]->(u2) return properties(u1) as user1, properties(u2) as user2",
        map[string]any{
            "spotifyID": spotifyID,
            "otherUserSpotifyID": otherUserSpotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return false, err
    }

    if len(resp.Records) < 1 {
        return false, nil
    }

    _, foundU1 := resp.Records[0].Get("user1")
    _, foundU2 := resp.Records[0].Get("user2")

    if !foundU1 || !foundU2 {
        return false, nil
    }
    
    return true, nil
}




func UnFollowUserBySpotifyID(){}

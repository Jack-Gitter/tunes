package db

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

func createPost(spotifyID string, songID string, songName string, albumName string, albumArtURI string, albumID string, rating int, text string) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) 
     MERGE (p:Post {songID: "ya", songName: "ya", albumName: "ya", albumArtURI: "ya", albumID: "ya", rating: "ya", text: "ya"})
     CREATE (u)-[:Posted]->(p)
     RETURN properties(p) `,
        map[string]any{ 
            "songID": songID,
            "songName": songName,
            "albumName": albumName,
            "albumArtURI": albumArtURI,
            "albumID": albumID,
            "rating": rating,
            "text": text,
            "spotifyID": spotifyID,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )
}



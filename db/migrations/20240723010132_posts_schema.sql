-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
	posterSpotifyID varchar(255) references users(spotifyid) ON DELETE CASCADE NOT NULL,
	songID varchar(255) NOT NULL,
    PRIMARY KEY (posterSpotifyID, songID),
	createdAt timestamp with time zone NOT NULL,
	updatedAt timestamp with time zone NOT NULL,
	songName varchar(255) NOT NULL,
	review varchar(255) NOT NULL,
	albumID varchar(255) NOT NULL,
	albumName varchar(255) NOT NULL,
    rating int NOT NULL,
	albumArtURI varchar(255)
	
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd

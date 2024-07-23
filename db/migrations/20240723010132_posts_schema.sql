-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
	albumArtURI varchar(255),
	albumID varchar(255),
	albumName varchar(255),
	createdAt timestamp with time zone,
	rating int,
	songID varchar(255),
	songName varchar(255),
	review varchar(255),
	updatedAt timestamp with time zone,
	posterSpotifyID varchar(255) references users(spotifyid),
    primary key (posterSpotifyID, songID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd

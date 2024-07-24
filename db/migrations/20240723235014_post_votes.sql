-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_votes (
    posterSpotifyID varchar(255),
    songID varchar(255),
    postsongID varchar(255),
	createdAt timestamp with time zone,
    updatedAt timestamp with time zone,
    liked boolean,
    FOREIGN KEY(posterSpotifyID, songID) references posts(posterspotifyid, songid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post_votes
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_votes (
    voterspotifyID varchar(255) references users(spotifyid),
    posterSpotifyID varchar(255),
    postsongID varchar(255),
	createdAt timestamp with time zone,
    updatedAt timestamp with time zone,
    liked boolean,
    FOREIGN KEY (posterSpotifyID, postsongID) references posts(posterspotifyid, songid),
    PRIMARY KEY (voterspotifyID, posterSpotifyID, postsongid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post_votes
-- +goose StatementEnd

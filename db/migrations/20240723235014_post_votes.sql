-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_votes (
    voterspotifyID varchar(255) references users(spotifyid) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    posterSpotifyID varchar(255) NOT NULL,
    postsongID varchar(255) NOT NULL,
    PRIMARY KEY (voterspotifyID, posterSpotifyID, postsongid),
    FOREIGN KEY (posterSpotifyID, postsongID) references posts(posterspotifyid, songid) ON DELETE CASCADE ON UPDATE CASCADE,
	createdAt timestamp with time zone NOT NULL,
    updatedAt timestamp with time zone NOT NULL,
    liked boolean NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post_votes
-- +goose StatementEnd

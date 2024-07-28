-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    commentorspotifyid varchar(255) references users(spotifyid) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    posterspotifyid varchar(255),
    songid varchar(255),
    commentText varchar(255),
    FOREIGN KEY (posterspotifyid, songid) references posts(posterspotifyid, songid) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments;
-- +goose StatementEnd

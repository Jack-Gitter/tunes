-- +goose Up
-- +goose StatementBegin
CREATE TABLE comment_votes (
    commentid int references comments(commentid) ON DELETE CASCADE ON UPDATE CASCADE, 
    liked boolean, 
    voterspotifyid varchar(255) references users(spotifyid) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY(commentid, voterspotifyid)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comment_votes
-- +goose StatementEnd

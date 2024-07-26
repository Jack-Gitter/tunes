-- +goose Up
-- +goose StatementBegin
CREATE TABLE FOLLOWERS (
	follower varchar(255) references users(spotifyid) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
	userFollowed varchar(255) references users(spotifyid) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
	PRIMARY KEY (follower, userFollowed)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE FOLLOWERS;
-- +goose StatementEnd

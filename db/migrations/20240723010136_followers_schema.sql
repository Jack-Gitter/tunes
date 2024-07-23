-- +goose Up
-- +goose StatementBegin
CREATE TABLE FOLLOWERS (
	follower varchar(255) references users(spotifyid),
	userFollowed varchar(255) references users(spotifyid),
	primary key(follower, userFollowed)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE FOLLOWERS;
-- +goose StatementEnd

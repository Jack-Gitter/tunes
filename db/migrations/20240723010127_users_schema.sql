-- +goose Up
-- +goose StatementBegin
CREATE TABLE USERS
(
    bio character varying(255),
    userrole character varying(255),
    spotifyid character varying(255) PRIMARY KEY,
    username character varying(255) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE USERS;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE USERS
(
    spotifyid character varying(255) PRIMARY KEY,
    userrole character varying(255) NOT NULL,
    username character varying(255),
    bio character varying(255),
    email character varying(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE USERS;
-- +goose StatementEnd

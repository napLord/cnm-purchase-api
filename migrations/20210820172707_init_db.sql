-- +goose Up
CREATE TABLE purchase (
  id BIGSERIAL PRIMARY KEY,
  foo BIGINT NOT NULL
);

-- +goose Down
DROP TABLE purchase;

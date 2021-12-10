-- migrate:up
CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL
);

-- migrate:down
DROP TABLE users;

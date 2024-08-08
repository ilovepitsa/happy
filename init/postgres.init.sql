DROP TABLE IF EXISTS user_follow;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id serial,
    last_name varchar(100) not null,
    first_name varchar(100) not null,
    email varchar(255) not null,
    birthday date not null
    PRIMARY KEY(id)
);

CREATE TABLE user_follow (
    subscriber_id integer references users(id),
    user_id integer references users(id),
    notify_before varchar(5)
);


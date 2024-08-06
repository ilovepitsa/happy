DROP TABLE IF EXISTS user_follow;
DROP TABLE IF EXISTS notification;
DROP TABLE IF EXISTS users;
DROP TYPE  IF EXISTS notification_type;
DROP TYPE  IF EXISTS notification_template;

CREATE TYPE notification_type as enum('email', 'web');

CREATE TABLE users (
    id serial,
    last_name varchar(100) not null,
    first_name varchar(100) not null,
    email varchar(255) not null,
    birthday date not null
    PRIMARY KEY(id)
);


CREATE TABLE follows (
    subscriber_id integer references users(id),
    target_id integer references users(id),
    notify_before varchar(5)
);

CREATE TABLE notification_template (
    id serial,
    template text
);

CREATE TABLE notification (
    id serial,
    target integer references users(id),
    type notification_type, 
    text text
);
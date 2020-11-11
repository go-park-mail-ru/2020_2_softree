create database if not exists users;
create user app_user with encrypted password 'NeverGonnaGiveYouUp';
grant all privileges on database users to app_user;

DROP TABLE IF EXISTS user_trade;

CREATE TABLE user_trade
(
    id       BIGSERIAL NOT NULL PRIMARY KEY,
    email    text,
    password text
);

CREATE TABLE watchlist
(
    user_id        bigint     NOT NULL,
    base_title     varchar(3) NOT NULL,
    currency_title varchar(3) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_trade (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS history_currency_by_minutes;
DROP TABLE IF EXISTS history_currency_by_hours;
DROP TABLE IF EXISTS history_currency_by_week;

CREATE TABLE history_currency_by_minutes
(
    title      text,
    value      decimal,
    updated_at timestamp
);

create index idx_history_currency_by_minutes on history_currency_by_minutes(updated_at);

CREATE TABLE history_currency_by_hours
(
    title      text,
    value      decimal,
    updated_at timestamp
);

CREATE TABLE history_currency_by_week
(
    title      text,
    value      decimal,
    updated_at timestamp
);

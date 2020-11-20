DROP TABLE IF EXISTS user_trade Cascade;
DROP TABLE IF EXISTS watchlist;
DROP TABLE IF EXISTS wallet;
DROP TABLE IF EXISTS payment_history;

CREATE TABLE user_trade
(
    id       BIGSERIAL NOT NULL PRIMARY KEY,
    email    text,
    password text,
    avatar   text default ''
);

CREATE TABLE watchlist
(
    user_id        bigint     NOT NULL,
    base_title     varchar(3) NOT NULL,
    currency_title varchar(3) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_trade (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE accounts
(
    user_id bigint NOT NULL,
    title   text,
    value   decimal,

    FOREIGN KEY (user_id) REFERENCES user_trade (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE payment_history
(
    user_id    bigint NOT NULL,
    from_title text,
    to_title   text,
    value      decimal,
    amount     decimal,
    updated_at timestamp,

    FOREIGN KEY (user_id) REFERENCES user_trade (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

truncate user_trade cascade;
truncate watchlist;
truncate accounts;
truncate payment_history;

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

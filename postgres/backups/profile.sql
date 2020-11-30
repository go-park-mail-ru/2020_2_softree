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

create index idx_watchlist on watchlist(user_id);

CREATE TABLE accounts
(
    user_id bigint NOT NULL,
    title   text,
    value   decimal,

    FOREIGN KEY (user_id) REFERENCES user_trade (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

create index idx_accounts on accounts(user_id, title);

truncate user_trade cascade;

create index idx_payment_history on payment_history(user_id);

CREATE TABLE payment_history
(
    user_id    bigint NOT NULL,
    base       text,
    curr       text,
    value      decimal,
    amount     decimal,
    sell       text,
    updated_at timestamp,

    FOREIGN KEY (user_id) REFERENCES user_trade (id)
        ON DELETE CASCADE ON UPDATE CASCADE
);
truncate watchlist;
truncate accounts;
truncate payment_history;
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

insert into user_trade (email, password) values ('test@test.com', '$2a$10$C5lRv07hMZS6VUzPq27Xu.l7j9zFLLl02hyVW9KyuZ5FDRfFRcf16');

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

insert into accounts(user_id, title, "value") values (1, 'USD', 1000.0);

create index idx_accounts on accounts(user_id, title);

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

insert into payment_history(user_id, base, curr, "value", amount, sell, updated_at)
values (1, 'USD', 'RUB', 76.7, 500, true, '2020-11-29');
insert into payment_history(user_id, base, curr, "value", amount, sell, updated_at)
values (1, 'USD', 'RUB', 73.6, 200, true, '2020-10-23');

CREATE TABLE wallet_history
(
    user_id    bigint NOT NULL,
    value      decimal,
    updated_at timestamp with time zone
);

insert into wallet_history(user_id, "value", updated_at)
values (1, 950, '2020-12-19');
insert into wallet_history(user_id, "value", updated_at)
values (1, 700, '2020-09-29');
insert into wallet_history(user_id, "value", updated_at)
values (1, 900, '2020-11-29');

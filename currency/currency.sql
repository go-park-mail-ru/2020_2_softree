DROP TABLE IF EXISTS history_currency_by_minutes;
DROP TABLE IF EXISTS history_currency_by_hours;
DROP TABLE IF EXISTS history_currency_by_day;

CREATE TABLE history_currency_by_minutes
(
    title      text,
    value      decimal,
    updated_at timestamp
);

create index idx_history_currency_by_minutes on history_currency_by_minutes(title);
cluster history_currency_by_minutes using idx_history_currency_by_minutes;

CREATE TABLE history_currency_by_hours
(
    title      text,
    value      decimal,
    updated_at timestamp
);

create index idx_history_currency_by_hours on history_currency_by_hours(updated_at, title);

CREATE TABLE history_currency_by_day
(
    title      text,
    value      decimal,
    updated_at timestamp
);

create index idx_history_currency_by_day on history_currency_by_day(updated_at, title);

truncate history_currency_by_minutes;
truncate history_currency_by_hours;
truncate history_currency_by_day;

create user app_rates with encrypted password 'NeverGonnaGiveYouUp';
grant all privileges on database rates to app_rates;

DROP TABLE IF EXISTS history_currency_by_minutes;
DROP TABLE IF EXISTS history_currency_by_hours;
DROP TABLE IF EXISTS history_currency_by_week;

CREATE TABLE history_currency_by_minutes(
	title text,
	value decimal,
	updated_at timestamp
);

CREATE TABLE history_currency_by_hours(
	title text,
	value decimal,
	updated_at timestamp
);

CREATE TABLE history_currency_by_week(
	title text,
	value decimal,
	updated_at timestamp
);
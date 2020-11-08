create database if not exists users;
create user app_user with encrypted password 'NeverGonnaGiveYouUp';
grant all privileges on database users to app_user;

DROP TABLE IF EXISTS user_trade;

CREATE TABLE user_trade(
	id BIGSERIAL NOT NULL PRIMARY KEY,
	email text,
	password text
);

CREATE TABLE watchlist(
	user_id bigint NOT NULL,
	base_title varchar(3) NOT NULL,
	currency_title varchar(3) NOT NULL,

	FOREIGN KEY (user_id) REFERENCES user_trade (id)
	ON DELETE CASCADE ON UPDATE CASCADE
);

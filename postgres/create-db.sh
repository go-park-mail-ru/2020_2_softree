#!/bin/bash

set -e
set -u

function create_user_and_database() {
	echo "Creating user and database"

	local database=$(echo $1 | tr ',' ' ' | awk '{print $1}')
	local owner=$(echo $1 | tr ',' ' ' | awk '{print $2}')
	local password=$(echo $1 | tr ',' ' ' | awk '{print $3}')

	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
		CREATE USER $owner WITH PASSWORD '$password';
		CREATE DATABASE $database;
		GRANT ALL PRIVILEGES ON DATABASE $database TO $owner;
	EOSQL
}

function restore_from_backup() {
	echo "Restore database from backup"

	local database=$(echo $1 | tr ',' ' ' | awk  '{print $1}')
	local owner=$(echo $1 | tr ',' ' ' | awk  '{print $2}')
	local backup=$(echo /backups/$database.sql)

	psql -U $owner -d $database -f $backup
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ':' ' '); do
		create_user_and_database $db
		restore_from_backup $db
	done
	echo "Multiple databases created from backups"
fi

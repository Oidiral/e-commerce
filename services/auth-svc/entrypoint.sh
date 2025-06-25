#!/bin/sh
set -e

timeout 30s sh -c "until pg_isready -h '$AUTH_DB_HOST' -p '$AUTH_DB_PORT' -U '$AUTH_DB_USER'; do sleep 1; done"

DSN="user=$AUTH_DB_USER password=$AUTH_DB_PASSWORD host=$AUTH_DB_HOST port=$AUTH_DB_PORT dbname=$AUTH_DB_NAME sslmode=disable"

echo "Applying migrations from: $(ls -d /db/migrations)"
/usr/local/bin/goose -dir /db/migrations postgres "$DSN" up

exec /usr/local/bin/auth "$@"
#!/bin/sh
set -e

timeout 30s sh -c "until pg_isready -h '$CART_DB_HOST' -p '$CART_DB_PORT' -U '$CART_DB_USER'; do sleep 1; done"

DSN="user=$CART_DB_USER password=$CART_DB_PASSWORD host=$CART_DB_HOST port=$CART_DB_PORT dbname=$CART_DB_NAME sslmode=disable"

echo "Applying migrations from: $(ls -d db/migrations)"
/usr/local/bin/goose -dir db/migrations postgres "$DSN" up

exec ./cart "$@"
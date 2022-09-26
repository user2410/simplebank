#!/bin/sh

# This script is expected to be executed inside container

# Exit immediately if a simple command exits with a non-zero status
set -e

echo "run db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "run the app"
# Example: entrypoint.sh server start -> server start
exec "$@"
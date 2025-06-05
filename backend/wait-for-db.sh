#!/bin/bash
# wait-for-db.sh

# Erasmo Cardoso da Silva
# Desenvolvedor Full Stack

set -e

host="$1"
shift
cmd="$@"

# Add a short initial delay
sleep 2

>&2 echo "Waiting for Postgres server at $host..."

# Wait for the main postgres database to be available
until PGPASSWORD="$DB_PASSWORD" psql -h "$host" -U "$DB_USER" -d "postgres" -c '\q'; do
  >&2 echo "Postgres server is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres server is up. Waiting for database '$DB_NAME'..."

# Now wait for the specific database to be available
until PGPASSWORD="$DB_PASSWORD" psql -h "$host" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "Database '$DB_NAME' is unavailable - sleeping"
  sleep 1
done

>&2 echo "Database '$DB_NAME' is up - executing command: $cmd"
exec $cmd 
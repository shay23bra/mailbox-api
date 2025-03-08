#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
count=0
max_tries=60
until pg_isready -h postgres -U ${DB_USER:-postgres} > /dev/null 2>&1 || [ $count -eq $max_tries ]; do
  echo "Waiting for PostgreSQL ($count/$max_tries)..."
  count=$((count+1))
  sleep 5
done

if [ $count -eq $max_tries ]; then
  echo "PostgreSQL did not become ready in time!"
  exit 1
fi

echo "PostgreSQL is ready! Starting application..."
exec "$@"
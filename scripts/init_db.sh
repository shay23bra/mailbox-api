#!/bin/bash
set -e

# Load environment variables from .env file - safer approach
if [ -f ../.env ]; then
    # Read .env line by line and export valid variables
    while IFS='=' read -r key value || [ -n "$key" ]; do
        # Skip comments and empty lines
        if [[ $key != \#* ]] && [[ ! -z "$key" ]]; then
            # Remove any leading/trailing whitespace
            key=$(echo $key | xargs)
            value=$(echo $value | xargs)
            # Export non-empty values
            if [[ ! -z "$key" ]]; then
                export "$key=$value"
            fi
        fi
    done < ../.env
fi

# Set default values if not in .env
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-mailbox}

echo "Creating database $DB_NAME if it doesn't exist..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME"

echo "Running migrations..."
for migration in ../migrations/*.sql; do
    echo "Applying migration: $migration"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration
done

echo "Database initialization completed."
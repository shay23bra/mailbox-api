#!/bin/bash
set -e

# Load environment variables from .env file
if [ -f ../.env ]; then
    export $(grep -v '^#' ../.env | xargs)
fi

# Set default values if not in .env
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-mailbox}

echo "Seeding departments..."
cat ../data/departments.csv | tail -n +2 | while IFS=, read -r id name; do
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "INSERT INTO departments (department_id, department_name) VALUES ($id, '$name') ON CONFLICT (department_id) DO NOTHING;"
done

echo "Seeding mailboxes..."
cat ../data/mailboxes.csv | tail -n +2 | while IFS=, read -r identifier full_name job_title department_id manager_identifier; do
    if [ "$manager_identifier" = "null" ]; then
        manager_identifier="NULL"
    else
        manager_identifier="'$manager_identifier'"
    fi
    
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "INSERT INTO mailboxes (mailbox_identifier, user_full_name, job_title, department_id, manager_mailbox_identifier) VALUES ('$identifier', '$full_name', '$job_title', $department_id, $manager_identifier) ON CONFLICT (mailbox_identifier) DO NOTHING;"
done

echo "Calculating organization metrics..."
cd ../ && go run main.go calculate-metrics

echo "Database seeding completed."
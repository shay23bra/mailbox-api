#!/bin/bash
set -e

# Get the actual container name
POSTGRES_CONTAINER=$(docker-compose ps -q postgres)
if [ -z "$POSTGRES_CONTAINER" ]; then
  echo "Error: PostgreSQL container not found. Make sure it's running with 'docker-compose up -d postgres'"
  exit 1
fi
echo "Using PostgreSQL container: $POSTGRES_CONTAINER"

# Set database variables
DB_USER=postgres
DB_NAME=mailbox

echo "Creating database schema..."
# Apply migrations directly through Docker Compose
docker-compose exec -T postgres psql -U ${DB_USER} -d ${DB_NAME} -c "
CREATE TABLE IF NOT EXISTS departments (
    department_id INT PRIMARY KEY,
    department_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS mailboxes (
    mailbox_identifier VARCHAR(100) PRIMARY KEY,
    user_full_name VARCHAR(100) NOT NULL,
    job_title VARCHAR(100) NOT NULL,
    department_id INT NOT NULL,
    manager_mailbox_identifier VARCHAR(100),
    org_depth INT NOT NULL DEFAULT 0,
    sub_org_size INT NOT NULL DEFAULT 0,
    FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

CREATE INDEX IF NOT EXISTS idx_mailboxes_department_id ON mailboxes(department_id);
CREATE INDEX IF NOT EXISTS idx_mailboxes_manager_id ON mailboxes(manager_mailbox_identifier);
CREATE INDEX IF NOT EXISTS idx_mailboxes_org_depth ON mailboxes(org_depth);
CREATE INDEX IF NOT EXISTS idx_mailboxes_sub_org_size ON mailboxes(sub_org_size);
"

echo "Seeding departments..."
# Copy departments.csv to container
docker cp ../data/departments.csv ${POSTGRES_CONTAINER}:/tmp/departments.csv
# Seed departments
docker-compose exec -T postgres bash -c "cat /tmp/departments.csv | tail -n +2 | while IFS=, read -r id name; do psql -U ${DB_USER} -d ${DB_NAME} -c \"INSERT INTO departments (department_id, department_name) VALUES (\$id, '\$name') ON CONFLICT (department_id) DO NOTHING;\"; done"

echo "Seeding mailboxes..."
# Copy mailboxes.csv to container
docker cp ../data/mailboxes.csv ${POSTGRES_CONTAINER}:/tmp/mailboxes.csv
# Seed mailboxes
docker-compose exec -T postgres bash -c "cat /tmp/mailboxes.csv | tail -n +2 | while IFS=, read -r identifier full_name job_title department_id manager_identifier; do if [ \"\$manager_identifier\" = \"null\" ]; then manager_identifier=\"NULL\"; else manager_identifier=\"'\$manager_identifier'\"; fi; psql -U ${DB_USER} -d ${DB_NAME} -c \"INSERT INTO mailboxes (mailbox_identifier, user_full_name, job_title, department_id, manager_mailbox_identifier) VALUES ('\$identifier', '\$full_name', '\$job_title', \$department_id, \$manager_identifier) ON CONFLICT (mailbox_identifier) DO NOTHING;\"; done"

echo "Database setup completed successfully!"
echo ""
echo "NEXT STEPS:"
echo "1. Run the application to calculate organization metrics:"
echo "   go run main.go"
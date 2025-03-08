# Mailbox API

A REST API for managing mailboxes in an organization and providing organizational hierarchy metrics.

## Features

- List mailboxes with support for searching, filtering, sorting, and pagination
- Query mailboxes by organizational hierarchy metrics (depth, sub-organization size)
- Role-based access control (CEO and CTO roles)
- Automatic filtering based on user role (CEO sees all, CTO sees only their sub-organization)
- Scalable architecture designed for large organizations

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Docker and Docker Compose (optional)

## Project Structure

```
mailbox-api/
├── main.go
├── .env
├── .env.test
├── Makefile
├── api/
│   ├── handler/
│   ├── middleware/
│   └── router/
├── config/
├── db/
├── dto/
├── logger/
├── model/
├── repository/
├── service/
├── util/
├── migrations/
├── scripts/
├── test/
├── data/
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Configuration

The application uses a `.env` file for configuration. You can customize the following parameters:

```
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mailbox
DB_SSLMODE=disable

# Authentication
JWT_SECRET=secure-jwt-secret-key-should-be-long-and-complex
TOKEN_EXPIRY=60 # minutes

# Logging
LOG_LEVEL=info
USE_SYSLOG=false
```

Make sure to set a strong, unique JWT_SECRET for production environments.

## Getting Started

### Running with Docker Compose

The easiest way to run the application is using Docker Compose:

```bash
docker-compose up -d
```

This will start both the PostgreSQL database and the API server.

### Running Locally

1. Start PostgreSQL:

```bash
docker-compose up -d postgres
```

2. Initialize and seed the database:

**Option A: Using psql client (requires PostgreSQL client tools)**
```bash
cd scripts
chmod +x init_db.sh
./init_db.sh
chmod +x seed_db.sh
./seed_db.sh
```

**Option B: Using Docker (no PostgreSQL client tools required)**
```bash
cd scripts
chmod +x docker_setup_db.sh
./docker_setup_db.sh
```
This script will:
- Start the PostgreSQL container if not running
- Create the database schema
- Import all test data from CSV files
- Calculate organization metrics
All without requiring PostgreSQL client tools on your local machine.

3. Run the application:

```bash
go run main.go
```

Alternatively, use the Makefile:

```bash
# Start everything with docker-compose, initialize and seed the database, and run the application
make dev
```

## Testing

In the file `mailbox_api_tests.md` you can find examples for local tests using curl.

The project includes comprehensive test suites:

- Unit tests for service layer
- Integration tests for repositories
- API tests for endpoints

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Test Environment

Tests use a separate test database (`mailbox_test`) which is created and populated automatically during test execution. Configuration for the test environment is in `.env.test`.

The test suite:

1. Creates a test database
2. Applies the schema
3. Seeds it with test data
4. Runs the tests
5. Cleans up after completion

This ensures tests are isolated and reproducible.

## API Endpoints

### Authentication

- `GET /api/token/ceo` - Get CEO token
- `GET /api/token/cto` - Get CTO token

### Mailboxes (Role-based access)

- `GET /api/mailboxes` - List mailboxes (CEO sees all, CTO sees only their sub-organization)
- `GET /api/mailboxes/:id` - Get a specific mailbox (CEO can access any, CTO can only access those in their sub-organization)
- `POST /api/mailboxes/calculate-metrics` - Recalculate organization metrics (CEO only)

### Query Parameters

- `search`: Search by name/title/department (partial match)
- `department`: Filter by department ID
- `org_depth_exact`: Filter by exact org depth
- `org_depth_gt`: Filter by org depth greater than
- `org_depth_lt`: Filter by org depth less than
- `sub_org_size_min`: Filter by minimum sub-org size
- `sub_org_size_max`: Filter by maximum sub-org size
- `sort_by`: Sort by field (can specify multiple)
- `sort_dir`: Sort direction (asc/desc, can specify multiple)
- `fields`: Select specific fields (comma-separated)
- `page`: Page number (default: 1)
- `page_size`: Page size (default: 10)

## Authentication

The API uses JWT for authentication. To access protected endpoints, include the JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

You can obtain a token from the `/api/token/ceo` or `/api/token/cto` endpoints.

## Role-Based Access Control

The API implements automatic data filtering based on user roles:

- **CEO**: Has access to all mailboxes across the organization
- **CTO**: Can only view mailboxes within their sub-organization (direct and indirect reports)

This approach eliminates the need for separate endpoints and ensures that users only see data they are authorized to access.

## Example Requests

### List mailboxes with filters

```
GET /api/mailboxes?department=2&org_depth_lt=2&sort_by=user_full_name&sort_dir=desc&page=1&page_size=10
```

### Search mailboxes

```
GET /api/mailboxes?search=software
```

### Get mailboxes with specific fields

```
GET /api/mailboxes?fields=mailbox_identifier,user_full_name,job_title,org_depth
```

### Multiple sorting criteria

```
GET /api/mailboxes?sort_by=department_id&sort_by=org_depth&sort_dir=asc&sort_dir=desc
```

## Continuous Integration

The project is set up for CI/CD with the `make ci` command which runs all tests and builds the application.

# Future Enhancements and Improvements

The current implementation provides a solid foundation, but there are numerous opportunities for enhancement. Below are several areas where the application could be expanded and improved:

## Authentication and Security

- **Firebase Authentication Integration**: Replace the current JWT implementation with Firebase Auth for more robust authentication, including:
  - Social login options (Google, Microsoft, etc.)
  - Multi-factor authentication
  - Email verification
  - Password reset functionality
  - User management dashboard

- **OAuth 2.0 Support**: Add support for OAuth 2.0 to allow integration with enterprise identity providers:
  - Microsoft Entra ID (formerly Azure AD) for Microsoft 365 organizations
  - Google Workspace authentication
  - SAML integration for enterprise SSO solutions

- **Enhanced RBAC**: Extend the role-based access control system:
  - Support for custom roles with fine-grained permissions
  - Dynamic role assignment based on organizational structure
  - Activity auditing and access logs

## Performance & Scalability

- **Caching Layer**: Implement Redis to cache frequently accessed organization data and pre-computed hierarchies, reducing database load
- **Database Optimization**: Add read replicas and advanced indexing strategies to support large organizations with thousands of employees
- **Asynchronous Processing**: Move resource-intensive operations like hierarchy calculations to background workers using message queues
- **Horizontal Scaling**: Enhance the application to support containerization for Kubernetes deployment with auto-scaling capabilities
- **Multi-Organization Support**: Transform the system into a multi-tenant solution:
  - Organization isolation and data segregation
  - Custom branding per organization
  - Organization-specific configurations

## User Experience

- **Web Dashboard**: Develop a React-based frontend featuring:
  - Interactive organization chart visualization
  - Role-specific views and dashboards
  - Responsive design for mobile and desktop access
- **Real-time Updates**: Implement WebSockets to push organizational changes to connected clients instantly
- **Enhanced Authentication**: Integrate with Firebase Auth for improved authentication including social logins, MFA, and SSO options
- **Search Improvements**: Add full-text search and advanced filtering for quicker access to organization data

## CI/CD & DevOps

- **Automated Pipeline**: Expand the CI/CD pipeline with:
  - Automated performance testing to catch regressions
  - Security scanning for vulnerabilities
  - Smoke tests in staging environments
  - Blue/green deployment strategy
- **Monitoring Stack**: Implement comprehensive monitoring with Prometheus and Grafana to track system health and performance metrics
- **Infrastructure as Code**: Define all infrastructure components using Terraform or AWS CloudFormation for reproducible environments
- **Feature Flags**: Implement feature flag system to safely roll out new functionality to subsets of users

## License

This project is licensed under the MIT License.
# Mailbox API Test Results

This document contains test results for the Mailbox API, demonstrating role-based access control and various query capabilities.

## Authentication

First, we obtain authentication tokens for both CEO and CTO roles:

```bash
# Get CEO token
CEO_TOKEN=$(curl -s "http://localhost:8080/api/token/ceo" | jq -r '.token')

# Get CTO token
CTO_TOKEN=$(curl -s "http://localhost:8080/api/token/cto" | jq -r '.token')
```

```bash
# Store tokens for later use
echo "CEO Token: $CEO_TOKEN"
> CEO Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiY2VvIiwiZXhwIjoxNzQxNDc1NjQ4LCJuYmYiOjE3NDE0NzIwNDgsImlhdCI6MTc0MTQ3MjA0OH0.-lxS3MqYQJtIvlywhT_agkXtcWIFKWup4LmtrBsIU3U
echo "CTO Token: $CTO_TOKEN"
> CTO Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiY3RvIiwiZXhwIjoxNzQxNDc1NjUwLCJuYmYiOjE3NDE0NzIwNTAsImlhdCI6MTc0MTQ3MjA1MH0.t1AlTDRHSQHZx5Ty1sy4WMpoxFgXdtoBkGBJvd_4w_8
```


## Listing Mailboxes

### CEO Access (All Mailboxes)

The CEO has access to all mailboxes in the organization:

```bash
curl -H "Authorization: Bearer $CEO_TOKEN" "http://localhost:8080/api/mailboxes"
```

**Response:**
```json
{
  "data": [
    {
      "mailbox_identifier": "alice.johnson@falafel.org",
      "user_full_name": "Alice Johnson",
      "job_title": "Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "bob.smith@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "bob.smith@falafel.org",
      "user_full_name": "Bob Smith",
      "job_title": "Project Manager",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "david.brown@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "carol.green@falafel.org",
      "user_full_name": "Carol Green",
      "job_title": "Marketing Analyst",
      "department_id": 3,
      "department": "Marketing",
      "manager_mailbox_identifier": "emma.davis@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "charlie.lee@falafel.org",
      "user_full_name": "Charlie Lee",
      "job_title": "HR Coordinator",
      "department_id": 4,
      "department": "Human Resources",
      "manager_mailbox_identifier": "isabella.white@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "david.brown@falafel.org",
      "user_full_name": "David Brown",
      "job_title": "CTO",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "isabella.white@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "emma.davis@falafel.org",
      "user_full_name": "Emma Davis",
      "job_title": "Marketing Lead",
      "department_id": 3,
      "department": "Marketing",
      "manager_mailbox_identifier": "isabella.white@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "frank.wilson@falafel.org",
      "user_full_name": "Frank Wilson",
      "job_title": "Sales Executive",
      "department_id": 5,
      "department": "Sales",
      "manager_mailbox_identifier": "isabella.white@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "grace.miller@falafel.org",
      "user_full_name": "Grace Miller",
      "job_title": "Customer Support Lead",
      "department_id": 3,
      "department": "Marketing",
      "manager_mailbox_identifier": "emma.davis@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "henry.moore@falafel.org",
      "user_full_name": "Henry Moore",
      "job_title": "Finance Manager",
      "department_id": 6,
      "department": "Finance",
      "manager_mailbox_identifier": "isabella.white@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "isabella.white@falafel.org",
      "user_full_name": "Isabella White",
      "job_title": "CEO",
      "department_id": 1,
      "department": "Executive",
      "manager_mailbox_identifier": "",
      "org_depth": 0,
      "sub_org_size": 0
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 16,
    "total_pages": 2
  }
}
```

### CTO Access (Sub-organization Only)

The CTO only has access to mailboxes within their sub-organization:

```bash
curl -H "Authorization: Bearer $CTO_TOKEN" "http://localhost:8080/api/mailboxes"
```

**Response:**
```json
{
  "data": [
    {
      "mailbox_identifier": "steve.johnson@falafel.org",
      "user_full_name": "Steve Johnson",
      "job_title": "Product Manager",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "david.brown@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "bob.smith@falafel.org",
      "user_full_name": "Bob Smith",
      "job_title": "Project Manager",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "david.brown@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "lisa.evans@falafel.org",
      "user_full_name": "Lisa Evans",
      "job_title": "Senior Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "bob.smith@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "alice.johnson@falafel.org",
      "user_full_name": "Alice Johnson",
      "job_title": "Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "bob.smith@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "paul.martin@falafel.org",
      "user_full_name": "Paul Martin",
      "job_title": "Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "alice.johnson@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 5,
    "total_pages": 1
  }
}
```

## Filtering Mailboxes

The API supports various filters including department, organization depth, and pagination:

```bash
curl -H "Authorization: Bearer $CEO_TOKEN" "http://localhost:8080/api/mailboxes?department=2&org_depth_lt=2&page=1&page_size=5"
```

**Response:**
```json
{
  "data": [
    {
      "mailbox_identifier": "alice.johnson@falafel.org",
      "user_full_name": "Alice Johnson",
      "job_title": "Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "bob.smith@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "bob.smith@falafel.org",
      "user_full_name": "Bob Smith",
      "job_title": "Project Manager",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "david.brown@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "david.brown@falafel.org",
      "user_full_name": "David Brown",
      "job_title": "CTO",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "isabella.white@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "lisa.evans@falafel.org",
      "user_full_name": "Lisa Evans",
      "job_title": "Senior Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "bob.smith@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    },
    {
      "mailbox_identifier": "paul.martin@falafel.org",
      "user_full_name": "Paul Martin",
      "job_title": "Software Engineer",
      "department_id": 2,
      "department": "Technology",
      "manager_mailbox_identifier": "alice.johnson@falafel.org",
      "org_depth": 0,
      "sub_org_size": 0
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 5,
    "total_items": 6,
    "total_pages": 2
  }
}
```

## Field Selection

The API allows selecting specific fields to return:

```bash
curl -H "Authorization: Bearer $CEO_TOKEN" "http://localhost:8080/api/mailboxes?fields=mailbox_identifier,user_full_name,job_title"
```

**Response:**
```json
{
  "data": [
    {
      "job_title": "Software Engineer",
      "mailbox_identifier": "alice.johnson@falafel.org",
      "user_full_name": "Alice Johnson"
    },
    {
      "job_title": "Project Manager",
      "mailbox_identifier": "bob.smith@falafel.org",
      "user_full_name": "Bob Smith"
    },
    {
      "job_title": "Marketing Analyst",
      "mailbox_identifier": "carol.green@falafel.org",
      "user_full_name": "Carol Green"
    },
    {
      "job_title": "HR Coordinator",
      "mailbox_identifier": "charlie.lee@falafel.org",
      "user_full_name": "Charlie Lee"
    },
    {
      "job_title": "CTO",
      "mailbox_identifier": "david.brown@falafel.org",
      "user_full_name": "David Brown"
    },
    {
      "job_title": "Marketing Lead",
      "mailbox_identifier": "emma.davis@falafel.org",
      "user_full_name": "Emma Davis"
    },
    {
      "job_title": "Sales Executive",
      "mailbox_identifier": "frank.wilson@falafel.org",
      "user_full_name": "Frank Wilson"
    },
    {
      "job_title": "Customer Support Lead",
      "mailbox_identifier": "grace.miller@falafel.org",
      "user_full_name": "Grace Miller"
    },
    {
      "job_title": "Finance Manager",
      "mailbox_identifier": "henry.moore@falafel.org",
      "user_full_name": "Henry Moore"
    },
    {
      "job_title": "CEO",
      "mailbox_identifier": "isabella.white@falafel.org",
      "user_full_name": "Isabella White"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 16,
    "total_pages": 2
  }
}
```

## Getting Specific Mailboxes

### CEO Access to CTO Mailbox

```bash
curl -H "Authorization: Bearer $CEO_TOKEN" "http://localhost:8080/api/mailboxes/david.brown@falafel.org"
```

**Response:**
```json
{
  "mailbox_identifier": "david.brown@falafel.org",
  "user_full_name": "David Brown",
  "job_title": "CTO",
  "department_id": 2,
  "department": "Technology",
  "manager_mailbox_identifier": "isabella.white@falafel.org",
  "org_depth": 0,
  "sub_org_size": 0
}
```

### CTO Access to Own Mailbox

```bash
curl -H "Authorization: Bearer $CTO_TOKEN" "http://localhost:8080/api/mailboxes/david.brown@falafel.org"
```

**Response:**
```json
{
  "mailbox_identifier": "david.brown@falafel.org",
  "user_full_name": "David Brown",
  "job_title": "CTO",
  "department_id": 2,
  "department": "Technology",
  "manager_mailbox_identifier": "isabella.white@falafel.org",
  "org_depth": 0,
  "sub_org_size": 0
}
```

### CTO Access to Subordinate's Mailbox

```bash
curl -H "Authorization: Bearer $CTO_TOKEN" "http://localhost:8080/api/mailboxes/bob.smith@falafel.org"
```

**Response:**
```json
{
  "mailbox_identifier": "bob.smith@falafel.org",
  "user_full_name": "Bob Smith",
  "job_title": "Project Manager",
  "department_id": 2,
  "department": "Technology",
  "manager_mailbox_identifier": "david.brown@falafel.org",
  "org_depth": 0,
  "sub_org_size": 0
}
```

### CTO Attempting to Access Outside Their Sub-organization

CTO attempts to access the Marketing Lead's mailbox, which is outside their sub-organization:

```bash
curl -H "Authorization: Bearer $CTO_TOKEN" "http://localhost:8080/api/mailboxes/emma.davis@falafel.org"
```

**Response:**
```json
{
  "error": "Access denied"
}
```

## Organization Metrics Calculation

### CEO Calculating Metrics (Allowed)

```bash
curl -X POST -H "Authorization: Bearer $CEO_TOKEN" "http://localhost:8080/api/mailboxes/calculate-metrics"
```

**Response:**
```json
{
  "message": "Org metrics calculated successfully"
}
```

### CTO Attempting to Calculate Metrics (Denied)

```bash
curl -X POST -H "Authorization: Bearer $CTO_TOKEN" "http://localhost:8080/api/mailboxes/calculate-metrics"
```

**Response:**
```json
{
  "error": "Access denied"
}
```

## Observations

1. **Role-Based Access Control** - CEO can access all mailboxes, while CTO can only access their own sub-organization
2. **Filtering and Pagination** - API supports filtering by department, org depth, and includes pagination
3. **Field Selection** - API allows selecting specific fields for more focused responses
4. **Authorization** - Properly denies access to resources outside of a user's permissions
5. **Metrics Calculation** - Only CEO role is allowed to recalculate organization metrics

All of these features demonstrate a well-designed API with proper security controls and flexible query capabilities.

## Note on Sub-organization Size and Org Depth

It appears that the org_depth and sub_org_size values are currently set to 0. These might need to be recalculated to reflect the actual organizational hierarchy.
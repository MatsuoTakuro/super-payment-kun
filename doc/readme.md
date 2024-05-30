# Super Payment Kun

This is a simple API for managing invoices and payments.

## Setup and Test API

### Prerequisites

- Docker installed

### Initialization

1. Clone the repository:

    ```sh
    git clone github.com/MatsuoTakuro/super-payment-kun
    cd super-payment-kun
    ```

2. Initialize and start Docker containers:

    ```sh
    make init
    ```

This will build the images for the Go app and MySQL, migrate/seed the database, and start the containers.

## API Usage

### Test Login

- **Endpoint:** `POST localhost:8080/api/testlogin`
- **Description:** Logs in as a test user and returns a JWT token.
- **Curl Command:**

    ```sh
    curl -X POST http://localhost:8080/api/testlogin
    ```

Example response:

```json
{
    "api_code": "000000",
    "data": {
        "会社A_佐藤太郎_admin_access_token": "eyJhbGciOi..."
    }
}
```

Use the token from the response for authorization in subsequent requests.

### Create an Invoice

- **Endpoint:** `POST localhost:8080/api/invoices`
- **Description:** Creates a new invoice.
- **Curl Command:**

    ```sh
    curl -X POST http://localhost:8080/api/invoices \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer {{token}}" \
        -d '{
            "vendor_id": "323e4567-e89b-12d3-a456-426614174000",
            "vendor_bank_account_id": "423e4567-e89b-12d3-a456-426614174000",
            "payment_amount": 120000,
            "due_date": "2024-07-01"
        }'
    ```

- In the above example, `vendor_id` is for ベンダーA and `vendor_bank_account_id` is for ベンダーA口座.
- Replace `{{token}}` with the token obtained from the test login response.

### Get a List of Invoices

- **Endpoint:** `GET localhost:8080/api/invoices`
- **Description:** Retrieves a list of invoices.
- **Query Parameters:**
  - `from_due_date` (required): Starting due date (e.g., `2024-05-28`)
  - `to_due_date` (required): Ending due date (e.g., `2024-07-01`)
  - `limit` (optional): Maximum number of results (default: 30)
  - `cursor` (optional): Cursor for pagination
  - `direction` (optional): Pagination direction (`fwd` for forward, `bwd` for backward)

#### Retrieve initial list

- **Curl Command:**

    ```sh
    curl -X GET "http://localhost:8080/api/invoices?from_due_date=2024-05-28&to_due_date=2024-07-01" \
        -H "Authorization: Bearer {{token}}"
    ```

- Note that at the initial stage, 30 invoices have been seeded in the database and so, you will see the first 30 invoices if the limit is not specified.

#### Retrieve next set using cursor

- **Curl Command:**

    ```sh
    curl -X GET "http://localhost:8080/api/invoices?from_due_date=2024-05-28&to_due_date=2024-07-01&cursor={{cursor}}&direction=fwd" \
        -H "Authorization: Bearer {{token}}"
    ```

- Replace `{{cursor}}` with the cursor value obtained from the previous response.
- Note that you will see the invoice you created in the previous step if the due date is within the specified range.

### OpenAPI Specification

You can preview the OpenAPI specification using the following command:

```sh
make watch_api_spec
```

You have to install the [Redocly CLI](https://github.com/Redocly/redocly-cli#redocly-cli) if not already installed and start a local server to preview the API documentation.

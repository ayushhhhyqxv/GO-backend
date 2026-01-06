**A collection of Go Snippets having various backend implementations using Echo and PostgreSQL**


###  Rate Limiter
Implements API rate limiting :
- User signup with unique API key generation
- Request rate limiting middleware (15 requests per 15 seconds)
- Database-backed rate limit tracking
- Retry-After header support

**Endpoints:**
- `POST /signup` - Create user and generate API key
- `GET /data` - Protected endpoint with rate limiting

###  Pagination
Product management API with advanced filtering and pagination:
- Paginated product listing with configurable page size
- Sorting by any field (ascending/descending)
- Text-based filtering on product name and category
- Sample data seeding

**Endpoints:**
- `POST /seed` - Populate database with sample products
- `GET /products` - Fetch products with pagination, filtering, and sorting

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 5)
- `sortField` - Field to sort by
- `sortOrder` - Sort order (asc/desc)
- `filter` - Search term for name and category

###  Binary Upload
File upload and retrieval system:
- Store files as binary data in PostgreSQL
- Download previously uploaded files
- Support for any file type

**Endpoints:**
- `POST /upload` - Upload file
- `GET /file/:id` - Download file by ID

###  Email Setup
Email sending utility with attachment support using SMTP configuration from environment variables.

## Stack
- **Framework:** Echo (v4)
- **Database:** PostgreSQL with GORM ORM
- **Configuration:** .env file for environment variables
- **Email:** gomail package for SMTP

## Instructions

1. Create `.env` file with required variables:
    ```
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=<your_password>
    DB_NAME=<db_name>
    DB_PORT=5432
    SERVER_PORT=<any_port_which_isn't_busy>
    ```

2. Install dependencies:
    ```
    go mod download
    ```

3. Run project:
    ```
    go run main.go
    ```

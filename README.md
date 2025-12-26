
A RESTful API built with the **GoLang**, the **Echo Web Framework**, and **PostgreSQL**. This demonstrates GET,POST,PUT,PATCH capabilities with proper HTTP status codes and SQL injection protection.



* **RESTful Routing**: Clean and organized endpoints using the Echo framework.
* **Database Integration**: Persistent storage with PostgreSQL using the `pq` driver.
* **Data Binding**: Automatic JSON-to-Struct mapping using `c.Bind`.
* **Robust Error Handling**: Standardized HTTP status codes (200, 201, 400, 404, 500).
* **Middleware**: Integrated request logging for debugging and performance monitoring.

---

## Requirements: 

* **Go**: 1.18 or higher.
* **PostgreSQL**: Installed and running on `localhost:5432`.
* **Postman**: Recommended for testing the API endpoints.

---

## Setup

1. **Clone the repository** to your local machine.
2. **Initialize the Go module**:
   ```bash
   go mod init restAPI
   go get [github.com/labstack/echo/v4](https://github.com/labstack/echo/v4)
   go get [github.com/lib/pq](https://github.com/lib/pq) ```
3. **Set your own DSN**:
    ```bash
    dsn := "host=localhost port=5432 user=postgres password=YOUR_PW dbname=newdb sslmode=disable"
    ```
4. **Run the Program**:
    ```bash
    go run main.go
    ```
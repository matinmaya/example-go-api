

# 🧪 example-go-api

A modular, clean-architecture-style Golang API using [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io). This project demonstrates a scalable approach to building RESTful APIs with features like:

- Request binding and validation
- Query filtering using struct tags
- Pagination and sorting
- Standardized API responses
- Role management example module

---

## 🚀 Features

- 🧩 Modular project structure
- 🔍 Flexible query filtering (e.g. `like`, `equal`, `between`)
- 📄 Pagination with `page`, `limit`, `sort_by`, `sort_dir`
- ✅ Centralized validation and binding utilities
- 🔄 Clean mapping between DTOs and models
- 📦 Designed for real-world extensions (auth, logging, etc.)

---

## 📁 Project Structure

```
example-go-api/
├── cmd/                    # App entrypoint (e.g. main.go)
├── internal/
│   ├── helpers/            # Context and DB helpers
│   └── modules/
│       └── user/
│           ├── rolehandler/    # HTTP handlers
│           ├── roleservice/    # Business logic
│           ├── rolerepository/ # DB access layer
│           ├── rolemodel/      # GORM model and query DTO
│           └── roledto/        # (optional) Response/Request DTOs
├── pkg/
│   ├── binding/            # Request validation and param helpers
│   ├── filterscopes/       # Query param filters and GORM scopes
│   ├── paginator/          # Pagination logic and helpers
│   ├── response/           # Standard API response wrapper
│   ├── validators/         # Reusable validation logic
│   └── reqctx/       # Field reflection utilities
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Requirements

- Go 1.20 or later
- SQLite (for testing) or configure any GORM-compatible database

---

## 🛠️ Setup & Run

```bash
git clone https://github.com/matinmaya/example-go-api.git
cd example-go-api
go mod tidy
go run main.go
```

If using another database (e.g. PostgreSQL), update the DB config in your helper/init file.

---

## 📥 Sample API: Role Management

### Create Role
```http
POST /roles
Content-Type: application/json

{
  "name": "admin",
  "description": "Administrator role"
}
```

### List Roles with Filters
```http
GET /roles?name=admin&status=1&sort_by=created_at&sort_dir=desc&page=1&limit=10
```

### Sample Response
```json
{
  "message": "success",
  "data": {
    "total": 1,
    "total_page": 1,
    "rows": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ]
  },
  "errors": null
}
```

---

## 🧠 Query Filtering Syntax

```go
type RoleListQuery struct {
  Name   string `form:"name" filter:"like,column=name"`
  Status int    `form:"status" filter:"equal"`
}
```

Supports:
- `equal`
- `like`
- `in`
- `between`, `not_between`
- `is_null`, `is_not_null`

---

## 📄 Response Format

All responses follow this standard format:

```json
{
  "message": "success",
  "data": { ... },
  "errors": null
}
```

---

## ✅ TODO

- [ ] JWT authentication & middleware
- [ ] Swagger/OpenAPI docs
- [ ] Automated tests
- [ ] CI/CD pipeline

---

## 📄 License

This project is licensed under the MIT License. See `LICENSE` for details.
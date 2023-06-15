# GoTodo REST API

## Golang 101 Mini Project

Simple ToDos list REST API using Go

### Features

- Add Todo Item
- Update Todo Item
- Delete Todo Item
- Batch Delete Todo Item
- Get Todo Item By ID
- Get All Todo Items

## Authentication

- JWT token based authentication and authorization. The token stored in `Cookie`.
- User email stored in request context. So, the handler function can access the current user email through the request context.

## How To Run?

1. Clone the repository

```bash
git clone https://github.com/hakimamarullah/gotodo-api.git
```

2. Run using Go >= v1.20

```bash
cd gotodo-api
```

```bash
go run main.go
```

Note: You can also run this app using Docker by building the image using the Dockerfile which has been provided in the reposisory.

## Next Development

- [] Add due date field to TodoItem
- [] Use environment variable for database credential, application's port, etc.

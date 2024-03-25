# Feature Toggles

## Stack

### Frontend

- Angular v17
- Angular Material v17

### Backend

- Golang v1.22
- Gin v1.9
- GORM v1.25
- SQlite

## How to run

### Docker (Recommended)

Be sure that ports 8081 and 80 are available. Else you can change the ports in the `compose.yml` file.

```yaml
    # Backend
    ports:
      - "{PORT_FOR_BACKEND}:8081" # Default: 8081
    # NGINX
    ports:
      - "{PORT_FOR_FRONTEND}:80" # Default: 80
```

Then run:

```bash
docker-compose up -d
```

Then you can access the application at `http://localhost:{PORT_FOR_FRONTEND}`

### Local

#### Start Backend

```bash
cd backend
export DEMO="true" # Populate the database with fake data
go run cmd/main.go
```

#### Start Frontend

```bash
cd frontend
npm install
npm start
```

Then you can access the application at `http://localhost:4200`

### Additional REST Endpoint

The endpoint that determines if feature is active for user it at
`POST /api/v1/customers/:customerID`

The tests in `backend/pkg/model/toggle_test.go`, corresponds to the Example API response in the homework document.

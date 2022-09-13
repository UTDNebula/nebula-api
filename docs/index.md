```
NEBULA-API
├───api             - API
│   ├───.env        - API developer environment variables (MONGODB_URI)
│   ├───server.go   - Server (go run server.go)
│   ├───configs     - MondoDB configuration
│   ├───controllers - Route controllers
│   ├───models      - Schema models
│   ├───responses   - Route responses
│   ├───routes      - Define routes
│   └───ts          - Old TypeScript Implementation (to be removed)
│       ├───controllers
│       ├───models
│       └───routes
├───docs
│   └───schemas     - Database schemas
├───node_modules    - Old Node Modules (to be removed)
└───scraper         - Scraper
    ├───configs
    ├───data
    └───scripts
```

To connect to your Mongo database, in /api/.env: set MONGODB_URI accordingly:
```
MONGODB_URI=<insert_connection_string_here>
```

To start the api server in your local environment from /nebula-api/api/:
```
go run server.go
```

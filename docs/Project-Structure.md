# Project Structure

This file provides a brief overview of important directories and files.

```text
nebula-api/
├── .github/              # Configuration for the GitHub repo. Used for GitHub workflows.
├── docs/                 # Pages for developer wiki documentation.
├── rest/                 # Directory for our RESTful API endpoints.
│   ├── configs/          # Congfiguration to set up for endpoints.
│   ├── controllers/      # Main part of API. Handles HTTP requests and responds to them.
│   ├── docs/             # Auto-generated docs on API routes from Swagger.
│   ├── responses/        # Was used for cloud storage endpoints. No longer used.
│   ├── routes/           # Defines all possible API routes.
│   ├── schema/           # Defines structure for the data. Shared with api-tools.
│   ├── Dockerfile        # Configuration for building a container image for Nebula API.
│   └── server.go         # Initializes API endpoints
├── .env.template         # Template for environment variables used in this project.
├── .gitignore            # For specifying files we do not want GitHub to track.
├── build.bat             # For building Nebula API on Windows.
├── go.mod                # For Go dependencies.
├── go.sum                # Also for Go dependencies.
├── Makefile              # For building Nebula API on Mac and Linux.
├── README.me             # Documentation mainly for users of Nebula API.
├── update-gateway.bat    # For updating the API gateway configuration on Windows.
└── update-gateway.sh     # For updating the API gateway configuration on Mac or Linux.
```

## Next Step

See [How-to-Contribute.md](How-to-Contribute.md)

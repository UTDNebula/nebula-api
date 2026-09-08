# API Architecture
Nebula API provides UTD data for others to use, including information on courses, professors, events, and more.

## API Overview
Here's a quick look at Nebula API.

### REST API
Currently, the Nebula API is a RESTful API. [What does that mean?](https://aws.amazon.com/what-is/restful-api/)

The Nebula API offers information about the following:
- Courses
- Professors
- Sections
- Grades
- Comet Calendar Events
- Mazevo Events
- Astra Events
- Club Event - provided by UTD Clubs
- Rooms
- Discounts

> [!NOTE]
> Nebula API powers only UTD Trends. UTD Clubs and Notebook have their own API internal to their own projects. Specifically, the Clubs endpoints are powered by the internal Clubs database.

# API Diagram
``` mermaid
flowchart LR
    A((Nebula API)) --> B((REST))
    C[(Club DB)] --> A
    D[(Mongo DB)] --> A
    F@{shape: bucket, label: "S3"}--> A
```

## API Documentation
We currently use Swagger for our API documentation. Swagger automatically generates documentation for our API routes.

If you make any changes, make sure to run either the Makefile or build.bat to update the documentation.

## Languages and Core Technologies
### Golang
[Tour of Go](https://go.dev/tour/welcome/1) - Start here if you're new to Golang.
[Documentation](https://go.dev/ref/spec) - For more technical details.

The main language we use for Nebula API and the API tools.

### MongoDB
[Go Mongo Driver Documentation](https://www.mongodb.com/docs/drivers/go/current/) - For using MongoDB with Golang.

Nebula API stores the majority of its data in a MongoDB database. We interact with the Mongo database with the Mongo Go driver.

## Libraries

### Gin
[Documentation](https://gin-gonic.com/en/docs/)

Gin is a HTTP web framework for Golang that allows to build or RESTful API.

### Swagger
[Our Swagger API Docs](https://api.utdnebula.com/swagger/index.html)

Swagger is responsible for auto-generating our documentation for our API endpoints.

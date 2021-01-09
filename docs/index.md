# Admin Console

Heroku endpoint: `https://comet-data-service.herokuapp.com/`
## Courses

1. `GET /`:  UI for CRUD console
2. `POST /courses`: add a new course
3. `GET /courses/name/<name>`: get course info by name
4. `GET /courses/id/<id>`: get course info by id
5. `GET /courses`: get all courses
6. `PUT /courses/<id>`: edit course with id
7. `DELETE /course/<id>`: delete course with id

### Course format:

```javascript
{
    "id": 0, 
    "course": "ACCT 2301", 
    "name": "Introductory Financial Accounting", 
    "hours": "3", 
    "description": "An introduction to financial reporting...", 
    "inclass": "3", 
    "outclass": "0", 
    "period": "S", 
    "prerequisites": [] // work-in-progress
}
```
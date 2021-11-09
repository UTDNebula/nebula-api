# API Reference

## Authentication
Every request to the Nebula Data API requires a valid API key included in the HTTP header `Authorization`.
```HTTP
GET /v1/example/ HTTP/1.1
Authorization: apikey
```

## Sections

`/v1/sections/id`
| Value | Description | Example |
| ------------- | ------------- | ------------- |
| id  | Section Name  | acct2301.001.21f |
```js
Returns
{
    "term": "string",
    "title": "string",
    "course_number": "string",
    "school": "string",
    "location": "string",
    "activity_type": "string",
    "class_number": "string",
    "days": "string",
    "assistants": "string",
    "times": "string",
    "topic": "string",
    "core_area": "string",
    "department": "string",
    "section_name": "string",
    "course_prefix": "string",
    "instructors": "string",
    "section_number": "string"
}
```
---

`/v1/sections/search?property=value`
| Property | Description | Example |
| ------------- | ------------- | ------------- |
| activity_type  | Whether the section is a lecture or a lab  | Lecture |
| assistants | Teaching assistants in the section  | Brad%20Nathan |
| class_number | Class number | 80642 |
| core_area | Core area of the class | ∅ |
| course_prefix | Abbreviation of the course | acct |
| days | On what days the section occurs | Monday%2C%20Wednesday%2C%20Friday |
| department | What department is the course under | mgmt |
| instructors | The instructors teaching the section | Cameron%20Holstead |
| location | Where the section meets | SOM%201.110 |
| section_name | Section name | acct2301.001.21f |
| section_number | Section number | 001 |
| term | Term in which the section occurs | 21f
| times | When the section occurs | 08:00%2008:50 |
| title | Course title | Introductory%20Financial%20Accounting |
| topic | Course topic | |
```js
Returns
[
   {
      "term": "string",
      "title": "string",
      "course_number": "string",
      "school": "string",
      "location": "string",
      "activity_type": "string",
      "class_number": "string",
      "days": "string",
      "assistants": "string",
      "times": "string",
      "topic": "string",
      "core_area": "string",
      "department": "string",
      "section_name": "string",
      "course_prefix": "string",
      "instructors": "string",
      "section_number": "string"
   },
]
```




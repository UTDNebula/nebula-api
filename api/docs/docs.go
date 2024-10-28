// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/course": {
            "get": {
                "description": "\"Returns all courses matching the query's string-typed key-value pairs\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "courseSearch",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The course's official number",
                        "name": "course_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The course's subject prefix",
                        "name": "subject_prefix",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The course's title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The course's description",
                        "name": "description",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The course's school",
                        "name": "school",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The number of credit hours awarded by successful completion of the course",
                        "name": "credit_hours",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The level of education that this course course corresponds to",
                        "name": "class_level",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The type of class this course corresponds to",
                        "name": "activity_type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The grading status of this course",
                        "name": "grading",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The internal (university) number used to reference this course",
                        "name": "internal_course_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The weekly contact hours in lecture for a course",
                        "name": "lecture_contact_hours",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The frequency of offering a course",
                        "name": "offering_frequency",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A list of courses",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/schema.Course"
                            }
                        }
                    }
                }
            }
        },
        "/course/{id}": {
            "get": {
                "description": "\"Returns the course with given ID\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "courseById",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the course to get",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A course",
                        "schema": {
                            "$ref": "#/definitions/schema.Course"
                        }
                    }
                }
            }
        },
        "/grades/overall": {
            "get": {
                "description": "\"Returns the overall grade distribution\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "gradeAggregationOverall",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The course's subject prefix",
                        "name": "prefix",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The course's official number",
                        "name": "number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professor's first name",
                        "name": "first_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professors's last name",
                        "name": "last_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The number of the section",
                        "name": "section_number",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A grade distribution array",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    }
                }
            }
        },
        "/grades/semester": {
            "get": {
                "description": "\"Returns grade distributions aggregated by semester\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "gradeAggregationBySemester",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The course's subject prefix",
                        "name": "prefix",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The course's official number",
                        "name": "number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professor's first name",
                        "name": "first_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professors's last name",
                        "name": "last_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The number of the section",
                        "name": "section_number",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "An array of grade distributions for each semester included",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/responses.GradeResponse"
                            }
                        }
                    }
                }
            }
        },
        "/professor": {
            "get": {
                "description": "\"Returns all professors matching the query's string-typed key-value pairs\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "professorSearch",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The professor's first name",
                        "name": "first_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professor's last name",
                        "name": "last_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "One of the professor's title",
                        "name": "titles",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professor's email address",
                        "name": "email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The professor's phone number",
                        "name": "phone_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The building of the location of the professor's office",
                        "name": "office.building",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The room of the location of the professor's office",
                        "name": "office.room",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A hyperlink to the UTD room locator of the professor's office",
                        "name": "office.map_uri",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A hyperlink pointing to the professor's official university profile",
                        "name": "profile_uri",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A link to the image used for the professor on the professor's official university profile",
                        "name": "image_uri",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The start date of one of the office hours meetings of the professor",
                        "name": "office_hours.start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The end date of one of the office hours meetings of the professor",
                        "name": "office_hours.end_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "One of the days that one of the office hours meetings of the professor",
                        "name": "office_hours.meeting_days",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The time one of the office hours meetings of the professor starts",
                        "name": "office_hours.start_time",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The time one of the office hours meetings of the professor ends",
                        "name": "office_hours.end_time",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The modality of one of the office hours meetings of the professor",
                        "name": "office_hours.modality",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The building of one of the office hours meetings of the professor",
                        "name": "office_hours.location.building",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The room of one of the office hours meetings of the professor",
                        "name": "office_hours.location.room",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A hyperlink to the UTD room locator of one of the office hours meetings of the professor",
                        "name": "office_hours.location.map_uri",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The _id of one of the sections the professor teaches",
                        "name": "sections",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A list of professors",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/schema.Professor"
                            }
                        }
                    }
                }
            }
        },
        "/professor/{id}": {
            "get": {
                "description": "\"Returns the professor with given ID\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "professorById",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the professor to get",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A professor",
                        "schema": {
                            "$ref": "#/definitions/schema.Professor"
                        }
                    }
                }
            }
        },
        "/section": {
            "get": {
                "description": "\"Returns all courses matching the query's string-typed key-value pairs\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "sectionSearch",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The section's official number",
                        "name": "section_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "An id that points to the course in MongoDB that this section is an instantiation of",
                        "name": "course_reference",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The name of the academic session of the section",
                        "name": "academic_session.name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The date of classes starting for the section",
                        "name": "academic_session.start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The date of classes ending for the section",
                        "name": "academic_session.end_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "One of the professors teaching the section",
                        "name": "professors",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The first name of one of the teaching assistants of the section",
                        "name": "teaching_assistants.first_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The last name of one of the teaching assistants of the section",
                        "name": "teaching_assistants.last_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The role of one of the teaching assistants of the section",
                        "name": "teaching_assistants.role",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The email of one of the teaching assistants of the section",
                        "name": "teaching_assistants.email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The internal (university) number used to reference this section",
                        "name": "internal_class_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The instruction modality for this section",
                        "name": "instruction_mode",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The start date of one of the section's meetings",
                        "name": "meetings.start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The end date of one of the section's meetings",
                        "name": "meetings.end_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "One of the days that one of the section's meetings",
                        "name": "meetings.meeting_days",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The time one of the section's meetings starts",
                        "name": "meetings.start_time",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The time one of the section's meetings ends",
                        "name": "meetings.end_time",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The modality of one of the section's meetings",
                        "name": "meetings.modality",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The building of one of the section's meetings",
                        "name": "meetings.location.building",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The room of one of the section's meetings",
                        "name": "meetings.location.room",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A hyperlink to the UTD room locator of one of the section's meetings",
                        "name": "meetings.location.map_uri",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "One of core requirement codes this section fulfills",
                        "name": "core_flags",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A link to the syllabus on the web",
                        "name": "syllabus_uri",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A list of sections",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/schema.Section"
                            }
                        }
                    }
                }
            }
        },
        "/section/{id}": {
            "get": {
                "description": "\"Returns the section with given ID\"",
                "produces": [
                    "application/json"
                ],
                "operationId": "sectionById",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the section to get",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A section",
                        "schema": {
                            "$ref": "#/definitions/schema.Section"
                        }
                    }
                }
            }
        },
        "/swagger/index.html": {
            "get": {
                "description": "Returns the OpenAPI/swagger spec for the API",
                "produces": [
                    "text/html"
                ],
                "operationId": "swagger",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.GradeResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "schema.AcademicSession": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                }
            }
        },
        "schema.Assistant": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "schema.CollectionRequirement": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "options": {
                    "type": "array",
                    "items": {}
                },
                "required": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "schema.Course": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "activity_type": {
                    "type": "string"
                },
                "attributes": {},
                "catalog_year": {
                    "type": "string"
                },
                "class_level": {
                    "type": "string"
                },
                "co_or_pre_requisites": {
                    "$ref": "#/definitions/schema.CollectionRequirement"
                },
                "corequisites": {
                    "$ref": "#/definitions/schema.CollectionRequirement"
                },
                "course_number": {
                    "type": "string"
                },
                "credit_hours": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "enrollment_reqs": {
                    "type": "string"
                },
                "grading": {
                    "type": "string"
                },
                "internal_course_number": {
                    "type": "string"
                },
                "laboratory_contact_hours": {
                    "type": "string"
                },
                "lecture_contact_hours": {
                    "type": "string"
                },
                "offering_frequency": {
                    "type": "string"
                },
                "prerequisites": {
                    "$ref": "#/definitions/schema.CollectionRequirement"
                },
                "school": {
                    "type": "string"
                },
                "sections": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "subject_prefix": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "schema.Location": {
            "type": "object",
            "properties": {
                "building": {
                    "type": "string"
                },
                "map_uri": {
                    "type": "string"
                },
                "room": {
                    "type": "string"
                }
            }
        },
        "schema.Meeting": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "end_time": {
                    "type": "string"
                },
                "location": {
                    "$ref": "#/definitions/schema.Location"
                },
                "meeting_days": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "modality": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                }
            }
        },
        "schema.Professor": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "image_uri": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "office": {
                    "$ref": "#/definitions/schema.Location"
                },
                "office_hours": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.Meeting"
                    }
                },
                "phone_number": {
                    "type": "string"
                },
                "profile_uri": {
                    "type": "string"
                },
                "sections": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "titles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "schema.Section": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "academic_session": {
                    "$ref": "#/definitions/schema.AcademicSession"
                },
                "attributes": {},
                "core_flags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "course_reference": {
                    "type": "string"
                },
                "grade_distribution": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "instruction_mode": {
                    "type": "string"
                },
                "internal_class_number": {
                    "type": "string"
                },
                "meetings": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.Meeting"
                    }
                },
                "professors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "section_corequisites": {
                    "$ref": "#/definitions/schema.CollectionRequirement"
                },
                "section_number": {
                    "type": "string"
                },
                "syllabus_uri": {
                    "type": "string"
                },
                "teaching_assistants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.Assistant"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "api_key": {
            "type": "apiKey",
            "name": "x-api-key",
            "in": "header"
        }
    },
    "security": [
        {
            "api_key": []
        }
    ],
    "x-google-backend": {
        "address": "REDACTED"
    },
    "x-google-endpoints": [
        {
            "allowCors": true,
            "name": "nebula-api-2lntm5dxoflqn.apigateway.nebula-api-368223.cloud.goog"
        }
    ],
    "x-google-management": {
        "metrics": [
            {
                "displayName": "Read Requests CUSTOM",
                "metricKind": "DELTA",
                "name": "read-requests",
                "valueType": "INT64"
            }
        ],
        "quota": {
            "limits": [
                {
                    "metric": "read-requests",
                    "name": "read-limit",
                    "unit": "1/min/{project}",
                    "values": {
                        "STANDARD": 1000
                    }
                }
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "nebula-api-2lntm5dxoflqn.apigateway.nebula-api-368223.cloud.goog",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "nebula-api",
	Description:      "The public Nebula Labs API for access to pertinent UT Dallas data",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

# Course

## Overview

The `Course` object represents a course available at the University of Texas at Dallas. A `Course` should not be confused with a `Section` which is the actual instantiation of a `Course` with a professor and dedicated meeting times.

## Object Representation

```ts
Course = {
    "_id": ObjectId,
    "course_number": string,
    "subject_prefix": string,
    "title": string,
    "description": string,
    "school": string,
    "credit_hours": string,
    "class_level": string,
    "activity_type": string,
    "grading": string,
    "internal_course_number": string,
    "prerequisites": CollectionRequirement,
    "corequisites": CollectionRequirement,
    "lecture_contact_hours": string,
    "laboratory_contact_hours": string,
    "offering_frequency": string,
    "attributes": Object,
}
```

## Properties

> `._id`
>
> **Type**: ObjectId
>
> The MongoDB database id for the `Course` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.course_number`
>
> **Type**: string
>
> The course's official number.
>
> **Example**: 2417

> `.subject_prefix`
>
> **Type**: string
>
> The course's subject prefix.
>
> **Example**: MATH

> `.title`
>
> **Type**: string
>
> The course's title.
>
> **Example**: Calculus I

> `.description`
>
> **Type**: string
>
> The course's description.
>
> **Example**: "Functions, limits, continuity, differentiation; integration of..."

> `.school`
>
> **Type**: string
>
> 
>
> **Example**: 

> `.credit_hours`
>
> **Type**: string
>
> The number of credit hours awarded by successful completion of the course. Will be "V" if variable credit hours.
>
> **Examples**: "4", "V"

> `.class_level`
>
> **Type**: string
>
> The level of education that this course course corresponds to
>
> **Example**: 4

> `.activity_type`
>
> **Type**: string
>
> The type of class this course corresponds to.
>
> **Example**: Lecture

> `.grading`
>
> **Type**: string
>
> The grading status of this course.
>
> **Example**: Graded

> `.internal_course_number`
>
> **Type**: string
>
> The internal (university) number used to reference this course.
>
> **Example**: 008613

> `.prerequisites`
>
> **Type**: CollectionRequirement
>
> A collection of prerequisites that must be met before a student may enroll in a section of this course.

> `.corequisites`
>
> **Type**: CollectionRequirement
>
> A collection of all course requirements that must be met either before or while a student enrolls in a section of this course.

> `.lecture_contact_hours`
>
> **Type**: string
>
> The weekly contact hours in lecture for a course. This information is outlined in the UTD Course Policies page.
>
> **Example**: 2

> `.laboratory_contact_hours`
>
> **Type**: string
>
> The weekly contact hours in laboratory for a course. This information is outlined in the UTD Course Policies page.
>
> **Example**: 4

> `.offering_frequency`
>
> **Type**: string
>
> The frequency of offering a course. The meanings of each letter can be found in the UTD Course Policies page.
>
> **Example**: "S", "Y", "T", "R"

> `.attributes`
>
> **Type**: Object
>
> Space for additional data describing the course not listed elsewhere.

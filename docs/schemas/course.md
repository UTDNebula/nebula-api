# Course

## Overview

The `Course` object represents a course available at the University of Texas at Dallas. A `Course` should not be confused with a `Section` which is the actual instantiation of a `Course` with a professor and dedicated meeting times.

## Object Representation

```ts
Course = {
    "_id": ObjectId,
    "subject_prefix": string,
    "course_number": string,
    "title": string,
    "description": string,
    "enrollment_reqs": string,
    "school": string,
    "credit_hours": string,
    "class_level": string,
    "activity_type": string,
    "grading": string,
    "internal_course_number": string,
    "prerequisites": CollectionRequirement,
    "corequisites": CollectionRequirement,
    "co_or_pre_requisites": CollectionRequirement,
    "sections": Array<ObjectId>,
    "lecture_contact_hours": string,
    "laboratory_contact_hours": string,
    "offering_frequency": string,
    "catalog_year": string,
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

> `.subject_prefix`
>
> **Type**: string
>
> The course's subject prefix.
>
> **Example**: MATH

> `.course_number`
>
> **Type**: string
>
> The course's official number.
>
> **Example**: 2417

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

> `.enrollment_reqs`
> 
> **Type**: string
> 
> The course's enrollment requirements.
> 
> **Example**: " A minimal placement score of 85% on ALEKS math placement exam or ..."

> `.school`
>
> **Type**: string
>
> The course's school. 
>
> **Example**: "School of Natural Sciences and Mathematics"

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
> A collection of all course requirements that must be met while a student enrolls in a section of this course.

> `.co_or_pre_requisites`
>
> **Type**: CollectionRequirement
>
> A collection of all course requirements that must be met either before or while a student enrolls in a section of this course.

> `.sections`
>
> **Type**: Array<ObjectId>
>
> A list of all sections that are instances of this course.

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

> `.catalog_year`
> 
> **Type**: string
> 
> The catalog year of the course. This is the year in which this instance of the course was published in the corresponding course catalog.
> Only the last two digits of the year are used and the first two digits are assumed to be "20".
> 
> **Example**: "19"

> `.attributes`
>
> **Type**: Object
>
> Space for additional data describing the course not listed elsewhere.

# Course
## Overview
The `Course` object represents a course available at the University of Texas at Dallas. A `Course` should not be confused with a `Section` which is the actual instantiation of a `Course` with a professor and dedicated meeting times. 

## Object Representation
```ts
Course = {
	"id": number,
	"course_number": string, 
	"subject_prefix": string, 
	"title": string, 
	"description": string,
	"school": string, 
	"credit_hours": number, 
	"class_level": string, 
	"activity_type": string, 
	"grading": string,
	"internal_course_number": string,
	"prerequisite_courses": Collection,
	"corequisite_courses": Collection,
	"postrequisite_courses": Collection,
	"membership": Array<Collection>,
	"attributes": {},
}
```

## Properties
> **`.id`**
>
> **Type**: number
> 
> The Firestore database reference number for the `Course` object. 


> **`.course_number`**
>
> **Type**: string
> 
> The course's official number. 
>
> **Example**: 2417


> **`.subject_prefix`**
>
> **Type**: string
> 
> The course's subject prefix.
>
> **Example**: MATH

> **`.title`**
>
> **Type**: string
> 
> The course's title.
>
> **Example**: Calculus I

> **`.description`**
>
> **Type**: string
> 
> The course's description.
>
> **Example**: Calculus I

> **`.school`**
>
> **Type**: string
> 
> The course's school. 
>
> **Example**: Natural Sciences and Mathematics

> **`.credit_hours`**
>
> **Type**: number
> 
> The number of credit hours awarded by successful completion of the course. 
>
> **Example**: 4

> **`.class_level`**
>
> **Type**: string
> 
> The level of education that this course corresponds to. 
>
> **Example**: Undergraduate

> **`.activity_type`**
>
> **Type**: string
> 
> The type of class this course corresponds to.
>
> **Example**: Lecture

> **`.grading`**
>
> **Type**: string
> 
> The grading status of this course.
>
> **Example**: Graded

> **`.internal_course_number`**
>
> **Type**: string
> 
> The internal (university) number used to reference this course.
>
> **Example**: 008613

> **`.prerequisite_courses`**
>
> **Type**: Collection
> 
> A collection of all course requirements that must be met before a student may enroll in a section of this course. 

>**` .corequiste_courses`**
>
> **Type**: Collection
> 
> A collection of all course requirements that must be met either before or while a student enrolls in a section of this course

> **`.postrequisite_courses`**
>
> **Type**: Collection
> 
> A collection of all courses that this course is a prerequisite for.

> **`.membership`**
>
> **Type**: Array\<Collection>
> 
> An array of `Collection` that this course is a part of.  

> **`.attributes`**
>
> **Type**: Object
> 
> Space for additional data describing the course not listed elsewhere. 
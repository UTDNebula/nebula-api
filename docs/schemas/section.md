# Section
## Overview

A `Section` object is the instantiation of a `Course` object with a professor, meeting times, and a grade distribution.

## Object Representation
```ts
Section = {
    "_id": ObjectId,
    "section_number": string,
    "course_reference": ObjectId,
    "section_corequisites": CollectionRequirement,
    "academic_session": AcademicSession,
    "professors": Array<ObjectId>,
    "teaching_assistants": Array<Assistant>,
    "internal_class_number": string,
    "instruction_mode": string,
    "meetings": Array<Meeting>,
    "syllabus_uri": string,
    "grade_distribution": Array<number>,
    "attributes": Object,
}
```

## Properties
> `._id`
> 
> **Type**: ObjectId
> 
> The MongoDB database id for the `Section` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.section_number`
> 
> **Type**: string
> 
> The section's official number.
> 
> **Example**: 002

> `.course_reference`
> 
> **Type**: ObjectId
> 
> An id that points to the course in MongoDB that this section is an instantiation of.
> 
> **Example**:

> `.section_corequisites`
> 
> **Type**: CollectionRequirement
> 
> A collection of all sections that must be taken alongside this section such as specific labs and exam sections.

> `.academic_session`
> 
> **Type**: AcademicSession
> 
> The academic session that the section takes place in.

> `.professors`
> 
> **Type**: Array<ObjectId>
> 
> An array of ids linking to each professor in MongoDB for this section.

> `.teaching_assistants`
> 
> **Type**: Array<Assistant>
> 
> An array of all teaching assistants for this section.

> `.internal_class_number`
> 
> **Type**: string
> 
> The internal (university) number used to reference this section.
> 
> **Example**: 82785

> `.instruction_mode`
> 
> **Type**: string
> 
> The instruction modality for this section.
> 
> **Example**: Traditional

> `.meetings`
> 
> **Type**: Array<Meeting>
> 
> An array of the locations and times that this section meets.

> `.core_flags`
> 
> **Type**: Array<string>
> 
> An array of core requirement codes this section fulfills. 
>
> **Example**: ["020", "050", ...]

> `.syllabus_uri`
> 
> **Type**: string
> 
> A link to the syllabus on the web.
> 
> **Example**: https://dox.utdallas.edu/syl118093

> `.grade_distribution`
> 
> **Type**: Array<number>
> 
> An array of how many students achieved a certain grade in this section decreating in the order of A+, A, A-, etc.
> 
> **Example**: [2, 3, 5, 4, 3, ...]

> `.attributes`
> 
> **Type**: Object
> 
> Space for additional data describing the section not listed elsewhere.

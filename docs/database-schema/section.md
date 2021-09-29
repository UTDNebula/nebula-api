# Section
## Overview
A `Section` object is the instantiation of 	a `Course` object with a professor, meeting times, and a grade distribution. 

## Object Representation
```ts
Section = {
	"id": number,
	"section_number": string,
	"course_reference": number,
	"section_corequisites": Collection, 
	"academic_session": {
		"start_date": date,
		"end_date": date,
	},
	"professor_reference": number,
	"teaching_assistants": Array<strings>,
	"internal_class_number": string,
	"instruction_mode": string,
	"meetings": Array
	"syllabus_reference": number,
	"grade_distribution": Array<number>,
	"attributes": {},
}
```

## Properties
> **`.id`**
>
> **Type**: number
> 
> The Firestore database reference number for the `Section` object. 

> **`.section_number`**
>
> **Type**: string
> 
> The section's official number.
>
> **Example**: 002

> **`.course_reference`**
>
> **Type**: number
> 
> A reference number that points to the course this section is an instantiation of. 

> **`.section_corequisites`**
>
> **Type**: Collection
> 
> A collection of all sections that must be taken alongside this section such as specific labs and exam sections.

> **`.academic_session`**
>
> **Type**: object
> 
> The dates over which this section was available. 
> 
> > **`.start_date`**
> > 
> > **Type**: Date
> >
> > The start date of the academic session.
>
> > **`.end_date`**
> >
> > **Type**: Date
> >
> > The end date of the academic session.

> **`.professor_reference`**
>
> **Type**: number
> 
> A reference number pointing to this section's professor

> **`.teaching_assistants`**
>
> **Type**: Array\<string>
> 
> An array of all teaching assistant names for this section.

> **`.internal_class_number`**
>
> **Type**: string
> 
> The internal (university) number used to reference this section.
>
> **Example**: 82785

> **`.instruction_mode`**
>
> **Type**: string
> 
> The instruction modality for this section.
>
> **Example**: Traditional

> **`.meetings`**
>
> **Type**: Array
> 
> An array of the locations and times that this section meets. 

> **`.section_number`**
>
> **Type**: string
> 
> The section's official number.
>
> **Example**: 002

> **`.syllabus_reference`**
>
> **Type**: number
> 
> A number referencing this section's syllabus. 

> **`.grade_distribution`**
>
> **Type**: Array\<number>
> 
> An array of how many students achieved a certain grade in this section decreasing in the order of A+, A, A-, etc. 

> **`.attributes`**
>
> **Type**: Object
> 
> Space for additional data describing the section not listed elsewhere. 
# Professor

## Overview

A `Professor` represents a professor employed at the University of Texas at Dallas.

## Object Representation

```ts
Professor = {
    "_id": ObjectId,
    "Summary": string,
    "Location": string,
    "StartTime": time.Time,
    "EndTime": time.Time,
    "Description":string,
	"EventType":Array<string>, 
	"TargetAudience":Array<string>,          
	"Topic":Array<string>,          
	"EventTags":Array<string>,          
	"EventWebsite":Array<string>,           
	"Department":Array<string>,        
	"ContactName":string,          
	"ContactEmail":string,         
	"ContactPhoneNumber":string,  
}
```

## Properties

> `._id`
>
> **Type**: ObjectId
>
> The MongoDB database id for the `Professor` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.summary`
>
> **Type**: string
>
> The title summary of the event
>
> **Example**: Strong HIIT

> `.location`
>
> **Type**: string
>
> Thelocation of the event
>
> **Example**: 800 W. Campbell Road, Richardson, Texas 75080-3021

> `.start_time`
>
> **Type**: time.Time
>
> The starting time of the event in RFC3339 format
>
> **Example**: 2023-10-31T06:30:00-05:00

> `.end_time`
>
> **Type**: time.Time
> The ending time of the event in RFC3339 format
>
> **Example**: 2023-10-31T06:30:00-05:00

> `.description`
>
> **Type**: string
>
> The longer description of the event and its details
>
> **Example**:  "Strong Nation combines body weight, muscle conditioning, cardio, and plyometric training moves synced to original music that has been designed to match every single lunge, squat, and burpee. Maximize your burn with the ultimate 60 minute, four part music inspired HIT workout. Stop counting the reps and start training to the beat! "

> `.event_type`
>
> **Type**: Array<string>
>
> The type of event 
>
>**Example**:  `["Sports & Recreation"]`

> `.target_audience`
>
> **Type**: Array<string>
>
> The target audences for the event
>
> **Example**: `["Undergraduate Students","Graduate Students",	"International Students"]`

> `.topic`
>
> **Type**: Array<string>
>
> The overall topics relating to the event
>
> **Example**: `["Health \u0026 Wellness"]`

> `.event_tags`
>
> **Type**: Array<string>
>
> Tags to help query the events
>
> **Example**: `["Lectures and workshops"]`

> `.event_website`
>
> **Type**: string
>
> The website of the organizer
>
> **Example**: "https://organizer.utdallas.edu/something/something"

> `.department`
>
> **Type**: Array<string>
>
> The department hosting the event
>
> **Example**: `["University Recreation"]`

> `.contact_name`
>
> **Type**: string
>
> Name of the contact for the event
>
> **Example**: "Department Name"

> `.contact_email`
>
> **Type**: string
>
> The email adrress of the contact for the event
>
> **Example**: "department@utdallas.edu"

> `.contact_phone_number`
>
> **Type**: string
>
> The phone number of the contact for the event
>
> **Example**:"4104489458"

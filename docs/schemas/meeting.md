# Meeting
## Overview

A 'Meeting' represents a recurring meeting. This schema can represent both recurring meetings and single meetings. Meetings occur repeatedly on the specified days of the week during a period. Non-recurring meetings should have the `start_date` equal to the `end_date`.

## Object Representation
```ts
ModalityType = "pending" | "traditional" | "hybrid" | "flexible" | "remote" | "online"

Meeting = {
    "start_date": string,
    "end_date": string,
    "meeting_days": Array<string>,
    "start_time": string,
    "end_time": string,
    "modality": string,
    "location": Location,
}
```

## Properties
> `.start_date`
>
> **Type**: string
>
> The start date of a meeting.
>
> **Example**: January 18, 2022

> `.end_date`
>
> **Type**: string
>
> The end date of a meeting.
>
> **Example**: May 5, 2022

> `.meeting_days`
>
> **Type**: Array<string>
>
> A list of all days the meeting occurs during the time period.
>
> **Example**: ["Monday", "Wednesday"]

> `.start_time`
>
> **Type**: string
>
> The time the meeting starts on each meeting day.
>
> **Example**: "10:00am"

> `.end_time`
>
> **Type**: string
>
> The time a meeting ends on each meeting day.
>
> **Example**: "11:15am"

> `.modality`
>
> **Type**: string
>
> The modality of the meeting following the modality types in [UTD's CourseBook](https://coursebook.utdallas.edu/modalities).
>
> **Example**: traditional

> `.location`
>
> **Type**: Location
>
> The location of the meeting.

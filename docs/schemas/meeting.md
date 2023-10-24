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
    "modality": ModalityType,
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
> **Example**: `2022-08-22T00:00:00-05:00`

> `.end_date`
>
> **Type**: string
>
> The end date of a meeting.
>
> **Example**: `2022-12-08T00:00:00-06:00`

> `.meeting_days`
>
> **Type**: Array\<string>
>
> A list of all days the meeting occurs during the time period.
>
> **Example**: `["Monday", "Wednesday"]`

> `.start_time`
>
> **Type**: string
>
> The time the meeting starts on each meeting day.
>
> **Example**: `0000-01-01T16:00:00-05:50`

> `.end_time`
>
> **Type**: string
>
> The time a meeting ends on each meeting day.
>
> **Example**: `0000-01-01T17:15:00-05:50`

> `.modality`
>
> **Type**: ModalityType
>
> The modality of the meeting following the modality types in [UTD's CourseBook](https://coursebook.utdallas.edu/modalities).
>
> NOTE: All observed entries of this have been blank, use with caution.
>
> **Example**: traditional

> `.location`
>
> **Type**: Location
>
> The location of the meeting.

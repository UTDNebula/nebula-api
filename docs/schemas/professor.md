# Professor

## Overview

A `Professor` represents a professor employed at the University of Texas at Dallas.

## Object Representation

```ts
Professor = {
    "_id": ObjectId,
    "first_name": string,
    "last_name": string,
    "titles": Array<string>,
    "email": string,
    "phone_number": string,
    "office": Location,
    "profile_uri": string,
    "image_uri": string,
    "office_hours": Array<Meeting>,
    "sections": Array<ObjectId>,
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

> `.first_name`
>
> **Type**: string
>
> The professor's first name.
>
> **Example**: John

> `.last_name`
>
> **Type**: string
>
> The professor's last name.
>
> **Example**: Doe

> `.titles`
>
> **Type**: Array\<string>
>
> The professor's titles.
>
> **Example**: `["Senior Mathematics Lecturer"]`, `["Lars Magnus Ericsson Chair", "Dean – Erik Jonsson School of Engineering and Computer Science"]`

> `.email`
>
> **Type**: string
>
> The professor's email address.
>
> **Example**: `xxx555555@utdallas.edu`

> `.phone_number`
>
> **Type**: string
>
> The professor's phone number.
>
> **Example**: `123-456-7890`

> `.office`
>
> **Type**: Location
>
> The location of the professor's office.

> `.profile_uri`
>
> **Type**: string
>
> A hyperlink pointing to the professor's official university profile.
>
> **Example**: `https://profiles.utdallas.edu/stephanie.adams`

> `.image_uri`
>
> **Type**: string
>
> A link to the image used for the professor on the professor's official university profile.
>
> **Example**: `https://profiles.utdallas.edu/storage/media/2384/conversions/Adams-headshot-1-medium.jpg`

> `.office_hours`
>
> **Type**: Array\<Meeting>
>
> A list of all office hours of the professor.

> `.sections`
>
> **Type**: Array\<ObjectId>
>
> A list of references to sections a professor is currently teaching or has taught. This will be sorted in increasing order with respect to `end_date` in the section's `academic_session`.

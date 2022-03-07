# Professor
## Overview

A `Professor` represents a professor employed at the University of Texas at Dallas.

## Object Representation
```ts
Professor = {
    "_id": ObjectId,
    "first_name": string,
    "last_name": string,
    "title": string,
    "email": string,
    "phone_number": string,
    "office": Location,
    "profile_uri": string,
    "office_hours": Array<Meeting>,
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

> `.title`
>
> **Type**: string
>
> The professor's title.
>
> **Example**: Senior Mathematics Lecturer

> `.email`
>
> **Type**: string
>
> The professor's email address.

> `.phone_number`
>
> **Type**: string
>
> The professor's phone number.

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

> `.office_hours`
>
> **Type**: Array<Meeting>
>
> A list of all office hours of the professor.

# Organization
## Overview

An `Organization` represents a student organization at the University of Texas at Dallas.

## Object Representation
```ts
Organization = {
    "_id": ObjectId,
    "title": string,
    "description": string,
    "categories": Array<string>,
    "president_name": string,
    "emails": Array<string>,
    "picture_data": string,
}
```

## Properties
> `._id`
>
> **Type**: ObjectId
>
> The MongoDB database id for the `Organization` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.title`
>
> **Type**: string
>
> The organization's title.
>
> **Example**: "200percent"

> `.description`
>
> **Type**: string
>
> The organization's description.
>
> **Example**: "The purpose of 200percent, as a K-pop performance group, is to ..."

> `.categories`
>
> **Type**: Array<string>
>
> The organization's categories.
>
> **Example**: ["Cultural", "Arts and Music", "Special Interest"]

> `.president_name`
>
> **Type**: string
>
> The organization's president.
> 
> **Example**: "Cheri Tang"

> `.emails`
>
> **Type**: Array<string>
>
> The organization's contact emails.
> 
> **Example**: ["ctt180000@utdallas.edu"]

> `.picture_data`
>
> **Type**: string
>
> The raw data of the organization's logo image, encoded in base64.

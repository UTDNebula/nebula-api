# Degree
## Overview

A `Degree` represents either a major, minor, or concentration received from The University of Texas at Dallas.

## Object Representation
```ts
DegreeSubtype = "major" | "minor" | "concentration"

Degree = {
    "_id": ObjectId,
    "subtype": DegreeSubtype,
    "name": string,
    "year": string,
    "abbreviation": string,
    "minimum_credit_hours": number,
    "requirements": CollectionRequirement,
}
```

## Properties
> `._id`
>
> **Type**: ObjectId
>
> The MongoDB database id for the `Degree` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.subtype`
>
> **Type**: DegreeSubtype
>
> The subtype of degree that this object represents.
>
> **Example**: Major

> `.name`
>
> **Type**: string
>
> The full name of the degree.
>
> **Example**: Bachelor of Science in Computer Science

> `.year`
> 
> **Type**: string
>
> The academic year to which this degree corresponds to.
>
> **Example**: 2021-2022

> `.abbreviation`
>
> **Type**: string
>
>The abbreviation of the degree.
>
> **Example**: B.S. in Computer Science

> `.minimum_credit_hours`
>
> **Type**: number
>
> The minimum credit hours required for the degree, which can be found on the UTD catalog.
>
> **Example**: 124

> `.requirements`
>
> **Type**: CollectionRequirement
>
> A collection of requirements required to satisfy the degree.

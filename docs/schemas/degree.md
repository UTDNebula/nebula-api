# Degree

## Overview

A `Degree` represents either a major, minor, or concentration received from The University of Texas at Dallas.

## Object Representation

```ts
DegreeSubtype = "major" | "minor" | "concentration" | "prescribed double major" | "certificate" | "track"

Degree = {
    "_id": ObjectId,
    "subtype": DegreeSubtype,
    "school": string,
    "name": string,
    "year": string,
    "abbreviation": string,
    "minimum_credit_hours": number,
    "catalog_uri": string,
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
> **Example**: `major`

> `.school`
>
> **Type**: string
>
> The school that the `degree` belongs to.
>
> **Example**: `School of Natural Sciences and Mathematics`

> `.name`
>
> **Type**: string
>
> The full name of the degree.
>
> **Example**: `Bachelor of Science in Computer Science`

> `.year`
>
> **Type**: string
>
> The academic year to which this degree corresponds to.
>
> **Example**: `2021-2022`

> `.abbreviation`
>
> **Type**: string
>
>The abbreviation of the degree.
>
> **Example**: `B.S. in Computer Science`

> `.minimum_credit_hours`
>
> **Type**: number
>
> The minimum credit hours required for the degree, which can be found on the UTD catalog.
>
> **Example**: 124

> `.catalog_uri`
>
> **Type**: string
>
> A link to the listing of the degree and its requirements in the UTD catalog.
>
> **Example**: `https://catalog.utdallas.edu/2021/undergraduate/programs/ah/philosophy`

> `.requirements`
>
> **Type**: CollectionRequirement
>
> A collection of requirements required to satisfy the degree.

# Assistant

## Overview

An 'Assistant' represents a teaching assistant at UT Dallas.

## Object Representation

```ts
Assistant = {
    "first_name": string,
    "last_name": string,
    "role": string,
    "email": string,
}
```

## Properties

> **Type**: ObjectId
>
> The MongoDB database id for the `Assistant` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.first_name`
>
> **Type**: string
>
> The first name of the assistant.
>
> **Example**: `John`

> `.last_name`
>
> **Type**: string
>
> The last name of the assistant.
>
> **Example**: `Doe`

> `.role`
>
> **Type**: string
>
> The role of the assistant.
>
> **Example**: `Teaching Assistant`

> `.email`
>
> **Type**: string
>
> The email address to contact the assistant.
>
> **Example**: `xxx555555@utdallas.edu`

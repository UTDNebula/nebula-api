# Academic Session

## Overview

An `AcademicSession` represents the time period in which courses takes place.

## Object Representation

```ts
AcademicSession = {
    "name": string,
    "start_date": string,
    "end_date": string,
}
```

## Properties

> `.name`
>
> **Type**: string
>
> The name of the academic session in question.
>
> **Examples**: `22S`, `18F`, `23U`

> `.start_date`
>
> **Type**: string
>
> The date of classes starting in the academic session.
>
> **Example**: `2022-01-18T00:00:00-06:00`

> `.end_date`
>
> **Type**: string
>
> The date of classes ending in the academic session.
>
> **Example**: `2022-05-13T00:00:00-05:00`

# Credit

## Overview

The `Credit` object represents an amount of 'semester credit hours' given by The University of Texas at Dallas. A `Credit` should not be confused with a `Course` as semester credit hours serve only to fulfill credit hour requirements.

## Object Representation

```ts
Credit = {
  "category": string,
  "credit_hours": number,
};
```

## Properties

> `.category`
>
> **Type**: string
>
> The catergory of the credit hours.
> If there is no category associated with the credit, the value is "general".
> "free" is a valid category.
>
> **Example**: "Geosciences", "Business Computer Information Systems", "free"

> `.credit_hours`
>
> **Type**: number
>
> The number of credit hours.
>
> **Example**: 3

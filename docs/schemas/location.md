# Location

## Overview

A location on the UT Dallas campus.

## Object Representation

```ts
Location = {
    "building": string,
    "room": string,
    "map_uri": string,
}
```

## Properties

> `.building`
>
> **Type**: string
>
> The building of the location.
>
> **Examples**: `SLC`, `ECSW`

> `.room`
>
> **Type**: string
>
> The room of the location.
>
> **Examples**: `2.203`, `1.315`

> `.map_uri`
>
> **Type**: string
>
> A hyperlink to the UTD room locator.
>
> **Example**: `https://locator.utdallas.edu/SLC_2.203`

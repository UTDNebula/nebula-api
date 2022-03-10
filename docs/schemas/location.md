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
> **Example**: "SLC", "ONLINE"

> `.room`
>
> **Type**: string
>
> The room of the location.
>
> **Example**: "2.203", "ONLINE"

> `.map_uri`
>
> **Type**: string
>
> A hyperlink to the UTD room locator.
>
> **Example**: https://locator.utdallas.edu/SLC_2.203

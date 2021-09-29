# Professor
## Overview
A `Professor` represents a professor who has taught or is currently teaching classes at the University of Texas at Dallas. 

## Object Representation
```ts
Professor = {
	"id": number,
	"first_name": string,
	"last_name": string,
	"title": string,
	"email": string,
	"phone_number": string,
	"office": string, 
	"profile_link": string,
}
```

## Properties
> **`.id`**
>
> **Type**: number
> 
> The Firestore database reference number for the `Professor` object. 

> **`.first_name`**
>
> **Type**: string
> 
> The professor's first name.

> **`.last_name`**
>
> **Type**: string
> 
> The professor's last name.

> **`.title`**
>
> **Type**: string
> 
> The professor's title. 
>
> **Example**: Senior Mathematics Lecturer

> **`.email`**
>
> **Type**: string
> 
> The professor's email address. 

> **`.phone_number`**
>
> **Type**: string
> 
> The professor's phone number. 

> **`.office`**
>
> **Type**: string
> 
> The location of the professor's office. 

> **`.profile_link`**
>
> **Type**: string
> 
> A hyperlink pointing to the professor's official university profile. 




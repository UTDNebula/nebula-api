# Collection
## Overview
A `Collection` object represents a grouping of  `Course` objects and other `Collection` objects. A `Collection` is best used to describe a complex grouping of courses where only certain conditions must be met such as in the case of prerequisite or corequisite courses. 

A `Collection` is best understood as a conditional statement object in that a `Collection` specifies how many of its own elements must evaluate to true for itself to evaluate to true. 

For example, in the case that a student is required to succesfully complete one of several courses to satisfy another course's prerequisites, we can represent it with a collection. Using MATH 2418 (Linear Algebra) as an example, it requires that either MATH 2306, MATH 2413, or MATH 2417 be completed. The `Collection` object in this case would be composed of the three courses and specify that only 1 out of the 3 total elements must be completed for the `Collection` itself to evaluate to true. 

Similarly, in the case that a student is required to sucessfully complete multiple courses as prerequisites, we can represent that prerequisite as a collection of courses of which N out of N courses must be completed for the `Collection `object to evaluate to true. 

If more complex collections are required, you can nest collections and courses to represent any finite requirements. 

For example, if a course requires that you take courses (A and B) or (C and D) as prerequisites, then this can be represented as two collections of (A and B) and (C and D) which require 2 out of 2 courses be completed. These two collections can then be put inside of another collection of which only 1 out of the 2 collections must evaluate to true for it to evaluate to true.

## Object Representation
```ts
Collection = {
	"id": number,
	"name": string,
	"type": string,
	"required": number,
	"total": number,
	"options": [],
}
```

## Properties
> **`.id`**
>
> **Type**: number
> 
> The Firestore database reference number for the `Collection` object. 

> **`.name`**
>
> **Type**: string
> 
> The name of the collection.
>
> **Example**: Minor in Mathematics

> **`.type`**
>
> **Type**: string
> 
> The type of course/collection grouping that is meant to be represented by this collection.
>
> **Example**: Prerequisites, Corequisites, Degree, Major, Minor, Concentration

> **`.required`**
>
> **Type**: number
> 
> The number of options in the collection that must be satisfied for the collection itself to evaluate to true.
>
> **Example**: 3

> **`.total`**
>
> **Type**: number
> 
> The total number of options within the collection. 
>
> **Example**: 5

> **`.options`**
>
> **Type**: Array
> 
> An array of all courses and collections that compose this collection.





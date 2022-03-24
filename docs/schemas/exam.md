# Exam

## Overview

An `Exam` represents a credit-bearing or placement examination at the University of Texas at Dallas.

## Object Representation

```ts
ExamType = "AP" | "ALEKS" | "CLEP" | "IB" | "CS placement"

type Outcome = {
  requirement: Requirement,
  outcome: Array<Array<ObjectId>>
}

abstract Exam = {
  "_id": ObjectId,
  "type": ExamType,
}

APExam extends Exam = {
  "_id": ObjectId,
  "type": "AP",
  "name": string,
  "yields": Array<Outcome>
}

ALEKSExam extends Exam = {
  "_id": ObjectId,
  "type": "ALEKS",
  "placement": Array<Outcome>
}

CLEPExam extends Exam = {
  "_id": ObjectId,
  "type": "CLEP",
  "name": string,
  "yields": Array<Outcome>
}

IBExam extends Exam = {
  "_id": ObjectId,
  "type": "IB",
  "name": string,
  "level": string,
  "yields": Array<Outcome>
}

CSPlacementExam extends Exam = {
  "_id": ObjectId,
  "type": "CS placement",
  "yields": Array<Outcome>
}
```

## Properties

> `._id`
>
> **Type**: ObjectId
>
> The MongoDB database id for the `Exam` object.
>
> **Example**: ObjectId("61ebbb126e3659537e8a14d6")

> `.type`
>
> **Type**: ExamType
>
> The type of exam object this object represents.
>
> **Example**: AP

> `.yields`
>
> **Type**: Array<Outcome>
>
> An array of `Outcomes` for which the credit for the `Course` or `Credit` is received. Does not include placement, only actual credit.
>
> **Example**: [{requirement: ExamRequirement, [[MATH2312._id, MATH 1325._id], [MATH2312._id, MATH2413._id]]}]

> `.placement`
>
> **Type**: Array<Outcome>
>
> An array of `Outcomes` for which the placement into the `Course` is earned. Does not include credit, only placement into the course.
>
> **Example**: [{requirement: ExamRequirement, [[MATH1325._id]]}]

> `.name`
>
> **Type**: string
>
> The name of the exact exam being taken.
>
> **Example**: "Macroeconomics", "American History I: Early Colonization to 1877", "Environmental Systems
> and Societies"

> `.level`
>
> **Type**: string
>
> The level of the IB exam.
>
> **Example**: Standard, Higher

## Outcome

> `.requirement`
>
> **Type**: Requirement
>
> The requirement to achieve the associated outcome
>
> **Example**: ExamRequirement.minimum_score === 4, MajorRequirement.major === "Physics"

> `.outcome`
>
> **Type**: Array<Array<ObjectId>>
>
> The set of sets of `Course`s and/or `Credit`s which can result (awarded/placed into) should the requirement be met.
> The outer array contains the possible choices.
>
> **Example**: [[MATH2312._id, MATH 1325._id], [MATH2312._id, MATH2413._id]]

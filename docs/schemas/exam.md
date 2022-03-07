# Exam
## Overview
  
An `Exam` represents a credit-bearing or placement examination at the University of Texas at Dallas.
  
## Object Representation
```ts
ExamType = "AP" | "ALEKS" | "CLEP" | "IB" | "CS placement"
  
abstract Exam = {
  "_id": ObjectId,
  "type": ExamType,
  "yields": Dictionary<int, Array<ObjectId>>,
}
  
APExam extends Exam = {
  "_id": ObjectId,
  "type": "AP",
  "name": string,
  "yields": Dictionary<int, Array<ObjectId>>,
}
  
ALEKSExam extends Exam = {
  "_id": ObjectId,
  "type": "ALEKS",
  "yields": {},
}
  
CLEPExam extends Exam = {
  "_id": ObjectId,
  "type": "CLEP",
  "name": string,
  "yields": Dictionary<int, Array<ObjectId>>,
}
 
IBExam extends Exam = {
  "_id": ObjectId,
  "type": "IB",
  "name": string,
  "level": string,
  "yields": Dictionary<int, Array<ObjectId>>,
}
  
CSPlacementExam extends Exam = {
  "_id": ObjectId,
  "type": "CS placement",
  "yields": {70: [CS1336._id, CS1136._id]},
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
> **Type**: Dictionary<int, Array<ObjectId>>
> 
> Maps scores to credit received as a reference to the `Course`. Does not include placement, only actual credit.
>
> **Example**: {}, {70: [CS1336._id, CS1136._id]}
  
> `.name`
> 
> **Type**: string
> 
> The name of the exact exam being taken.
>
> **Example**: "Macroeconomics", "American History I: Early Colonization to 1877", "Environmental Systems
and Societies"
  
> `.level`
> 
> **Type**: string
> 
> The level of the IB exam.
>
> **Example**: Standard, Higher

# IMPLEMENTING THE GRAPHQL RESOLVER

Most of the GraphQL frameworks, including gqlgen(Go), follow the **schema-first** approach & auto-generate the boilerplate code resolving the GraphQL schema.

Our responsibility is:
- Define or update the schema 
- Implement the resolver to fetch (aggregate) data from the database

## Define or update the schema (model)
Define or update schema in ```schema.graphqls```. You should rely on ```../api/schema/``` for objects we already have in the database.

Run 
```
cd graphql
make generate
# Or .\build.bat generate if you are using Windows??
```
This starts the boilerplate auto-generation.

After running it you will see: 
- New schema in Go struct types in ```graph/model/models_gen.go```
- Changes in ```graph/generated.go```

**IMPORTANT**:
- DO NOT edit ```models_gen.go``` directly
- DO NOT edit ```generated.go``` directly

If you want to modify the schema later on, update the schema in ```schema.graphqls``` and run the command again.

The content in ```generated.go``` essentially: 
- How each field in GraphQL is mapped to Go structs
- Resolver wiring and execution configuration


### Note:
Objects that have **ID REFERENCE** (course, section, etc) require separation of concerns.

This means: 
- One schema generated in ```models-gen.go``` on GraphQL side
- One schema defined in ```db_models.go``` on MongoDB side. 

For example, ```Course``` interacts with GraphQL and ```DBCourse``` interacts with the Mongo.

Currently, ```Course``` & ```DBCourse``` are identical, but once we start implementing reference resolver later, there will be difference and separation of concerns will come into plays.

For other objects without **ID REFERENCE** (astra events, rooms, etc), we can combine DB object and GraphQL object as one using ```@goTag```. Since we won't ever implement reference resolver on these objects, there's no need for different schemas.

## Implement the Resolver
When you run
```
make generate
```
Resolver methods for new objects will be auto-generated in 
```
schema.resolvers.go
```
This happens because gqlgen requires you to implement resolvers to all defined schemas (no missing), and by default all resolvers take place in that file.

**OUR STRUCTURE:**
We want more structured organization, so we will move resolver methods into dedicated files per object

For example, 
- ```course.resolvers.go``` contains all course-related resolvers
- ```astra.resolvers.go``` contains all astra-related resolvers.

**ISSUE:**
- When you run ```make generate``` again, the existing methods will be generated back to ```schema.resolvers.go```.
- Remove the re-generated methods as a temporary solution.

Implementation of GraphQL resolvers are very similar to REST controllers implementation, which includes querying and aggregating MongoDB and parse them into the schemas.
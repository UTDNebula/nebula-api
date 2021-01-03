# Comet Data Service API
The API allows any authenticated user to get various UTD data including catalog
information and more.

## Roadmap
- Catalog Information
  - Basically open data for Coursebook
- Degree plan information (2020Q4)
  - Machine-readable data for course prerequisites
- Course popularity (2021Q1)
  - Data sourced from Comet Planning revealing insights from schedules

## Current API Endpoints

Testing server: `https://comet-data-api.herokuapp.com/`

- `degreeAll`
  - Lists out all scraped information for core, major, and elective requirements
  - Core information is complete, major/elective are in progress
- `degree`
  - params = school (ah, ecs, ...), major (computer-science, latin-american-studies, ...)
  - Returns the degree plan information for one school + major
- `core`, `major`, `elective`
  - params = school, major
  - Returns the core/major/elective requirements for that major
- `prerequisite`
  - params = name (CS 3345, CS 1200, ...)
  - Returns the prerequisites for the given course name in a structured format

## Implementation

### Prerequisite Parser

[Pyparsing](https://github.com/pyparsing/pyparsing) is used for parsing prerequisite strings. This is now deprecated as we have moved to using [Chevrotain](https://github.com/SAP/chevrotain), which is a similar (but slightly more flexible and easier to use) language parser in JavaScript.

## Notes

+ `Procfile` and `runtime.txt` are for Heroku deployment

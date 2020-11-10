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
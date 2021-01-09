# About the Scraper
The scraper is a designed to obtain and consolidate assorted UTD data.

The data scraper collects data from two sources:
- The Catalog
- Coursebook

## Information

`fetchAllCourseNames.py` fetches all course names from UTD catalog.
- Course names are stored in `output/all_course_names.txt` in the format of `XXXX 0000`

`fetchAllMajors.py` fetches all majors and their core, major, and elective requirements.

- Final outputs are stored in `output/major_requirements_processed.txt`
- The JSON file has the different schools as the first layer, then each school has its majors, which then contains their Core, Major, and Elective requirements.

`fetchCoursePrereqs.py` fetches the prerequisites for each course. 

- Final outputs are stored in `output/course_catalog_parsed.txt`
- Each course has its `id`, `name`, `description`, `hours`, `prerequisites`, `inclass`, `outclass`, and `period`.
- Prerequisites are currently **not** formatted; implement `def format_prereq()` in this file to specify proper prerequisite format

## How to Run

1. Setup a virtual environment in `scraper/` directory: `virtualenv env`
2. Install the necessary prerequisites: `pip install -r requirements.txt`
3. Scripts:
    1. `python fetchAllCourseNames.py` to update course names
    2. `python fetchAllMajors.py` to fetch all majors and their degree requirements
    3. `python fetchCoursePrereqs.py` to fetch every course's prerequisites
4. See Information section for more details about output files
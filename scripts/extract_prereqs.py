import re
import copy
import json 

# TODO: update in JS scraper
courses = json.load(open("data/scheduler_info.json"))
updated_courses = []

count = 0
for course in courses:
    if course["id"] == 276:
        print(26)
    matches = re.findall(r"\. (Prerequisite[s]* or Corequisite|Prerequisite|Corequisite)[s]*: (.*?)(?=(\.))", course["description"])
    newCourse = copy.deepcopy(course)
    for match in matches:
        if "or Corequisite" in match[0]:
            newCourse["prerequisiteOrCorequisites"] = match[1]
            count += 1
        elif "Prerequisite" in match[0]:
            newCourse["prerequisites"] = match[1]
            count += 1
        elif "Corequisite" in match[0]:
            newCourse["corequisites"] = match[1]
            count += 1
    newCourse["description"] = re.sub(r"\. (Prerequisite[s]* or Corequisite|Prerequisite|Corequisite)[s]*: (.*?)(?=(\.))", "", course["description"])
    newCourse["course"] = course["subjectId"] + " " + course["number"]
    updated_courses.append(newCourse)

with open("data/scheduler_prereq.json", 'w') as f:
    json.dump(updated_courses, f);
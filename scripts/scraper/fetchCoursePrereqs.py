import re
import json
import requests
from bs4 import BeautifulSoup as bs4

##### Start scraping section #####


def getCourseWeb(courseName):
    try:
        courseName = courseName.replace(' ', '').lower()
        page = requests.get(
            'https://catalog.utdallas.edu/2020/undergraduate/courses/' + courseName)
        soup = bs4(page.text, 'html.parser')
        all_info = soup.select('#bukku-page > p')[0]
        title = all_info.select('.course_title')[0]
        return {'title': title.text, 'description': all_info.text}
    except:
        print('Error: ' + courseName)
        return {}


def retrieve_web():
    filename = 'output/course_catalog_complete.txt'
    course_info = {}
    with open('output/all_course_names.txt') as course_list:
        i = 0
        for course in course_list:
            course = course.rstrip()
            if course not in course_info:
                course_info[course] = getCourseWeb(course)
                print(course)
                if i % 10 == 0:
                    json.dump(course_info, open(filename, 'w'))
                i += 1
        json.dump(course_info, open(filename, 'w'))


##### Start processing (not prerequisites) section #####

checks = {
    'Prerequisite[s]*: ': 'prerequisites',
    'Corequisite[s]*: ': 'corequisites',
    'Recommended Corequisite[s]*: ': 'recommendedCorequisites',
    'Prerequisite[s]* or Corequisite[s]*: ': 'PrerequisitesOrCorequisites'
}


def getCourse(id, courseName, title, all_info):
    result = {}
    try:
        courseName = courseName.replace(' ', '').lower()
        hours = re.findall(r'[a-z]+[0-9]([a-z0-9])', courseName)[0]
        description = re.sub(
            r'.*?credit hour[s]*\)', '', all_info).strip()
        information = re.sub(
            r'([A-Za-z][A-Za-z ]+requisite[s]*: .*?)\.', '', description)
        inclass = 'N/A'
        outclass = 'N/A'
        period = 'N/A'

        matches = re.findall(r'\((.*?)-(.*?)\) ([A-Z])', information)
        if matches:
            inclass = matches[0][0]
            outclass = matches[0][1]
            period = matches[0][2]
            information = re.sub(r'\((.*?)-(.*?)\) ([A-Z])', '', information)
        else:
            print(desp)
            print('-------------------------')

        prerequisites = re.findall(
            r'([A-Za-z][A-Za-z ]+requisite[s]*: .*?)\.', description)

        result = {
            'id': id,
            'name': title,
            'hours': hours,
            'description': information,
            'inclass': inclass,
            'outclass': outclass,
            'period': period,
            'prerequisites': prerequisites
        }

    except:
        print('error!')
    return result


def parse_web_data():
    filename = 'output/course_catalog_parsed.txt'
    course_info = {}
    with open('output/course_catalog_complete.txt') as course_data:
        course_complete = json.loads(course_data.read())
        i = 0
        for course in course_complete:
            course = course.rstrip()
            current_data = course_complete[course]
            course_info[course] = getCourse(i,
                                            course, current_data['title'], current_data['description'])
            if i % 10 == 0:
                # dump periodically to save
                json.dump(course_info, open(filename, 'w'))
            i += 1
        json.dump(course_info, open(filename, 'w'))

##### Start processing prerequisites section (TODO) #####


def format_prereq(name):
    # TODO
    return name


def format_prereqs(course):
    for i, prereq in enumerate(course['prerequisites']):
        processed = format_prereq(prereq)
        course['prerequisites'][i] = processed
    return course


def process_prerequisites():
    course_info = {}
    filename = 'output/course_catalog_processed.txt'
    with open('output/course_catalog_parsed.txt') as course_data:
        course_info = json.loads(course_data.read())
        i = 0
        for course in course_info:
            course_info[course] = format_prereqs(course_info[course])
            if i % 10 == 0:
                json.dump(course_info, open(filename, 'w'))
            i += 1
        json.dump(course_info, open(filename, 'w'))


def main():
    retrieve_web()
    parse_web_data()
    # process_prerequisites()


if __name__ == '__main__':
    main()

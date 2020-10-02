from bs4 import BeautifulSoup as bs4
import requests
import re
import json

##### Start scraping section #####

# returns the level of a node using class name of format xind-L, where L is the level
def get_level(node):
    try:
        match = re.findall(r"xind-([0-9])", ''.join(node['class']))
        if match and len(match) > 0:
            return int(match[0][0])
    except:
        pass
    return 10000

# remove tooltips from text
def process(child):
    text = child.text
    if child.find('sup'):
        text = text[:-len(child.find('sup').text)]
    return text

# recursively store degree plan structure using website layout
def recursion(all_children, index):
    level = get_level(all_children[index])
    result = {}
    i = index + 1

    while i < len(all_children):
        child = all_children[i]
        child_level = get_level(child)
        # only process immediate children
        if child_level == level + 1:
            text = process(child)
            if text == "":
                pass
            elif i+1 < len(all_children) and (all_children[i+1].get('class') != None and "catreq-cont" in all_children[i+1].get('class')):
                # Example:
                # MATH 3323
                # or MATH 3333
                # or MATH 3343
                # Will be combined into {"Select one of the following": [MATH 3323, MATH 3333, MATH 3343]}
                options = [process(child)]
                while i+1 < len(all_children) and (all_children[i+1].get('class') != None and "catreq-cont" in all_children[i+1].get('class')):
                    # add all "or ..." into list
                    newText = process(all_children[i+1])[3:]
                    options.append(newText)
                    i += 1

                # append into results
                text = "Select one of the following"
                result[text] = options
            else:
                # normal path, process children recursively
                result[text] = recursion(all_children, i)
        elif child_level <= level:
            break
        i += 1
    return result


def fetch_major(url):
    # fetch one major using its URL
    page = requests.get(url)
    soup = bs4(page.text, 'html.parser')
    all = soup.select('#bukku-page')[0]
    all_children = all.find_all(recursive=False)

    i = 0
    indices = []
    foundBA = False # some majors have both BA/BS in one URL
    foundBS = False # BA always comes first, so load first result into BA, second into BS
    while i < len(all_children):
        children = all_children[i]
        try:
            # xind-2 is "I. Core Requirement", "II. Major Requirement" and maybe electives title heading
            if 'xind-2' in children['class']:
                text = children.findAll(text=True, recursive=False)
                if len(text) > 0:
                    match = re.findall(r"([I]+)\. (.*)\:", text[0])
                    if match:
                        if foundBA and match[0][0] == "I":
                            indices = [
                                {"index": x["index"], "name": x["name"] + " BA"} for x in indices]
                            foundBS = True
                        foundBA = True
                        if foundBS:
                            indices.append(
                                {"index": i, "name": match[0][1] + " BS"})
                        else:
                            indices.append({"index": i, "name": match[0][1]})
        except:
            # no classes
            pass
        i += 1

    result = {}
    # indices contains all the big headings (I. ... Requirement)
    for root in indices:
        result[root["name"]] = recursion(all_children, root["index"])
    # json.dump(result, open("major_requirements.txt", 'w'))
    return result


def retrieve_all():
    # Scan through program listing, feed all URL into fetchMajor function
    page = requests.get(
        'https://catalog.utdallas.edu/2020/undergraduate/programs')
    soup = bs4(page.text, 'html.parser')
    all_links = soup.find_all('a', href=True)

    all_results = {}
    for link_element in all_links:
        href = link_element['href']
        if href not in all_results:
            # find all urls in the format below
            matches = re.findall(
                r'http://catalog.utdallas.edu/2020/undergraduate/programs/(.*)/(.*)', href)
            # historical-studies is empty for some reason
            if matches and 'historical-studies' not in href:
                print(href)
                major_info = fetch_major(href)
                if major_info:
                    if matches[0][0] not in all_results:
                        # [0][0] is the school, [0][1] is the major
                        all_results[matches[0][0]] = {}
                    all_results[matches[0][0]][matches[0][1]] = major_info
                    json.dump(all_results, open("output/major_requirements.txt", 'w'))

# debug function to test for != 9 cores and no cores
def scan_anomaly():
    all_results = {}
    with open('output/major_requirements.txt') as prior:
        all_results = json.loads(prior.read())
    for school in all_results:
        for major in all_results[school]:
            has_title = False
            for title in all_results[school][major]:
                if title.startswith("Core"):
                    has_title = True
                    count = len(all_results[school][major][title])
                    if count != 9:
                        print("NOT 9: " + school + ", " + major +
                              ", " + title + ": " + str(count))
            if not has_title:
                print("NO CORE: " + school + ", " + major +
                      ", " + title + ": " + str(count))

##### Start processing section #####

# Constants
mapping = {
    "Select one of the following": "Choose 1",
    "And choose one course from the following": "Choose 1",
    "Choose one of the following": "Choose 1",
    "Choose two of the following": "Choose 2",
    "Choose one course from the following": "Choose 1",
    "Choose two courses from the following": "Choose 2",
    "Choose one from the following": "Choose 1",
    "AND one of the following": "Choose 1",
    "Choose three courses from the following": "Choose 3",
    "Select 3 semester credit hours of Related Courses from the following": "Any 3 from following",

    "3 semester credit hours from Mathematics Core": "Any 3 from 020",
    "6 semester credit hours from Communication Core": "Any 6 from 010",
    "6 semester credit hours from Life and Physical Sciences Core": "Any 6 from 030",
    "3 semester credit hours from Language, Philosophy and Culture Core": "Any 3 from 040",
    "3 semester credit hours from Creative Arts Core": "Any 3 from 050",
    "3 semester credit hours Creative Arts Core": "Any 3 from 050",
    "6 semester credit hours from American History Core": "Any 6 from 060",
    "3 semester credit hours from American History Core": "Any 3 from 060",
    "6 semester credit hours from Government/Political Science Core": "Any 6 from 070",
    "3 semester credit hours from Social and Behavioral Sciences Core": "Any 3 from 080",
    "6 semester credit hours from Social and Behavioral Sciences Core": "Any 6 from 080",
    "6 semester credit hours from Component Area Option Core": "Any 6 from 090",
    "one of the following laboratories:": "Choose 1",
    "Students take 9 semester credit hours from any 4000 level course from the list below. Independent Study in Computer Engineering (CE 4V97), Undergraduate Research in Computer Engineering (CE 4V98), or Senior Honors in Computer Engineering (CE 4399) may be used for up to 6 of these hours.": "9 hours",
    "Students pursuing a concentration in one of the following areas should take a minimum of two courses in that area:": "Choose 1 path and Choose 2 courses",
    "Students pursuing the general program take 29 semester credit hours from the list below:": "29 hours",
    "Core Curriculum Requirements": "Core",
    "Major Requirements": "Major",
    "Elective Requirements": "Elective",
    "Elective Requirements BS": "Elective BS",
    "Elective Requirements BA": "Elective BA",

}

starts_map = {
    "Communication": "010",
    "Mathematics": "020",
    "Life and Physical Sciences": "030",
    "Language, Philosophy and Culture": "040",
    "Creative Arts": "050",
    "American History": "060",
    "Government/Political Science": "070",
    "Government / Political Science": "070",
    "Social and Behavioral Sciences": "080",
    "Component Area Option": "090",
}

# End Constants

def process_string(string):
    global exception_count
    string = string.strip()
    # Match a course: ACCT 2301 Intro Accounting --> ACCT 2301
    match_course = re.findall(r"^([A-Z]+ [0-9][V0-9][0-9]+) ", string)
    if match_course:
        return match_course[0]

    # Match major type title: Major Prep Courses: 24 semester credit hours... --> 24 hours
    match_major_part = re.findall(r"Major(?:.*?)Courses: ([0-9]+)(?:-[0-9]+)* semester credit hours", string)
    if match_major_part:
        return match_major_part[0] + " hours"

    # Check if mapping exists
    mapped = None
    for key in mapping:
        if key.lower() in string.lower():
            mapped = mapping[key]
            break
    if mapped:
        return mapped
    
    for key in starts_map:
        if string.lower().startswith(key.lower()):
            mapped = starts_map[key]
            break
    if mapped:
        return mapped

    # Otherwise, log that the string is not processed, return original
    # print("not processed: " + string)
    exception_count += 1
    if string not in exceptions:
        exceptions[string] = 0
    exceptions[string] += 1
    exception_count += 1
    return string

def process_types(obj):
    global exception_count, exceptions
    if type(obj) is dict:
        # Dictionary type object
        new_tree = {}
        for key in obj:
            new_key = process_types(key)
            new_tree[new_key] = process_types(obj[key])
        return new_tree
    elif type(obj) is list:
        # List type object (mainly select one from n)
        new_tree = []
        for key in obj:
            new_key = process_types(key)
            new_tree.append(new_key)
        return new_tree
    elif type(obj) is str:
        # String type object
        return process_string(obj)
    else:
        # return original since no matches
        # print("not processed obj: " + obj)
        return obj

exception_count = 0
exceptions = {}

def process_major_requirements():
    global exceptions, exception_count
    data = {}
    new_data = {}
    with open('output/major_requirements.txt') as major_req:
        data = json.loads(major_req.read())

    for school_key in data:
        school = data[school_key]
        new_school = {}
        for major_key in school:
            major = school[major_key]
            new_school[major_key] = {}
            for type_key in major:
                new_school[major_key][process_types(type_key)] = process_types(major[type_key])
            # test major has no information
            if new_school[major_key] == {}:
                print(major)
        new_data[school_key] = new_school

    json.dump(new_data, open("output/major_requirements_processed.txt", 'w'))
    print(exception_count)

    # not processed strings
    exceptions = {k: v for k, v in sorted(exceptions.items(), key=lambda item: item[1], reverse=True)}
    print(len(exceptions))


def main(): 
    retrieve_all() # retrieve website, no processing
    process_major_requirements() # process scraped data

    # testing
    # scan_anomaly()
    # fetch_major("https://catalog.utdallas.edu/2020/undergraduate/programs/bbs/neuroscience")


if __name__=="__main__": 
    main() 
import pyparsing as pp
from pyparsing import oneOf
import json
import copy


# https://utdallas.collegescheduler.com/api/terms/2021%20Spring/subjects/CS/courses/3345

with open("data/new.json", 'r') as f:
    data = json.loads(f.read())
with open("data/major_requirements_processed.txt", 'r') as f:
    degreePlan = json.loads(f.read())
with open("data/ecs_computer-science.json", 'r') as f:
    ecs_cs = json.loads(f.read())


def prereq(name):
    if name in data:
        print(name)
        new_data = copy.deepcopy(data[name])
        for index, req in enumerate(new_data["prerequisites"]):
            reqName = req.split(": ")
            new_data["prerequisites"][index] = {reqName[0]: process(reqName[1])}
        return new_data["prerequisites"]
    return {"message": "course name not found."}


def recursion(obj, grade="CR"):
    if type(obj) == list and len(obj) == 1:
        return recursion(obj[0], grade=grade)
    if type(obj) == str:
        if obj in data:
            id = data[obj]["id"]
            return {"type": "course", "id": id, "grade": grade}
        return obj
    
    if type(obj) == list and len(obj) >= 3:
        comparator = obj[1] # and | or
        if comparator != "with":
            obj = {"node": {"type": "op", "id": comparator, "children": [recursion(req, grade=grade) for ind, req in enumerate(obj) if ind % 2 == 0]}}
        else:
            print(obj)
            obj = recursion(obj[0], grade=obj[2][0].split("_")[1])
        return obj
    
    for index, ob in enumerate(obj):
        obj[index] = recursion(ob, grade=grade)
    return obj

def process(st):
    complex_expr = pp.Forward()
    operator = pp.Regex(">=|<=|!=|>|<|=").setName("operator")
    logical = (pp.Keyword("AND") | pp.Keyword("OR") | pp.Keyword("WITH")).setName("logical")
    vars = pp.Regex(r"([A-Z]+ [0-9]+|GRADE_(A|B|C|B-|C-))")
    condition = (vars + operator + vars) | vars
    clause = pp.Group(condition ^ (pp.Suppress(
        "(") + complex_expr + pp.Suppress(")")))

    expr = pp.operatorPrecedence(clause, [
                                ("with", 2, pp.opAssoc.LEFT, ),
                                (oneOf("or and"), 2, pp.opAssoc.LEFT, )
                                 ])

    complex_expr << expr

    try:
        rep = complex_expr.parseString(st)
        strRep = str(rep).replace('\'', '"')
        print(strRep)
        js = json.loads(strRep)
        res = recursion(js)
        return res
    except:
        return {"message": "Prerequisites can not be parsed"}


def start():
    result = {}
    for course in data:
        result[course] = prereq(course)
    with open('prerequisites.json', 'w') as f:
        json.dump(result, f)

def courseToId():
    co_id = {}
    id_co = {}
    for course in data:
        co_id[course] = data[course]["id"]
        id_co[data[course]["id"]] = course
    with open('course_to_id.json', 'w') as f:
        json.dump(co_id, f)
    with open('id_to_course.json', 'w') as f:
        json.dump(id_co, f)

if __name__ == '__main__':
    #start()
    courseToId()

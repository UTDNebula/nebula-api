import json

def cinput(prompt: str, dep: int):
    return input(pad(dep) + prompt)

def pad(n: int):
    return '    ' * n

def getCourse(dep: int):
    name = cinput('Name of course: ', dep)
    return name

def getGroup(dep: int):
    group = {}

    group['name'] = cinput('Name of group: ', dep)
    choose = cinput('Use count(c) or hours(h): ', dep)
    group['choose'] = 'ChooseType.count' if choose == 'c' else 'ChooseType.hours'
    group['pick'] = int(cinput('Pick how many: ', dep))
    group['children'] = []

    numOfChildren = int(cinput('Number of children: ', dep))
    for _ in range(numOfChildren):
        child = getCourseOrGroup(dep + 1, group['name'])
        group['children'].append(child)
    return group

def getCourseOrGroup(dep: int, title: str):
    choose = cinput(f'Add new group(g) or course(c) for {title}: ', dep + 1)
    if choose == 'g':
        return getGroup(dep + 1)
    elif choose == 'c':
        return getCourse(dep + 1)

def getDegreePlan():
    plan = {
        'major': cinput('Major: ', 0),
        'degree': cinput('Degree: ', 0),
        'creditHours': int(cinput('Credit Hours: ', 0))
    }
    plan['children'] = []
    count = int(cinput('Number of children: ', 0))
    for _ in range(count):
        plan['children'].append(getCourseOrGroup(0, plan['major']))
    return plan

plan = getDegreePlan()
str = json.dumps(plan, indent=4)
with open("result.json", 'w') as f:
    f.write(str)
print(str)

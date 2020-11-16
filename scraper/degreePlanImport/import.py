import pandas as pd
import numpy as np
import json

degreePlan = pd.read_excel('20.21 CS Degree Plan1.xlsm', index_col=None, header=None)
df = degreePlan.fillna('N/A')
formatted = df.values.tolist()

# each group = [name row, data start row, data end row]
starts = [[3, 5, 19], [21, 22, 39], [41, 42, 49], [51, 52, 57]]

def lineToDict(row):
    lower_div = row[0]
    upper_div = row[1]
    title = row[3]
    course = row[4]
    return {"title": title, "course": course, "LD": lower_div, "UD": upper_div}

all_info = {}

for start in starts:
    name = formatted[start[0]][0]
    courses = []
    for i in range(start[1], start[2] + 1):
        course = lineToDict(formatted[i])
        courses.append(course)
    all_info[name] = courses

with open('data.json', 'w') as f:
    json.dump(all_info, f)
from firebase_admin import firestore, credentials
import firebase_admin
import json

# credentials
cred = credentials.Certificate("secret.json")
firebase_admin.initialize_app(cred)
db = firestore.client()

# get collection
courses_col = db.collection(u'catalogTest2')

# add data into collection
with open("output/major_requirements_processed.txt") as course_data:
    courses = json.loads(course_data.read())
    for key in courses:
        for major in courses[key]:
            print(key + "-" + major)
            doc = courses_col.document(key + "-" + major)
            doc.set(courses[key][major])
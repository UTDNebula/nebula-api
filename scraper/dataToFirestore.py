import dotenv
import json
import os
import sys

import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore

# import environment variables
dotenv.load_dotenv()

# Use the application default credentials
cred = credentials.Certificate(os.environ('GOOGLE_APPLICATION_CREDENTIALS'))
firebase_admin.initialize_app(cred)

# create database reference
db = firestore.Client()

# import data from scraped json
with open(sys.argv[1], 'r') as file:
    data = json.load(file)

# import data
for download in data['downloads']:

    # skip failed downloads 
    if download == {}:
        continue

    # import sections
    for section in download['report_data']:

        # skip empty sections 
        if section == {}:
            continue 

        # set section data form json 
        section_data = {
            'section_name': section['section_address'],
            'course_prefix': section['course_prefix'],
            'course_number': section['course_number'],
            'section_number': section['section'],
            'class_number': section['class_number'],
            'title': section['title'],
            'topic': section['topic'],
            'instructors': section['instructors'],
            'assistants': section['assistants'],
            'term': section['term'],
            'days': section['days'],
            'times': section['times'],
            'location': section['location'],
            'core_area': section['core_area'],
            'activity_type': section['activity_type'],
            'school': section['school'],
            'department': section['dept']
        }

        # add section to database
        db.collection('sections').document(section_data['section_name']).set(section_data)

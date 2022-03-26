import { readFileSync, writeFileSync } from 'fs';
import { stringify } from 'querystring';

const examReadFile = '../data/ExamsWithoutObjectIds.json';
const examWriteFile = '../data/Exams.json';
const courseFile = '../data/Courses.json';

writeExamsToFile(buildExamCourseAssociations(examReadFile, courseFile), examWriteFile);

function findCourseObjectId(courses: [any], prefix: string, number: string): string {
  for (const course of courses) {
    if (course.subject_prefix === prefix && course.course_number === number) {
      return course._id;
    }
  }

  // course not found
  return null;
}

function buildExamCourseAssociations(examReadFile: string, courseFile: string) {
  //var unknownCourseCount = 0;
  const exams = JSON.parse(readFileSync(examReadFile, 'utf8'));

  const courses = JSON.parse(readFileSync(courseFile, 'utf8'));

  // List of exams
  for (let e = 0; e < exams.length; e++) {
    const exam = exams[e];

    // list of yields
    for (let i = 0; i < exam.yields.length; i++) {
      const choices = exam.yields[i].outcome;

      // list of choices
      for (let j = 0; j < choices.length; j++) {
        const choice = choices[j];

        // list of courses in choice
        for (let k = 0; k < choice.length; k++) {
          const course = choice[k];

          // grab course prefix and number
          if (
            'prefix' in course &&
            'number' in course &&
            course.prefix !== null &&
            course.number !== null
          ) {
            const courseId = findCourseObjectId(courses, course.prefix, course.number.toString());
            if (courseId !== null) {
              // Replace the Course prefix and number object with the associated ObjectId
              exams[e].yields[i].outcome[j][k] = courseId;
            }
            /*else {
                            console.log("COURSE NOT IN COURSES DATABASE");
                            console.log(course.prefix + " " + course.number);
                            unknownCourseCount++;
                        }*/
          }
        }
      }
    }
  }

  //console.log("Number of Unknown Courses: " + unknownCourseCount);

  return exams;
}

function writeExamsToFile(exams: [any], filename = './Exams.json') {
  writeFileSync(filename, JSON.stringify(exams, null, '\t'), { flag: 'w' });
}

writeExamsToFile(buildExamCourseAssociations(examReadFile, courseFile), examWriteFile);

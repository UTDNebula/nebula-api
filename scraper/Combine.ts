import { readFileSync, writeFileSync } from 'fs';
import schemas from '../api/schemas';

// Combine courses together
let Courses = new Map<string, schemas.PsuedoCourse>();
for (let CourseObj of JSON.parse(readFileSync("./data/PsuedoCourses.json", { encoding: "utf-8" })) as schemas.PsuedoCourse[]) {
    if (!Courses.has(CourseObj.internal_course_number))
        Courses.set(CourseObj.internal_course_number, CourseObj);
    else
        Courses.get(CourseObj.internal_course_number).sections.push(...CourseObj.sections);
}

// Combine professors together
let Professors: schemas.Professor[] = JSON.parse(readFileSync("./data/Professors.json", { encoding: 'utf-8' }));
let CombinedProfessors: schemas.Professor[] = [];
for (let Professor of Professors) {
    let ExistingProfessor: schemas.Professor = CombinedProfessors.find((prof: schemas.Professor) => {
        return (prof.first_name == Professor.first_name && prof.last_name == Professor.last_name)
    });
    if (ExistingProfessor)
        ExistingProfessor.sections.push(...Professor.sections);
    else
        CombinedProfessors.push(Professor);
}

writeFileSync("./data/PsuedoCourses.json", JSON.stringify(Array.from(Courses.values()), null, '\t'), { flag: 'w' });
writeFileSync("./data/Professors.json", JSON.stringify(CombinedProfessors, null, '\t'), { flag: 'w' });
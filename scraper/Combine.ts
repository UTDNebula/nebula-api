import { readFileSync, writeFileSync, readdirSync, Dirent } from 'fs';
import { exit } from 'process';
import mongoose from 'mongoose';
import schemas from '../api/schemas';

// Get all semester directories in ./data dir
const SemesterDirectories: Dirent[] = readdirSync("./data", { withFileTypes: true }).filter((entry: Dirent) => {
    return entry.isDirectory();
})

// Get all courses, sections, and profs from data dirs

let Courses: schemas.PsuedoCourse[] = [];
let Sections: schemas.Section[] = [];
let Professors: schemas.Professor[] = [];

for (let directory of SemesterDirectories) {
    Courses.push(...(JSON.parse(readFileSync(`./data/${directory.name}/PsuedoCourses.json`, { encoding: "utf-8" })) as schemas.PsuedoCourse[]));
    Sections.push(...(JSON.parse(readFileSync(`./data/${directory.name}/Sections.json`, { encoding: "utf-8" })) as schemas.Section[]));
    Professors.push(...(JSON.parse(readFileSync(`./data/${directory.name}/Professors.json`, { encoding: "utf-8" })) as schemas.Professor[]));
};

// Combine courses together
let CombinedCourses: Map<string, schemas.PsuedoCourse> | schemas.PsuedoCourse[] = new Map<string, schemas.PsuedoCourse>();
for (let CourseObj of Courses) {
    let ValidSections = [];
    for (let CourseSectionRef of CourseObj.sections) {
        if (Sections.findIndex((section: schemas.Section) => { return section._id == CourseSectionRef }) == -1)
            console.log(`WARN: Course references missing section id ${CourseSectionRef}, removed reference`);
        else
            ValidSections.push(CourseSectionRef);
    }
    CourseObj.sections = ValidSections;
    if (!CombinedCourses.has(CourseObj.internal_course_number))
        CombinedCourses.set(CourseObj.internal_course_number, CourseObj);
    else
        CombinedCourses.get(CourseObj.internal_course_number).sections.push(...CourseObj.sections);
}

// Convert combinedcourses to array
CombinedCourses = Array.from(CombinedCourses.values());

// Combine professors together
let CombinedProfessors: schemas.Professor[] = [];
for (let Professor of Professors) {
    let ValidSections = [];
    for (let ProfSectionRef of Professor.sections) {
        if (Sections.findIndex((section: schemas.Section) => { return section._id == ProfSectionRef }) == -1)
            console.log(`WARN: Professor references missing section id ${ProfSectionRef}, removed reference`);
        else
            ValidSections.push(ProfSectionRef);
    }
    Professor.sections = ValidSections;
    let ExistingProfessor: schemas.Professor = CombinedProfessors.find((prof: schemas.Professor) => {
        return (prof.first_name == Professor.first_name && prof.last_name == Professor.last_name)
    });
    if (ExistingProfessor)
        ExistingProfessor.sections.push(...Professor.sections);
    else
        CombinedProfessors.push(Professor);
}

// Verify validity of section course and prof references
for (let Section of Sections) {
    let ReferencedCourse = Courses.find((course: schemas.PsuedoCourse) => { return course._id == Section.course_reference });
    if (!ReferencedCourse) {
        console.log(`ERROR: Section references missing course id ${Section.course_reference}, aborting`);
        exit(1);
    }
    else
        Section.course_reference = CombinedCourses.find((course: schemas.PsuedoCourse) => { return course.internal_course_number == ReferencedCourse.internal_course_number })._id;
    let NewProfRefs = [];
    for (let ProfessorRef of Section.professors) {
        let ReferencedProfessor = Professors.find((professor: schemas.Professor) => { return professor._id == ProfessorRef });
        if (!ReferencedProfessor) {
            console.log(`ERROR: Section references missing professor id ${ProfessorRef}, aborting`);
            exit(1);
        }
        else
            NewProfRefs.push(CombinedProfessors.find((professor: schemas.Professor) => { return professor.first_name == ReferencedProfessor.first_name && professor.last_name == ReferencedProfessor.last_name})._id);
    }
    Section.professors = NewProfRefs;
}

// Sections are already combined properly and can just be pushed to a file

writeFileSync("./data/CombinedPsuedoCourses.json", JSON.stringify(CombinedCourses, null, '\t'), { flag: 'w' });
writeFileSync("./data/CombinedSections.json", JSON.stringify(Sections, null, '\t'), { flag: 'w' });
writeFileSync("./data/CombinedProfessors.json", JSON.stringify(CombinedProfessors, null, '\t'), { flag: 'w' });

console.log(`Successfully combined data into ${CombinedCourses.length} courses, ${Sections.length} sections, and ${CombinedProfessors.length} professors.`);
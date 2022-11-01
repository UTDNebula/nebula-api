/////////////////////////////////
//	Script for uploading the combined data to combinedDB. Should be run after parsing the combined data data using Parse.ts.
////////////////////////////////

import 'dotenv/config';
import schemas from '../../api/ts/schemas';
import mongoose from 'mongoose';
import { readFileSync } from 'fs';

const Courses: schemas.Course[] = JSON.parse(readFileSync(`./data/CombinedCourses.json`, { encoding: "utf-8" })) as schemas.Course[];
const Sections: schemas.Section[] = JSON.parse(readFileSync(`./data/CombinedSections.json`, { encoding: "utf-8" })) as schemas.Section[];
const Professors: schemas.Professor[] = JSON.parse(readFileSync(`./data/CombinedProfessors.json`, { encoding: "utf-8" })) as schemas.Professor[];

type ModelPair = [mongoose.Model<any>, any];

const CourseModel: mongoose.Model<schemas.Course> = mongoose.model("Courses", schemas.CourseSchema);
const SectionModel: mongoose.Model<schemas.Section> = mongoose.model("Sections", schemas.SectionSchema);
const ProfessorModel: mongoose.Model<schemas.Professor> = mongoose.model("Professors", schemas.ProfessorSchema);

async function main() : Promise<void> {
    await mongoose.connect(`mongodb+srv://${process.env.MONGO_USERNAME}:${process.env.MONGO_PASSWORD}@development-0.gftz1.mongodb.net/combinedDB?retryWrites=true&w=majority`);
    for (let Pair of [[CourseModel, Courses], [SectionModel, Sections], [ProfessorModel, Professors]] as ModelPair[]) {
        await Pair[0].deleteMany();
        await Pair[0].insertMany(Pair[1]);
    };
    console.log("Done");
    await mongoose.disconnect();
}

main();
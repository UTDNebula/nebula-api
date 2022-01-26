/*
This script allows for support of multiple databases using the export schema model,
as each model is assigned to a different database under mongoose.
*/

import dotenv from 'dotenv';
import mongoose from 'mongoose';

import { CourseSchema } from './schemas/course';
import { DegreeSchema } from './schemas/degree';

dotenv.config();

const degreeDB = mongoose.createConnection(process.env.DEGREE_MONGODB_URI);
export const DegreeModel = degreeDB.model('degree', DegreeSchema);

const courseDB = mongoose.createConnection(process.env.COURSE_MONGODB_URI);
export const CourseModel = courseDB.model('course', CourseSchema);
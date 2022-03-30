/* 
commandline params: <rawDataPath> <academic_session>
ex. npx ts-node ./parser/parser.ts /acmutd/utd-grades/master/data/Fall%202019/Fall2019.json 19F
*/

import https from 'https';
import dotenv from 'dotenv';
import fs from 'fs';
import mongoose, { Types } from 'mongoose';

import { SectionModel } from '../api/models/section';

dotenv.config();

const DATA_PATH = process.argv[2];
const SESH = process.argv[3];

enum Grades {
  'A+',
  'A',
  'A-',
  'B+',
  'B',
  'B-',
  'C+',
  'C',
  'C-',
  'D+',
  'D',
  'D-',
  'F',
  'W',
}

type GradeSection = {
  subj: string;
  num: string;
  sect: string;
  CR: number;
  NC: number;
  P: number;
  prof: string;
  grades: object;
  term: string;
};

type ConciseSection = {
  _id: Types.ObjectId;
  section_number: string;
  courses: {
    course_number: string;
    subject_prefix: string;
  };
};

const options = {
  hostname: 'raw.githubusercontent.com',
  port: 443,
  path: DATA_PATH,
  method: 'GET',
};

const req = https.request(options, (res) => {
  if (!SESH) return req.destroy(new Error('Argument for session was not provided.'));
  console.log('Getting data from ' + options.hostname + options.path);
  console.log(`statusCode: ${res.statusCode}`);

  const chunks: Uint8Array[] = [];

  // get all the semester's grade data in chunks and convert
  res.on('data', (chunk) => {
    chunks.push(chunk);
  });
  res.on('end', function () {
    let body: GradeSection[];
    try {
      body = JSON.parse(Buffer.concat(chunks).toString());
    } catch (e) {
      return console.error(e);
    }
    processData(body);
  });
});

req.on('error', (error) => {
  console.error(error);
});

req.end();

const processData = async (data: GradeSection[]) => {
  if (!fs.existsSync('./parser/logs')) fs.mkdirSync('./parser/logs');
  const logger = fs.createWriteStream(__dirname + `/logs/debug.txt`, { flags: 'w' });
  await mongoose.connect(process.env.MONGODB_URI);
  // get all of the specified semester's sections
  const semesterData: ConciseSection[] = await SectionModel.aggregate([
    {
      $match: {
        'academic_session.name': SESH,
      },
    },
    {
      $lookup: {
        from: 'courses',
        localField: 'course_reference',
        foreignField: '_id',
        as: 'courses',
      },
    },
    {
      $unwind: '$courses',
    },
    {
      $project: {
        _id: 1,
        section_number: 1,
        'courses.course_number': 1,
        'courses.subject_prefix': 1,
      },
    },
  ]);
  let count = 0;
  // go through all grade data and find its counterpart in the mongoDB sections
  for (const sect of data) {
    const matchedSection: ConciseSection = semesterData.find(
      (section) =>
        section.courses.course_number == sect.num &&
        section.courses.subject_prefix == sect.subj &&
        section.section_number == sect.sect,
    );
    if (!matchedSection) {
      logger.write(
        `Couldn't find section ${sect.subj}.${sect.num}.${sect.sect} in DB from grade data for ${SESH}.\n`,
      );
      continue;
    }
    // update the grade data for the DB section
    const update = await SectionModel.updateOne(
      { _id: matchedSection._id },
      { grade_distribution: processGrades(sect.grades) },
    );
    if (update.modifiedCount == 0)
      logger.write(
        `${sect.subj}.${sect.num}.${sect.sect} (${matchedSection._id}) was NOT modified in the DB.\n`,
      );
    else
      logger.write(
        `${sect.subj}.${sect.num}.${sect.sect} (${matchedSection._id}) was modified in the DB.\n`,
      );
    updateProgress(count++, semesterData.length);
  }
  await mongoose.disconnect();
  logger.close();
};

// convert grade object to an array for schema
const processGrades = (gradesObject: object) => {
  const gradesArray: number[] = Array(14).fill(0); // A+ A A- B+ B B- C+ C C- D+ D D- F W
  for (const grade in gradesObject) {
    gradesArray[Grades[grade]] = gradesObject[grade];
  }
  return gradesArray;
};

async function updateProgress(sectionNum, maxSections) {
  const i = Math.round((20 * sectionNum) / maxSections);
  const dots = '.'.repeat(i);
  const left = 20 - i;
  const empty = ' '.repeat(left);

  process.stdout.write(`\r[${dots}${empty}] ${i * 5}%`);
}

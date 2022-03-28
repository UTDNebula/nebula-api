import { Schema, connection } from 'mongoose';

export interface Course {
  _id: string;
  course_number: string;
  subject_prefix: string;
  title: string;
  description: string;
  school: string;
  credit_hours: string;
  class_level: string;
  activity_type: string;
  grading: string;
  internal_course_number: string;
  prerequisites: object;
  corequisites: object;
  co_or_pre_requisites: object;
  sections: string[];
  lecture_contact_hours: string;
  laboratory_contact_hours: string;
  offering_frequency: string;
  attributes: object;
}

export const CourseSchema = new Schema<Course>({
  _id: { type: String, required: true },
  course_number: { type: String, required: true },
  subject_prefix: { type: String, required: true },
  title: { type: String, required: true },
  description: { type: String, required: true },
  school: { type: String, required: true },
  credit_hours: { type: String, required: true },
  class_level: { type: String, required: true },
  activity_type: { type: String, required: true },
  grading: { type: String, required: true },
  internal_course_number: { type: String, required: true },
  prerequisites: { type: Object, required: true },
  corequisites: { type: Object, required: true },
  co_or_pre_requisites: { type: Object, required: true },
  sections: { type: [String], required: true },
  lecture_contact_hours: { type: String, required: true },
  laboratory_contact_hours: { type: String, required: true },
  offering_frequency: { type: String, required: true },
  attributes: { type: Object, required: true },
});

const courseDB = connection.useDb('combinedDB');
export const CourseModel = courseDB.model<Course>('course', CourseSchema);

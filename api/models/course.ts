import { Schema, model } from 'mongoose';

export interface Course {
  name: string;
  subject: string;
  number: string;
  description: string;
  credit_hours: number;
  school: string;
  class_level: string;
  activity_type: string;
  grading: string;
  prerequisite_courses: object;
  corequisite_courses: object;
}

const schema = new Schema<Course>({
  name: { type: String, required: true },
  subject: { type: String, required: true },
  number: { type: String, required: true },
  description: { type: String, required: true },
  credit_hours: { type: Number, required: true },
  school: { type: String, required: true },
  class_level: { type: String, required: true },
  activity_type: { type: String, required: true },
  grading: { type: String, required: true },
  prerequisite_courses: { type: Object, required: true },
  corequisite_courses: { type: Object, required: true },
});

export const CourseModel = model<Course>('course', schema);

import { Schema, connection } from 'mongoose';

export type AcademicSession = {
  name: string;
  start_date: string;
  end_date: string;
};

export type Assistant = {
  first_name: string;
  last_name: string;
  role: string;
  email: string;
};

export type Location = {
  building: string;
  room: string;
  map_uri: string;
};

type ModalityType = 'pending' | 'traditional' | 'hybrid' | 'flexible' | 'remote' | 'online';

export type Meeting = {
  start_date: string;
  end_date: string;
  meeting_days: Array<string>;
  start_time: string;
  end_time: string;
  modality: ModalityType;
  location: Location;
};

export interface Section {
  _id: string;
  section_number: string;
  course_reference: Schema.Types.ObjectId;
  section_corequisites: object; // i was too lazy and did not code all the requirements, but it should still work with object
  academic_session: AcademicSession;
  professors: Array<Schema.Types.ObjectId>;
  teaching_assistants: Array<Assistant>;
  internal_class_number: string;
  instruction_mode: string;
  meetings: Array<Meeting>;
  core_flags: Array<string>;
  syllabus_uri: string;
  grade_distribution: Array<number>;
  attributes: object;
}

export const SectionSchema = new Schema<Section>({
  _id: { type: String, required: true },
  section_number: { type: String, required: true },
  course_reference: { type: Schema.Types.ObjectId, required: true },
  section_corequisites: { type: Object, required: true },
  academic_session: { type: Object, required: true },
  professors: { type: [Schema.Types.ObjectId], required: true },
  teaching_assistants: { type: [Object], required: true },
  internal_class_number: { type: String, required: true },
  instruction_mode: { type: String, required: true },
  meetings: { type: [Object], required: true },
  core_flags: { type: [String], required: true },
  syllabus_uri: { type: String, required: true },
  grade_distribution: { type: [Number], required: true },
  attributes: { type: Object, required: true },
});

const sectionDB = connection.useDb('combinedDB');
export const SectionModel = sectionDB.model<Section>('section', SectionSchema);

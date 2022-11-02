import { Schema, connection, Types } from 'mongoose';
import { Location, Meeting } from './section'; // dependency will be resolved with section pull request

export interface Professor {
  _id: Types.ObjectId;
  first_name: string;
  last_name: string;
  titles: Array<string>;
  email: string;
  phone_number: string;
  office: Location;
  profile_uri: string;
  image_uri: string;
  office_hours: Array<Meeting>;
  sections: Array<Schema.Types.ObjectId>;
}

export const ProfessorSchema = new Schema<Professor>({
  _id: { type: Types.ObjectId, required: true },
  first_name: { type: String, required: true },
  last_name: { type: String, required: true },
  titles: { type: [String], required: false },
  email: { type: String, required: false },
  phone_number: { type: String, required: false },
  office: { type: Object, required: false },
  profile_uri: { type: String, required: false },
  image_uri: { type: String, required: false },
  office_hours: { type: [Object], required: false },
  sections: { type: [Types.ObjectId], required: true },
});

const professorDB = connection.useDb('combinedDB');
export const ProfessorModel = professorDB.model<Professor>('professor', ProfessorSchema);

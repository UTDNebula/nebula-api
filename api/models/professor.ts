import { Schema, connection } from 'mongoose';
import { Location, Meeting } from './section'; // dependency will be resolved with section pull request

export interface Professor {
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
  first_name: { type: String, required: true },
  last_name: { type: String, required: true },
  titles: { type: [String], required: true },
  email: { type: String, required: true },
  phone_number: { type: String, required: true },
  office: { type: Object, required: true },
  profile_uri: { type: String, required: true },
  image_uri: { type: String, required: true },
  office_hours: { type: [Object], required: true },
  sections: { type: [Schema.Types.ObjectId], required: true }
});

const professorDB = connection.useDb('professorDB');
export const ProfessorModel = professorDB.model<Professor>('professor', ProfessorSchema);
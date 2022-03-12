import { Schema, connection } from 'mongoose';

type DegreeSubtype = "major" | "minor" | "concentration";

export interface Degree {
  subtype: DegreeSubtype;
  name: string;
  abbreviation: string;
  minimum_credit_hours: number;
  requirements: object;
}

export const DegreeSchema = new Schema<Degree>({
  subtype: { type: String, required: true },
  name: { type: String, required: true },
  abbreviation: { type: String, required: true },
  minimum_credit_hours: { type: Number, required: true },
  requirements: { type: Object, required: true }
});

const degreeDB = connection.useDb('degreeDB');
export const DegreeModel = degreeDB.model<Degree>('degree', DegreeSchema);
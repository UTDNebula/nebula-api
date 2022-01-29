import { Schema, connection } from 'mongoose';

export interface Degree {
  name: string;
  required: number;
  total: number;
  options: Array<object | string>;
}

export const DegreeSchema = new Schema<Degree>({
  name: { type: String, required: true },
  required: { type: Number, required: true },
  total: { type: Number, required: true },
  options: { type: Array, required: true },
});

const degreeDB = connection.useDb('degreeDB');
export const DegreeModel = degreeDB.model<Degree>('degree', DegreeSchema);
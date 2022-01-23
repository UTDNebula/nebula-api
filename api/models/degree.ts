import { Schema, model } from 'mongoose';

export interface Degree {
  name: string;
  required: number;
  total: number;
  options: Array<object | string>;
}

const schema = new Schema<Degree>({
  name: { type: String, required: true },
  required: { type: Number, required: true },
  total: { type: Number, required: true },
  options: { type: Array, required: true },
});

export const DegreeModel = model<Degree>('degree', schema);

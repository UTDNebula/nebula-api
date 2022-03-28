import { Schema, connection } from 'mongoose';

type ExamType = 'AP' | 'ALEKS' | 'CLEP' | 'IB' | 'CS placement';

export interface Exam {
  _id: string;
  type: ExamType;
  yields: Record<number, Schema.Types.ObjectId>;
}

export const ExamSchema = new Schema<Exam>({
  _id: { type: String, required: true },
  type: { type: String, required: true },
  yields: { type: Object, required: true },
});

const examDB = connection.useDb('examDB');
export const ExamModel = examDB.model<Exam>('exam', ExamSchema);

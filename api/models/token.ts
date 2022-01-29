import { Schema, connection } from 'mongoose';

export interface Token {
  creator: string,
  user: string,
  active: boolean
}

const TokenSchema = new Schema<Token>({
  creator: { type: String, required: true },
  user: { type: String, required: true },
  active: { type: Boolean, required: true }
}, {
  timestamps: true
});

const tokenDB = connection.useDb('tokenDB');
export const TokenModel = tokenDB.model<Token>('token', TokenSchema);

import { db, increment, decrement } from '../../../lib/firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from '../auth';
import all_handler from '../../../lib/handler/[id]';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  return all_handler(req, res, "courses", "course");
}

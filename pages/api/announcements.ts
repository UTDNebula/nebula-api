import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from './auth';
import { getByKey, update, remove, post, getAll } from '../../lib/handler/[id]';

const collection = 'announcements';

/**
 * Get all announcements
 * @param req request
 * @param res response
 * @returns JSON of all items in announcements
 */
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  await getAll(collection).then((data) => {
    res.json(data);
  });
  return 0;
}

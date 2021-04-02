import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from './auth';
import { getByKey, update, remove, post, getAll } from '../../lib/handler/[id]';

const collection = 'announcements';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  await getAll(collection).then((data) => {
    res.json(data);
  });
  return 0;
}

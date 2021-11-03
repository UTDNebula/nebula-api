import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from '../auth';
import { getByKey, update, remove, post } from '../../../lib/handler/[id]';

const collection = 'announcements';

/**
 * Handles all routes for announcements
 * @param req request
 * @param res response
 * @returns appropriate JSON response
 */
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const id = parseInt(req.query.id as string);
  if (req.method == 'GET')
    await getByKey(collection, 'id', id).then((data) => {
      // process courses
      res.json(data);
    });
  await authCheck(req)
    .then(async (_) => {
      if (req.method === 'PUT') {
        await update(collection, 'id', id, JSON.parse(req.body)).then((resp) => {
          res.json(resp);
        });
      } else if (req.method === 'DELETE') {
        await remove(collection, 'id', id).then((resp) => {
          res.json(resp);
        });
      } else if (req.method === 'POST') {
        await post(collection, 'id', id, JSON.parse(req.body), 'course').then((resp) => {
          res.json(resp);
        });
      }
    })
    .catch((_) => {
      res.json({ auth: false });
    });
  return 0;
}

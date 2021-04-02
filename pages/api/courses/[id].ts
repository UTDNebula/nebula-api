import { db, increment, decrement } from '../../../lib/firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from '../auth';
import { getAll, update, remove, post } from '../../../lib/handler/[id]';
import { courseType } from '../../../lib/types/types';

const collection = 'courses';

export function process_course(course: any) : courseType {
  return {
    id: course.id,
    course: course.course,
    description: course.description,
    title: course.titleLong,
    prerequisites: course.prerequisites,
    corequisites: course.corequisites,
    hours: course.hours,
    inclass: course.inclass,
    outclass: course.outclass
  }
}

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const id: number = parseInt(req.query.id as string);
  if (req.method == 'GET')
    await getAll(collection, 'id', id).then((data) => {
      // process courses
      res.json(process_course(data));
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

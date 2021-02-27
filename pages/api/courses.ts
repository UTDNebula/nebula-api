import { db, increment } from '../../lib/firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';

// TODO: fallback in case of empty ID, should find better way
export default async (req: NextApiRequest, res: NextApiResponse) => {
  if (req.method === 'POST') {
    const course = JSON.parse(req.body);

    // get counter to set new id
    const courseCol = db.collection('courses');
    const counter = courseCol.doc('_counter');
    const count = await counter.get();
    const newId = count.data()['count'];
    course['id'] = parseInt(newId);

    await courseCol.doc(course['course']).set(course);
    await counter.update({ count: increment });

    res.json({ message: 'Course updated and counter updated' });
  }

  return 0;
};

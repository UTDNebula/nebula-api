import { db } from '../../../lib/firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';

export default async (req: NextApiRequest, res: NextApiResponse) => {
  const name = (req.query.name as string).toUpperCase();
  const end = name.replace(/.$/, (c) => String.fromCharCode(c.charCodeAt(0) + 1));
  const courseCol = db.collection('courses');
  const snapshot = await courseCol.where('course', '>=', name).where('course', '<', end).get();
  if (snapshot.empty) {
    res.json([]);
  } else {
    const result = [];
    snapshot.forEach((doc) => {
      result.push(doc.data());
    });
    res.json(result);
  }
};

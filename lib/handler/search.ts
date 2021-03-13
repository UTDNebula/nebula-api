import { db } from '../firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function search_all (req: NextApiRequest, res: NextApiResponse, collection: string, name: string, end: string, titleName: string) {
  const courseCol = db.collection(collection);

  const snapshot = await courseCol.where(titleName, '>=', name).where(titleName, '<', end).get();
  if (snapshot.empty) {
    res.json([]);
  } else {
    const result = [];
    snapshot.forEach((doc) => {
      result.push(doc.data());
    });
    console.log(result);
    res.json(result);
  }
};

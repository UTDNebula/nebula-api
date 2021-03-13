import { db, increment, decrement } from '../firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from './auth';

export default async function all_handler(req: NextApiRequest, res: NextApiResponse, collection: string, titleName: string) {
  const id = parseInt(req.query.id as string);
  if (req.method === 'GET') {
    await db
      .collection(collection)
      .where('id', '==', id)
      .get()
      .then((querysnapshot) => {
        res.json(querysnapshot.docs[0].data());
      })
      .catch((error) => {
        res.json({ error });
      });
  } else if (req.method === 'PUT') {
    const authorized = await authCheck(req);
    if (!authorized) {
      res.json({ message: 'not authorized.' });
      return;
    }
    const course = JSON.parse(req.body);
    await db
      .collection(collection)
      .where('id', '==', id)
      .get()
      .then((snapshot) => {
        if (snapshot.empty) {
          res.json({ updated: false });
        } else {
          snapshot.docs[0].ref.update(course);
          res.json({ updated: true });
        }
      });
  } else if (req.method === 'DELETE') {
    const authorized = await authCheck(req);
    if (!authorized) {
      res.json({ message: 'not authorized.' });
      return;
    }
    const result = db.collection(collection).where('id', '==', id);
    await result.get().then(async (snapshot) => {
      if (snapshot.empty) {
        res.json({ deleted: false });
      } else {
        snapshot.docs[0].ref.delete();
        const courseCol = db.collection(collection);
        const counter = courseCol.doc('_counter');
        await counter.update({ count: decrement });
        res.json({ deleted: true });
      }
    });
  } else if (req.method === 'POST') {
    const authorized = await authCheck(req);
    if (!authorized) {
      res.json({ message: 'not authorized.' });
      return;
    }

    const course = JSON.parse(req.body);
    // get counter to set new id
    const courseCol = db.collection(collection);
    const counter = courseCol.doc('_counter');
    const count = await counter.get();
    const newId = count.data()['count'];
    course['id'] = parseInt(newId);
    console.log(course);

    await courseCol.doc(course[titleName]).set(course);
    await counter.update({ count: increment });

    res.json({ message: 'Data updated and counter updated' });
  }
  return 0;
}

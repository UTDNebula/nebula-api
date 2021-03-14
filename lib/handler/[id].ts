import { db, increment, decrement } from '../firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from './auth';

export function getAll(collection: string, key: string, value: number) {
  return db
    .collection(collection)
    .where(key, '==', value)
    .get()
    .then((querysnapshot) => {
      return querysnapshot.docs[0].data();
    });
}

export function update(collection: string, key: string, value: number, updated: any) {
  return db.collection(collection)
    .where(key, '==', value)
    .get()
    .then((snapshot) => {
      if (snapshot.empty) {
        return { updated: false };
      } else {
        snapshot.docs[0].ref.update(updated);
        return { updated: true };
      }
    });
}

export function remove(collection: string, key: string, value: number) {
  const result = db.collection(collection).where(key, '==', value);
  return result.get().then(async (snapshot) => {
    if (snapshot.empty) {
      return { deleted: false };
    } else {
      snapshot.docs[0].ref.delete();
      const col = db.collection(collection);
      const counter = col.doc('_counter');
      await counter.update({ count: decrement });
      return { deleted: true };
    }
  });
}

export function post(
  collection: string,
  key: string,
  value: number,
  data: any,
  titleName: string,
) {
  // get counter to set new id
  const col = db.collection(collection);
  const counter = col.doc('_counter');
  return counter.get().then(async (count) => {
    const newId = count.data()['count'];
    data['id'] = parseInt(newId);

    await col.doc(data[titleName]).set(data);
    await counter.update({ count: increment });

    return { added: true };
  });
}

// export default async function all_handler(
//   req: NextApiRequest,
//   res: NextApiResponse,
//   collection: string,
//   titleName: string,
// ) {
//   const id = parseInt(req.query.id as string);
//   if (req.method === 'GET') {
//     return getAll(collection, 'id', id);
//   }
//   return authCheck(req)
//     .then(async (x) => {
//       if (req.method === 'PUT') {
//         return update(collection, 'id', id, JSON.parse(req.body));
//       } else if (req.method === 'DELETE') {
//         return remove(collection, 'id', id);
//       } else if (req.method === 'POST') {
//         return post(collection, 'id', id, JSON.parse(req.body), titleName);
//       }
//     })
//     .catch((_) => {
//       return { auth: false, updated: false };
//     });
// }

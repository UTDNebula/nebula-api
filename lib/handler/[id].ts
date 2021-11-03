import { db, increment, decrement } from '../firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from './auth';

/**
 * Get all documents from a collection
 * @param collection collection name
 * @returns all documents in that collection
 */
export function getAll(collection: string) {
  return db
    .collection(collection)
    .get()
    .then((snapshot) => {
      if (snapshot.empty) {
        return [];
      } else {
        const result = [];
        snapshot.forEach((doc) => {
          result.push(doc.data());
        });
        return result;
      }
    });
}

/**
 * Find first snapshot in collection where key = value
 * @param collection collection name
 * @param key name of key to match
 * @param value value corresponding to the key
 * @returns
 */
export function getByKey(collection: string, key: string, value: number) {
  return db
    .collection(collection)
    .where(key, '==', value)
    .get()
    .then((querysnapshot) => {
      return querysnapshot.docs[0].data();
    });
}

/**
 * Update collection entry
 * @param collection collection name
 * @param key key name
 * @param value value corresponding to the key
 * @param updated updated data
 * @returns
 */
export function update(collection: string, key: string, value: number, updated: any) {
  return db
    .collection(collection)
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

/**
 * Remove entry from collection
 * @param collection collection name
 * @param key key name
 * @param value matching value
 * @returns if delete was successful
 */
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

/**
 * Add new entry to collection
 * @param collection collection name
 * @param key key name
 * @param value value for the key
 * @param data data for the new entry
 * @param titleName title of the document
 * @returns post status
 */
export function post(collection: string, key: string, value: number, data: any, titleName: string) {
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

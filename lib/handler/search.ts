import { db } from '../firebaseAdmin';

/**
 * Find all matching results from collection
 * @param collection collection name
 * @param name name range start
 * @param end name range end
 * @param titleName key to match against
 * @returns
 */
export default async function search_all(
  collection: string,
  name: string,
  end: string,
  titleName: string,
) {
  const courseCol = db.collection(collection);
  return courseCol
    .where(titleName, '>=', name)
    .where(titleName, '<', end)
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

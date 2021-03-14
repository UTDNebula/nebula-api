import { db } from '../firebaseAdmin';

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

import { db, decrement } from '../../../lib/firebase';

export default async function handler(req, res) {
    const id = parseInt(req.query.id);
    if (req.method === 'GET') {
        db.collection('courses')
            .where('id', '==', id)
            .get()
            .then((querySnapshot) => {
                res.json(querySnapshot.docs[0].data());
            })
            .catch((error) => {
                res.json({ error });
            });
    } else if (req.method === 'PUT') {
        const course = req.body;
        db.collection('courses')
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
        console.log(`Deleting course with id ${id}`);
        const result = db.collection('courses').where('id', '==', id);
        result.get().then(async (snapshot) => {
            if (snapshot.empty) {
                console.log('not found');
                res.json({ deleted: false });
            } else {
                snapshot.docs[0].ref.delete();
                const counter = courseCol.doc('_counter');
                await counter.update({ count: decrement });
                res.json({ deleted: true });
            }
        });
    }
}

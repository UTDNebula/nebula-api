import { db, increment, decrement } from "../../../lib/firebaseAdmin"
import { NextApiRequest, NextApiResponse } from 'next';
import { authCheck } from '../auth';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
    const id = parseInt(req.query.id as string);
    if (req.method === 'GET') {
        await db.collection('courses')
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
            res.json({ "message": "not authorized." });
            return;
        }
        const course = JSON.parse(req.body);
        console.log(course);
        await db.collection('courses')
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
            res.json({ "message": "not authorized." });
            return;
        }
        console.log(`deleting course with id ${id}`);
        const result = db.collection('courses').where('id', '==', id);
        await result.get().then(async (snapshot) => {
            if (snapshot.empty) {
                console.log('not found');
                res.json({ deleted: false });
            } else {
                snapshot.docs[0].ref.delete();
                const courseCol = db.collection('courses');
                const counter = courseCol.doc('_counter');
                await counter.update({ count: decrement });
                res.json({ deleted: true });
            }
        });
    } else if (req.method === 'POST') {
        const authorized = await authCheck(req);
        if (!authorized) {
            res.json({ "message": "not authorized." });
            return;
        }

        const course = JSON.parse(req.body);
        console.log(course);
        console.log(req.body);

        // get counter to set new id
        const courseCol = db.collection('courses');
        const counter = courseCol.doc('_counter');
        const count = await counter.get();
        console.log(count.data());
        const newId = count.data()['count'];
        course['id'] = parseInt(newId);

        await courseCol.doc(course['course']).set(course);
        await counter.update({ count: increment });

        res.json({ message: 'Course updated and counter updated' });
    }
    return 0;
}

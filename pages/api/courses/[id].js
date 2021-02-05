import { db, decrement } from '../../../lib/firebaseAdmin';
import { parseCookies } from 'nookies';
import { auth } from '../../../lib/firebaseAdmin';

const authCheck = async (req) => {
    const cookies = parseCookies({ req });
    if (!cookies) return false;
    try {
        const token = await auth.verifyIdToken(cookies.token);
        const { uid, email } = token;
        console.log(email);
        return true;
    } catch (err) {
        // not logged in
        console.log("no login")
        return false;
    }
}

export default async function handler(req, res) {
    const id = parseInt(req.query.id);
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
                const counter = courseCol.doc('_counter');
                await counter.update({ count: decrement });
                res.json({ deleted: true });
            }
        });
    }
    return 0;
}

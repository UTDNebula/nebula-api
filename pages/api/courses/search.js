import { db } from '../../../lib/firebaseAdmin';

export default async (req, res) => {
    const name = req.query.name.toUpperCase();
    console.log(name);

    const end = name.replace(/.$/, (c) =>
        String.fromCharCode(c.charCodeAt(0) + 1)
    );
    const courseCol = db.collection('courses');
    const snapshot = await courseCol
        .where('course', '>=', name)
        .where('course', '<', end)
        .get();
    if (snapshot.empty) {
        console.log('not found');
        res.json([]);
    } else {
        console.log('found!');
        const result = [];
        snapshot.forEach((doc) => {
            result.push(doc.data());
        });
        res.json(result);
    }
};

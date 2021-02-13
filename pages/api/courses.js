import { db, increment } from '../../lib/firebaseAdmin';

export default async (req, res) => {
    console.log(req.method);
    if (req.method === 'POST') {
        const course = JSON.parse(req.body);
        console.log(course);

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
};

import express from 'express';
import db from '../firestore';

const router = express.Router();

router.route('/').get((req, res) => {
  res.json({ success: true });
});

router.route('/v1/sections/:id').get(async (req, res) => {
  const section = await db.collection('sections').doc(req.params.id).get();
  if (section.exists) {
    res.status(200).json(section.data());
  } else {
    res.status(404).json({ message: 'Section not found by id.' });
  }
});

export default router;

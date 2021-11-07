import express from 'express';
import db from '../firestore';

const router = express.Router();

router.route('/').get((req, res) => {
  res.json({ success: true });
});

export default router;

import express from 'express';

const router = express.Router();

router.route('/').get((req, res) => {
  res.json({ success: true });
});

export default router;

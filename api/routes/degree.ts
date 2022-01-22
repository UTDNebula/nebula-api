import express from 'express';

import { degreeSearch, degreeById } from '../controllers/degree';

const router = express.Router();

router.get('/search?', degreeSearch);
router.get('/:id', degreeById);

export default router;

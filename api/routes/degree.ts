import express from 'express';

import { degreeSearch, degreeById, degreeInsert } from '../controllers/degree';

const router = express.Router();

router.get('/search?', degreeSearch);
router.get('/:id', degreeById);
router.post('/insert', degreeInsert);

export default router;

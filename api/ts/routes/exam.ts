import express from 'express';

import { examSearch, examById } from '../controllers/exam';

const router = express.Router();

router.get('/', examSearch);
router.get('/:id', examById);

export default router;

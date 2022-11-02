import express from 'express';

import { professorSearch, professorById } from '../controllers/professor';

const router = express.Router();

router.get('/', professorSearch);
router.get('/:id', professorById);

export default router;

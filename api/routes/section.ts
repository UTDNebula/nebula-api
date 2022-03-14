import express from 'express';

import { sectionSearch, sectionById } from '../controllers/section';

const router = express.Router();

router.get('/', sectionSearch);
router.get('/:id', sectionById);

export default router;

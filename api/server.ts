import express from 'express';
import dotenv from 'dotenv';
import middlewareController from './middleware/controller';

import * as router from './routes/router';
import * as sections from './routes/sections';

const app = express();

dotenv.config();

app.use(middlewareController);
app.use('/', router.default);
app.use('/v1/sections', sections.default);

app.listen(process.env.PORT, () => {
  console.log(`The server has started on port ${process.env.PORT}`);
});

import express from 'express';
import dotenv from 'dotenv';
import middlewareController from './middleware/controller';

import router from './routes/router';
import sections from './routes/sections';

const app = express();
const DEFAULT_PORT = 3000;

dotenv.config();

app.use(middlewareController);
app.use('/', router);
app.use('/v1/sections', sections);

app.listen(DEFAULT_PORT, () => {
  console.log(`The server has started on port ${DEFAULT_PORT}`);
});

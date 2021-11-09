import express from 'express';
import dotenv from 'dotenv';

import * as router from './routes/router';

const app = express();

dotenv.config();

app.use('/', router.default);

app.listen(process.env.PORT, () => {
  console.log(`The server has started on port ${process.env.PORT}`);
});

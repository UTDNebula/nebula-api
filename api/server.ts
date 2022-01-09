import dotenv from 'dotenv';
import express from 'express';
import mongoose from 'mongoose';

import router from './routes/router';

dotenv.config();

mongoose.connect(process.env.MONGODB_URI).then(() => {
  const app = express();

  app.use('/', router);

  app.listen(process.env.PORT, () => {
    console.log(`The server has started on port ${process.env.PORT}`);
  });
});

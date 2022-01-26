import dotenv from 'dotenv';
import express from 'express';
import mongoose from 'mongoose';

import course from './routes/course';
import degree from './routes/degree';

dotenv.config();

mongoose.connect(process.env.MONGODB_URI).then(() => {
  const app = express();

  app.use('/course', course);
  app.use('/degree', degree);

  app.listen(process.env.PORT, () => {
    console.log(`The server has started on port ${process.env.PORT}`);
  });
});
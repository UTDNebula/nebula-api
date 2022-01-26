import dotenv from 'dotenv';
import express from 'express';

import course from './routes/course';
import degree from './routes/degree';

dotenv.config();

const app = express();

app.use('/course', course);
app.use('/degree', degree);

app.listen(process.env.PORT, () => {
  console.log(`The server has started on port ${process.env.PORT}`);
});

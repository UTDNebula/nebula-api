import all_search from '../../../lib/handler/search';
import { NextApiRequest, NextApiResponse } from 'next';
import { process_course } from './[id]';

/**
 * Search courses
 */
export default async (req: NextApiRequest, res: NextApiResponse) => {
  const name = req.query.name as string;
  const end = name.replace(/.$/, (c) => String.fromCharCode(c.charCodeAt(0) + 1));
  await all_search('courses', name, end, 'course').then((resp) => {
    let processed_data = [];
    resp.forEach((data) => {
      processed_data.push(process_course(data));
    });
    res.json(processed_data);
  });
  return 0;
};

import all_search from '../../../lib/handler/search';
import { NextApiRequest, NextApiResponse } from 'next';

export default async (req: NextApiRequest, res: NextApiResponse) => {
  const name = (req.query.name as string);
  const end = name.replace(/.$/, (c) => String.fromCharCode(c.charCodeAt(0) + 1));
  await all_search("announcements", name, end, "announcement").then(resp => {
    res.json(resp);
  });
  return 0;
};

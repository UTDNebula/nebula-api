import all_search from '../../../lib/handler/search';
import { NextApiRequest, NextApiResponse } from 'next';

export default async (req: NextApiRequest, res: NextApiResponse) => {
  const name = (req.query.name as string);
  const end = name.replace(/.$/, (c) => String.fromCharCode(c.charCodeAt(0) + 1));
  return all_search(req, res, "announcements", name, end, "title");
};

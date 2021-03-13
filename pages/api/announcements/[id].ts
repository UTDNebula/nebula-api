import all_handler from '../../../lib/handler/[id]';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  return all_handler(req, res, "announcements", "title");
}

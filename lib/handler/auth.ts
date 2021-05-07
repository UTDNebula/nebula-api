import { parseCookies } from 'nookies';
import { auth } from '../firebaseAdmin';
import { NextApiRequest } from 'next';

/**
 * Check if authorized using cookies
 * @param req request
 * @returns resolve/reject appropriately
 */
export const authCheck = async (req: NextApiRequest) => {
  const cookies = parseCookies({ req });
  if (!cookies) return new Promise((res, rej) => rej);
  try {
    const token = await auth.verifyIdToken(cookies.token);
    const { uid, email } = token;
    // logged in
    return new Promise((res, rej) => res);
  } catch (err) {
    // not logged in
    return new Promise((res, rej) => rej);
  }
};

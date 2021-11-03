import { parseCookies } from 'nookies';
import { auth } from '../../lib/firebaseAdmin';
import { NextApiRequest } from 'next';

/**
 * @param req request
 * @returns if user is authorized
 */
export const authCheck = async (req: NextApiRequest) => {
  const cookies = parseCookies({ req });
  if (!cookies) return false;
  try {
    const token = await auth.verifyIdToken(cookies.token);
    const { uid, email } = token;
    return true;
  } catch (err) {
    // not logged in
    return false;
  }
};

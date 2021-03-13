import { parseCookies } from 'nookies';
import { auth } from '../firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';

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

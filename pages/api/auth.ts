import { parseCookies } from 'nookies';
import { auth } from '../../lib/firebaseAdmin';
import { NextApiRequest, NextApiResponse } from 'next';

export const authCheck = async (req: NextApiRequest) => {
    const cookies = parseCookies({ req });
    if (!cookies) return false;
    try {
        const token = await auth.verifyIdToken(cookies.token);
        const { uid, email } = token;
        console.log(email);
        return true;
    } catch (err) {
        // not logged in
        console.log("no login")
        return false;
    }
}
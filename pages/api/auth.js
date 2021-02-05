import { parseCookies } from 'nookies';
import { auth } from '../../lib/firebaseAdmin';

const authCheck = async (req) => {
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

export default async (req, res) => {
    const authorized = await authCheck(req);
    if (!authorized) {
        res.json({ "authorized": false });
    } else {
        res.json({ "authorized": true })
    }
    return;
};

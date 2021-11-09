import db from '../database/auth';
import bcrypt from 'bcryptjs';

// ----------
// MIDDLEWARE
// ----------

async function middlewareController (req, res, next) {
    const apiHash = req.header('Authorization');
    var verified: boolean = false;
  
    // stop the program before we get a read if there is no key
    if (!apiHash) {
      res.status(403).json({ message: 'API key was not provided.' });
      return;
    }

    // check each hash for a match
    const hashes = await db.collection('api_hashes').get();
    hashes.forEach(hash => {
      if (bcrypt.compareSync(apiHash.toString(), hash.id)) {
        verified = true;
      }
    });
  
    if (verified) {
      next();
    } else {
      res.status(403).json({ message: 'Could not authenticate api key.' });
    }
};

export default middlewareController;
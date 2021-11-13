import db from '../database/auth';
import bcrypt from 'bcryptjs';

// ----------
// MIDDLEWARE
// ----------

async function middlewareController (req, res, next) {
  const apiHash = req.header('Authorization');
  const API_HASHES_COLLECTION = 'api_hashes';
  
  // Decline request if it is not authenticated
  if (!apiHash) {
    res.status(403).json({ message: 'API key was not provided.' });
    return;
  }

  // Ensure API keys are valid
  try {
    const hashes = await db.collection(API_HASHES_COLLECTION).get();
    for (const hash of hashes.docs) {
      if (bcrypt.compareSync(apiHash.toString(), hash)) {
        next();
        return;
      }
    }
    res.status(403).json({ message: 'Could not authenticate API key.' });
  } catch (error) {
    res.status(500).json({ message: 'Internal server error while processing API key.'});
    console.error(error);
    return;
  }
};

export default middlewareController;
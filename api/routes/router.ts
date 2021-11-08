import express from 'express';
import db from '../firestore';
import bcrypt from 'bcryptjs';

const router = express.Router();

router.route('/').get((req, res) => {
  res.json({ success: true });
});

router.use(async (req, res, next) => {
  const apiHash = req.header('Authorization');
  var verified: boolean = false;
  if (!apiHash) {
    res.status(403).json({ message: 'API key was not provided.' });
    return;
  }
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
});

router.route('/v1/sections/search/').get(async (req, res, next) => {
  var query = db.collection('sections');

  // go through each http requested query and add them to the firestore query
  Object.keys(req.query).forEach(element => {
    var filter: string = req.query[element].toString();
    if (element == 'times') { // converted from XX:XX_XX:XX to XX:XX - XX:XX
      let times = filter.split('_');
      filter = times[0] + ' - ' + times[1];
    } else if (element == 'days') { // converted from Monday_Wednesday to Monday & Wednesday
      let days = filter.split('_');
      filter = '';
      days.forEach((day, i) => {
        filter += day;
        if (i != days.length - 1) filter += ' & ';
      });
    }
    query = query.where(element, '==', filter);
  });

  // get an array of matching sections
  query = await query.get();
  var sections: object[] = [];
  if(query.empty) {
    res.status(404).json({ message: 'Section(s) not found with these query parameters.' });
    return;
  }

  // create an array of section data and send
  query.forEach(section => {
    sections.push(section.data());
  });
  res.status(200).json(sections);
});

router.route('/v1/sections/:id').get(async (req, res, next) => {
  const section = await db.collection('sections').doc(req.params.id).get();
  if (section.exists) {
    res.status(200).json(section.data());
  } else {
    res.status(404).json({ message: 'Section not found by id.' });
  }
});

export default router;

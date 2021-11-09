import express from 'express';
import db from '../database/auth';

const router = express.Router();

// --------
// SECTIONS
// --------

router.route('/search').get(async (req, res, next) => {

  // go through each http requested query and add them to the firestore query
  var query = db.collection('sections');
  Object.keys(req.query).forEach(element => {
    var filter: string = req.query[element].toString();

    // converted from XX:XX_XX:XX to XX:XX - XX:XX
    if (element == 'times') { 
      let times = filter.split('_');
      filter = times[0] + ' - ' + times[1];

    // converted from Monday_Wednesday to Monday & Wednesday
    } else if (element == 'days') { 
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

router.route('/:id').get(async (req, res, next) => {
  const section = await db.collection('sections').doc(req.params.id).get();
  if (section.exists) {
    res.status(200).json(section.data());
  } else {
    res.status(404).json({ message: 'Section not found by id.' });
  }
});

export default router;

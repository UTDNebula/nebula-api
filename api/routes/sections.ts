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
    query = query.where(element, '==', req.query[element].toString());
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

require("dotenv").config();

const { firestore, db } = require("../../lib/firebase");
const {
  _getAll,
  _post,
  _findExact,
  _findFuzzy,
  _deleteById,
  _patch,
} = require("./crudHandler.js");

// DEBUG init database data
// const course_info = require('./data/21s_courses.json');
const debug = false;

async function course_getAll(req, res) {
  const result = await _getAll(req.app.get("db").collection("courses"));
  res.json(result);
}

async function course_post(req, res) {
  await _post(
    req.app.get("db").collection("courses"),
    req.body,
    req.app.get("increment"),
    "course"
  );
  res.json({ message: "Course updated and counter updated" });
}

async function course_findById(id) {
  return _findExact(db.collection("courses"), "id", parseInt(id));
}

async function course_findByName(req, res) {
  const name = req.params.name.toUpperCase();
  const result = await _findFuzzy(
    req.app.get("db").collection("courses"),
    "course",
    name
  );
  res.json(result);
}

async function course_deleteById(req, res) {
  const result = await _deleteById(
    req.app.get("db").collection("courses"),
    "id",
    parseInt(req.params.id),
    req.app.get("decrement")
  );
  res.json({ deleted: result });
}

async function course_patch(req, res) {
  const result = await _patch(
    req.app.get("db").collection("courses"),
    parseInt(req.params.id),
    req.body
  );
  return res.json({ updated: result });
}

async function course_init(req, res) {
  if (!debug) res.json({ message: "action not allowed." });
  let batch = req.app.get("db").batch();
  course_info.forEach((course, index) => {
    let ref = req.app.get("db").collection("courses").doc(course["course"]);
    batch.set(ref, course);
    // Firestore only allows 500 docs every batch write
    if (index % 400 == 0) {
      batch.commit();
      batch = req.app.get("db").batch();
      console.log(`batch ${index / 400} saved.`);
    }
  });
  batch.commit();
  res.json({ message: "data saved successfully!" });
}

module.exports = {
  course_getAll,
  course_deleteById,
  course_findById,
  course_findByName,
  course_patch,
  course_post,
  course_init,
};

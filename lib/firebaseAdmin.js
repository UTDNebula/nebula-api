import admin from "firebase-admin";

if (!admin.apps.length) {
  admin.initializeApp({
    credential: admin.credential.cert({
      project_id: process.env.firestore_project_id,
      private_key: process.env.firestore_private_key.replace(/\\n/g, "\n"),
      client_email: process.env.firestore_client_email,
    })
  });
}

const db = admin.firestore();
const auth = admin.auth();
const increment = admin.firestore.FieldValue.increment(1);
const decrement = admin.firestore.FieldValue.increment(-1);

module.exports = {
    db, auth, increment, decrement
}

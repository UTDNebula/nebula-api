import admin from "firebase-admin";

require("dotenv").config();

const firebase = require('firebase/app');
require('firebase/auth');

const firebaseConfig = {
  apiKey: process.env.FIREBASE_API_KEY,
  authDomain: 'cometplanning.firebaseapp.com',
  databaseURL: 'https://cometplanning.firebaseio.com',
  projectId: 'cometplanning',
  storageBucket: 'cometplanning.appspot.com',
  messagingSenderId: process.env.FIREBASE_MESSAGING_ID,
  appId: process.env.FIREBASE_APP_ID
};

if (!firebase.apps.length) {
    firebase.initializeApp(firebaseConfig);
}

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
const increment = admin.firestore.FieldValue.increment(1);
const decrement = admin.firestore.FieldValue.increment(-1);

module.exports = {
    firebase, db, increment, decrement
}

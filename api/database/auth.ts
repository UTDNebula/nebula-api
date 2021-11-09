const { initializeApp, applicationDefault, cert } = require('firebase-admin/app');
const { getFirestore, Timestamp, FieldValue } = require('firebase-admin/firestore');

var serviceAccount: object;

try {
  serviceAccount = require('../../api/database/serviceAccountKey.json');
  console.log('Starting application using Service Account.');
} catch (error) {
  console.log('Starting application using GCP.');
  serviceAccount = null;
}

initializeApp({
  credential: serviceAccount ? cert(serviceAccount) : applicationDefault()
});

const db = getFirestore();

export default db;
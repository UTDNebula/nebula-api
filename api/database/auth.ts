const { initializeApp, applicationDefault, cert } = require('firebase-admin/app');
const { getFirestore, Timestamp, FieldValue } = require('firebase-admin/firestore');

var serviceAccount: object;

try {
  serviceAccount = require('../../serviceAccountKey.json');
} catch (error) {
  serviceAccount = null;
}

initializeApp({
  credential: serviceAccount ? cert(serviceAccount) : applicationDefault()
});

const db = getFirestore();

export default db;
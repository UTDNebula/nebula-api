import admin from 'firebase-admin';

if (!admin.apps.length) {
  admin.initializeApp({
    credential: admin.credential.cert({
      projectId: process.env.firestore_project_id,
      privateKey: process.env.firestore_private_key.replace(/\\n/g, '\n'),
      clientEmail: process.env.firestore_client_email,
    }),
  });
}

export const db = admin.firestore();
export const auth = admin.auth();
export const increment = admin.firestore.FieldValue.increment(1);
export const decrement = admin.firestore.FieldValue.increment(-1);

import firebase from 'firebase/app';
import 'firebase/auth';

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_firebase_apiKey,
  authDomain: 'cometplanning.firebaseapp.com',
  databaseURL: 'https://cometplanning.firebaseio.com',
  projectId: 'cometplanning',
  storageBucket: 'cometplanning.appspot.com',
  messagingSenderId: process.env.NEXT_PUBLIC_firebase_messagingSenderId,
  appId: process.env.NEXT_PUBLIC_firebase_appId,
};

if (!firebase.apps.length) {
  firebase.initializeApp(firebaseConfig);
}

export default firebase;

import firebase from '../lib/firebase';
import { useEffect } from 'react';
import { useAuth } from './use-auth';
import { useRouter } from 'next/router';

const Home: React.FunctionComponent = () => {
  const auth = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (auth.user) {
      router.push('/console');
    }
  });

  return (
    <div className="min-h-screen min-w-screen flex">
      <button
        className="button py-2 px-4 rounded-md shadow-md bg-blue-600 hover:bg-blue-700 mx-auto place-self-center text-xl text-white"
        onClick={() => {
          var provider = new firebase.auth.GoogleAuthProvider();
          firebase
            .auth()
            .signInWithPopup(provider)
            .then((result) => {
              router.push('/console');
            })
            .catch((error) => {
              console.error(error);
            });
        }}
      >
        Log In
      </button>
    </div>
  );
};

export default Home;

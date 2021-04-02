import firebase from '../lib/firebase';
import { useEffect } from 'react';
import { useAuth } from '../components/use-auth';
import { useRouter } from 'next/router';

const Home: React.FunctionComponent = () => {
  const auth = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (auth.user) {
      router.push('/console/courses');
    }
  });

  return (
    <div className="min-h-screen min-w-screen flex">
      <button
        className="btn1 mx-auto place-self-center"
        onClick={() => {
          var provider = new firebase.auth.GoogleAuthProvider();
          firebase
            .auth()
            .signInWithPopup(provider)
            .then((result) => {
              router.push('/console/courses');
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

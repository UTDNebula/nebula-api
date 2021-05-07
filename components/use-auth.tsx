import React, { useState, useEffect, useContext, createContext } from 'react';
import firebase from '../lib/firebase';
import nookies from 'nookies';

const authContext = createContext(null);

/**
 * Auth provider for all front-end pages
 */
export default function ProvideAuth({ children }) {
  const auth = useProvideAuth();
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
}

/**
 * Provides auth context for components
 */
export const useAuth = () => {
  return useContext(authContext);
};

function useProvideAuth() {
  const [user, setUser] = useState(null);

  const signout = () => {
    return firebase
      .auth()
      .signOut()
      .then(() => {
        setUser(false);
      });
  };

  useEffect(() => {
    const unsubscribe = firebase.auth().onAuthStateChanged(async (user) => {
      if (user) {
        const token = await user.getIdToken();
        setUser(user);
        nookies.destroy(null, 'token');
        nookies.set(null, 'token', token, {});
      } else {
        setUser(false);
        nookies.destroy(null, 'token');
        nookies.set(null, 'token', '', {});
      }
    });

    return () => unsubscribe();
  }, []);

  return {
    user,
    signout,
  };
}

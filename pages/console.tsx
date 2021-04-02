import React from 'react';
import { useState, useEffect } from 'react';
import { useAuth } from '../components/use-auth';
import { useRouter } from 'next/router';
import Courses from './console/courses';
import Announcements from './console/announcements';

const AuthenticatedPage: React.FunctionComponent = () => {
  const [showCourses, setShowCourses] = useState(true);

  const auth = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!auth.user) {
      router.push('/');
    }
  }, []);

  return (
    <div>
      <div className="bg-blue-100 px-8 py-4 flex">
        <h1 className="text-xl font-bold flex-1 place-self-center">Admin Console</h1>
        <button
          className="btn1 mr-4"
          onClick={() => {
            setShowCourses(!showCourses);
          }}
        >
          Show {!showCourses ? "Courses" : "Announcements"}
        </button>
        <button
          className="btn1"
          onClick={() => {
            auth.signout().then((_) => router.push('/'));
          }}
        >
          Sign out
        </button>
      </div>
      {auth.user ? <>{showCourses ? <Courses /> : <Announcements />}</> : <p>No</p>}
    </div>
  );
};

export default AuthenticatedPage;

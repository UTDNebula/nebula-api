import React from 'react';
import { useState, useEffect } from 'react';
import { useAuth } from '../components/use-auth';
import { useRouter } from 'next/router';
import Courses from '../components/courses/courses';
import Announcements from '../components/announcements/announcements';

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
          className="mr-4 p-2 font-light rounded-xl bg-blue-300 hover:bg-blue-500"
          onClick={() => {
            setShowCourses(!showCourses);
          }}
        >
          Show {!showCourses ? "Courses" : "Announcements"}
        </button>
        <button
          className="p-2 font-light rounded-xl bg-blue-300 hover:bg-blue-500"
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

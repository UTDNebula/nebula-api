import React, { useEffect } from 'react';
import { useState } from 'react';
import Course from '../../components/courses/course';
import Modal from '../../components/courses/courseModal';
import { useAuth } from '../../components/use-auth';
import { useRouter } from 'next/router';
import Link from 'next/link';

const Courses: React.FunctionComponent = () => {
  const [input, setInput] = useState('');
  const [data, setData] = useState([]);
  const [open, setOpen] = useState(null);
  const [method, setMethod] = useState('');
  const [message, setMessage] = useState('No results found.');

  const auth = useAuth();
  const router = useRouter();

  const search = async () => {
    setMessage(`Searching database for ${input}...`);
    setData(await fetch(`/api/courses/search?name=${input}`).then((res) => res.json()));
    setMessage('No results found.');
  };

  const editCourse = (course) => {
    setMethod('PUT');
    setOpen(course);
  };

  const addModal = (course) => {
    setMethod('POST');
    setOpen({
      course: '',
      description: '',
      titleLong: '',
      hours: '',
      inclass: '',
      number: '',
      outclass: '',
      period: '',
      prerequisites: '',
      title: '',
    });
  };

  const deleteCourse = async (id) => {
    await fetch(`/api/courses/${id}`, {
      method: 'DELETE',
    }).then((res) => console.log(res.json()));
  };

  const close = async (course, update = false) => {
    if (update) {
      await fetch(`/api/courses/${course.id}`, {
        method: method,
        body: JSON.stringify(course),
      })
        .then((res) => res.json())
        .then((msg) => console.log(msg))
        .catch((err) => console.error(err));
    }
    setOpen(null);
    setMethod('');
  };

  useEffect(() => {
    if (!auth.user) {
      router.push('/');
    }
  }, []);

  return (
    <div>
      <div className="bg-blue-100 px-8 py-4 flex">
        <h1 className="text-xl font-bold flex-1 place-self-center">Admin Console</h1>
        <Link href="/console/announcements">
          <button className="btn1 mr-4">Show Announcements</button>
        </Link>
        <button
          className="btn1"
          onClick={() => {
            auth.signout().then((_) => router.push('/'));
          }}
        >
          Sign out
        </button>
      </div>
      <div>
        {open ? <Modal info={open} close={close} /> : <></>}
        <>
          <div className="m-8">
            <div className="flex mb-8">
              <input
                value={input}
                className="ring-blue-200 mr-4 py-2 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block shadow-md focus:outline-none ring-2 focus:ring-blue-600"
                placeholder="search term"
                onInput={(e) => {
                  const value = e.currentTarget.value;
                  setInput(value);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') search();
                }}
              ></input>
              <button className="btn1" onClick={search}>
                Search
              </button>
              <button className="btn1 mx-4" onClick={addModal}>
                Add Course
              </button>
            </div>
            <div className="mt-4 grid grid-cols-1 gap-8">
              {data && data.length != 0 ? (
                data.map((course) => {
                  return (
                    <Course
                      key={course.id}
                      course={course}
                      editCourse={editCourse}
                      deleteCourse={deleteCourse}
                    />
                  );
                })
              ) : (
                <p className="text-center">{message}</p>
              )}
            </div>
          </div>
        </>
      </div>
    </div>
  );
};

export default Courses;

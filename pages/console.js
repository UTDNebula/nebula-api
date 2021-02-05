import React from "react";
import { useState, useEffect } from 'react';
import { useAuth } from './use-auth';
import { useRouter } from 'next/router';
import Course from './course';
import Modal from './modal';

const AuthenticatedPage = (props) => {

    const [input, setInput] = useState('');
    const [data, setData] = useState('');
    const [open, setOpen] = useState(null);

    const auth = useAuth();
    const router = useRouter();

    const search = async () => {
        setData(await fetch(`/api/courses/search?name=${input}`).then(res => res.json()));
    }

    const editCourse = (course) => {
        setOpen(course);
    }

    const deleteCourse = (id) => {
        console.log("deleting " + id);
    }

    useEffect(() => {
        if(!auth.user) {
            router.push("/");
        }
    })

    const close = async (course, update=false) => {
        if (update) {
            await fetch(`/api/courses/${course.id}`, {
                method: 'PUT',
                body: JSON.stringify(course)
            }).then(res => res.json()).then(msg => console.log(msg));
        }
        setOpen(null);
    }

    return (
        <div>
            {open ? <Modal info={open} close={close} /> : <></>}
            {auth.user ? 
            <>
            <div className="bg-blue-100 px-8 py-4 flex">
                <h1 className="text-xl font-bold flex-1 place-self-center">Admin Console</h1>
                <button className="p-2 font-light rounded-xl bg-blue-300 hover:bg-blue-500"
                onClick={() => {
                    auth.signout().then(_ => router.push("/"));
                }}>Sign out</button>
            </div>
            <div className="m-8">
                <div className="flex mb-8">
                    <input 
                        value={input} 
                        className="ring-blue-200 mr-4 py-2 px-4 bg-white rounded-lg placeholder-gray-400 text-gray-900 appearance-none inline-block w-full shadow-md focus:outline-none ring-2 focus:ring-blue-600" 
                        placeholder="search term"
                        onInput={e => setInput(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') search();
                        }}
                    ></input>
                    <button className="p-2 font-light rounded-lg bg-blue-300 hover:bg-blue-500"
                        onClick={search}
                    >Search</button>
                </div>
                <div className="mt-4 grid grid-cols-1 gap-8">
                    {data ?
                        data.map((course) => {
                            return <Course key={course.id} course={course} editCourse={editCourse} deleteCourse={deleteCourse} />
                        })
                    : <p className="text-center">No results found.</p>
                    }
                </div>
            </div>
            </>
            : <p>No</p>}
        </div>
    );
};

export default AuthenticatedPage;
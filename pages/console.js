import React from "react";
import nookies from "nookies";
import { auth } from "../lib/firebaseAdmin";
import useSWR from 'swr';
import { useAuth } from './use-auth';
import { useRouter } from 'next/router';

export const getServerSideProps = async (ctx) => {
    try {
        const cookies = nookies.get(ctx);
        const token = await auth.verifyIdToken(cookies.token);
        // the user is authenticated
        return {
            props: {},
        };
    } catch (err) {
        // user not logged in
        return {
            redirect: {
                permanent: false,
                destination: "/",
            },
            props: {}
        };
    }
};

const fetcher = async (...args) => {
    const res = await fetch(...args);
    return res.json();
};

const AuthenticatedPage = (props) => {
    const { data } = useSWR(`/api/courses/10`, fetcher);
    const auth = useAuth();
    const router = useRouter();

    if (!data) {
        return 'Loading...';
    }

    return (
        <div className="p-8">
            <button onClick={() => {
                auth.signout().then(_ => router.push("/"));
            }}>Sign out</button>
            <h1 className="text-2xl">Title: {data.titleLong}</h1>
            <p className="text-blue-600">Description: {data.description}</p>
        </div>
    );
};

export default AuthenticatedPage;
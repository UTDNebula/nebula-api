import Head from 'next/head';
import firebase from '../lib/firebase';
import { useAuth } from './use-auth';
import nookies from "nookies";
import { auth } from "../lib/firebaseAdmin";
import Link from "next/link";
import { useRouter } from 'next/router'


export const getServerSideProps = async (ctx) => {
    try {
        const cookies = nookies.get(ctx);
        const token = await auth.verifyIdToken(cookies.token);
        return {
            redirect: {
                permanent: false,
                destination: "/console",
            },
            props: {}
        };
    } catch (err) {
        // user not logged in
        return { props: {} };
    }
}

export default function Home() {

    // const auth = useAuth();
    const router = useRouter();

    return (
        <div>
            <p>Please log in:</p>
            <button className="button" onClick={() => {
                var provider = new firebase.auth.GoogleAuthProvider();
                firebase.auth()
                    .signInWithPopup(provider)
                    .then((result) => {
                        console.log(result.user);
                        router.push("/console")
                    }).catch((error) => {
                        console.error(error);
                    });
            }}>Click me</button>
        </div>
    );
}

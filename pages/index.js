import firebase from '../lib/firebase';
import { useEffect } from 'react';
import { useAuth } from './use-auth';
import { useRouter } from 'next/router'

export default function Home() {

    const auth = useAuth();
    const router = useRouter();

    useEffect(() => {
        if(auth.user) {
            router.push("/console");
        }
    })

    return (
        <div className="min-h-screen min-w-screen flex">
            <button className="button rounded-lg p-4 bg-blue-200 hover:bg-blue-400 mx-auto place-self-center text-xl font-light" onClick={() => {
                var provider = new firebase.auth.GoogleAuthProvider();
                firebase.auth()
                    .signInWithPopup(provider)
                    .then((result) => {
                        console.log(result.user);
                        router.push("/console")
                    }).catch((error) => {
                        console.error(error);
                    });
            }}>Log In</button>
        </div>
    );
}

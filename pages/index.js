import Head from 'next/head';
import Console from './console';

const Home = () => {
    return (
        <div>
            <Head>
                <title>Admin Console</title>
            </Head>
            <Console />
        </div>
    );
};

export default Home;

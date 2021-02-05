import 'tailwindcss/tailwind.css';
import ProvideAuth from './use-auth';
function MyApp({ Component, pageProps }) {
    return (
        <ProvideAuth>
            <Component {...pageProps} />
        </ProvideAuth>
    );
}

export default MyApp;

import useSWR from 'swr';

const fetcher = async (...args) => {
    const res = await fetch(...args);

    return res.json();
};

function Console() {
    const { data } = useSWR(`/api/courses/10`, fetcher);

    if (!data) {
        return 'Loading...';
    }

    return (
        <div class="p-8">
            <h1 class="text-2xl">Title: {data.titleLong}</h1>
            <p class="text-blue-600">Description: {data.description}</p>
        </div>
    );
}

export default Console;

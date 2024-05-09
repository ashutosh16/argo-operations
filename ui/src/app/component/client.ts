//Wrapper for fetch
export const apiCall = (url: string, headers: Record<string, string>) => {
    return fetch(url, { headers })
        .then((res) => res.json())
        .then((res) => res)
        .catch((err) => {
            throw err;
        });
};
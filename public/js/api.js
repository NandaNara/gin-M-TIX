window.apiFetch = async function apiFetch(path, options = {}) {
    const headers = new Headers(options.headers || {});
    const token = localStorage.getItem('mtix_token');

    if (token) headers.set('Authorization', `Bearer ${token}`);
    if (options.body && !(options.body instanceof FormData) && !headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json');
    }

    const response = await fetch(path, { ...options, headers });
    const isJSON = response.headers.get('content-type')?.includes('application/json');
    const data = isJSON ? await response.json() : await response.blob();

    if (!response.ok) {
        if (response.status === 401) localStorage.removeItem('mtix_token');
        throw new Error(data?.error || `Request failed (${response.status})`);
    }
    return data;
};

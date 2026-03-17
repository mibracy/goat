let jwtToken = localStorage.getItem("jwtToken") || "";

export const updateJwtToken = (token) => {
    jwtToken = token;
    if (token) {
        localStorage.setItem("jwtToken", token);
    } else {
        localStorage.removeItem("jwtToken");
    }
};

export const getUserIdFromJwt = () => {
    try {
        const base64Url = jwtToken.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
        const payload = JSON.parse(jsonPayload);
        return payload.sub; // 'sub' contains the user ID
    } catch (e) {
        console.error("Error decoding JWT:", e);
        return null;
    }
};

export async function callApi(method, path, body = null, setApiResponse) {
    const baseUrl = window.location.origin;
    const options = { method };
    options.headers = { "Content-Type": "application/json" };

    if (body) {
        options.body = JSON.stringify(body);
    }

    if (jwtToken && path !== `${baseUrl}/login` && path !== `${baseUrl}/register` && path !== `${baseUrl}/forgot-password` && path !== `${baseUrl}/reset-password`) {
        options.headers["Authorization"] = `Bearer ${jwtToken}`;
    }

    try {
        const response = await fetch(path, options);
        const textData = await response.text();

        try {
            const jsonData = JSON.parse(textData);
            setApiResponse(JSON.stringify(jsonData, null, 2));
            return jsonData; // Return parsed JSON data for further processing
        } catch (jsonError) {
            setApiResponse(textData); // If not JSON, display as plain text
            return textData; // Return raw text data
        }
    } catch (error) {
        setApiResponse("Error: " + error.message);
        console.error("API Call Error:", error);
        return null; // Indicate error
    }
}

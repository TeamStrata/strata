import router from "@/router";

//this line is proof that there either is no god, or he is powerless to prevent what I have done
const URL = (import.meta.env.MODE === "development" ? "http://localhost:8080" : window.location.origin);


export async function apiFetch(route, verb = 'GET', body = false, content = "application/json") {
    if (route != "/login") {
        route = "/api" + route;
    }

    const res = await fetch(URL + route, {
        method: verb,
        credentials: "include",
        headers: {
            'Content-Type': `${content}`,
        },
        body: body ? body : null
    })

    if (res.status == 401) {
        // Go back to login
        router.push('/login');
    }


    return res;
}
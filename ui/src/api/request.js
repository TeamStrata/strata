import router from "@/router";

//this line is proof that there either is no god, or he is powerless to prevent what I have done
const URL = (import.meta.env.MODE === "development" ? "http://localhost:8080/" : window.location.href).slice(0, -1);


export async function apiFetch(route, verb = 'GET', body = false, content = "application/json") {
    const res = await fetch(URL + route, {
        method: verb,
        credentials: "include",
        headers: {
            'Content-Type': `${content}`,
        },
        body: body ? body : null
    })
    console.log(res.status)

    if (route != "/login") {
        // handle not logged in from request
        if (res.status == 401) {
            router.push('/');
        }
    }

    return res;
}
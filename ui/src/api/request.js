import router from "@/router";

export async function apiFetch(route, verb, body) {
    const res = await fetch("http://localhost:8080" + route, {
        method: verb,
        credentials: "include",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body)    
    })
    console.log(res.status)

    if(route != "/login") {
        // handle not logged in from request
        if(res.status == 401) {
            router.push('/');
        }
    } 

    return res;
}
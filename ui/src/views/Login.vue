<script setup>
import router from '@/router';
import { ref } from 'vue';

var user = ref("");
var pass = ref("");
var error = ref(false);

function login() {
    // submit user and pass to form
    fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
            username: user.value,
            password: pass.value
        }),
    })
        .then(response => {
            if (!response.ok) {
                error = true;
                throw new Error('Something went wrong');
            }
            router.push("/")
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

</script>

<template>
    <div class="flex h-screen bg-gray-100">
        <div class="flex flex-col w-96 m-auto p-8 bg-white shadow-md rounded-lg">
            <h1 class="text-2xl font-semibold text-gray-800 mb-4">Log In</h1>
            <label for="username" class="text-gray-600">Username</label>
            <input id="username" v-model="user"
                class="mt-1 border border-gray-300 rounded-md p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter your username">

            <label for="password" class="mt-4 text-gray-600">Password</label>
            <input id="password" type="password" v-model="pass"
                class="mt-1 border border-gray-300 rounded-md p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter your password">

            <button @click="login"
                class="mt-6 bg-blue-500 text-white font-semibold py-2 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200">Login</button>

            <p v-if="error" class="mt-4 text-red-500">Something went wrong</p>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
// placeholder user data 
let users =
    [
        { name: "Alistair", stuff: "hello" },
        { name: "Sebastian", stuff: "hello" }
    ]
var newUsername = ref("");
var newPassword = ref("");

var addModal = ref(false);

function addUser() {
    // submit user and pass to form
    fetch('http://localhost:8080/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: newUsername.value,
            password: newPassword.value
        }),
    })
        .then(response => {
            if (!response.ok) {
                error = true;
                throw new Error('Something went wrong');
            }
            let data = response.json();
            console.log(data)
        })
        .catch(error => {
            console.error('Error:', error);
        });
}
</script>

<template>
    <!-- add user modal -->
    <div class="static" v-if="addModal">
        <div class="bg-black/20 w-screen h-screen absolute left-0 top-0 backdrop-blur-xs" @click="addModal = false">
        </div>
        <div class="bg-white p-5 left-1/2 transform -translate-x-1/2 top-1/3 -translate-y-1/2 absolute rounded-lg">
            <h1>Add User</h1>
            <div>
                <label for="newUser">Username:</label>
                <input id="newUser" v-model="newUsername">
            </div>
            <div>
                <label for="newPass">Password:</label>
                <input id="newPass" type="password" v-model="newPassword">
            </div>
            <div class="flex flex-row justify-between mt-4">
                <button class="bg-neutral-200 px-2 rounded-sm cursor-pointer" @click="addModal = false">Cancel</button>
                <button class="bg-blue-100 rounded-sm w-fit px-2 cursor-pointer" @click="addUser">Add User</button>
            </div>

        </div>
    </div>

    <h1 class="font-bold">Manage Users</h1>


    <button @click="addModal = true" class="bg-amber-200 p-3 rounded-xl cursor-pointer">Add User</button>

    <div class="flex flex-col">
    </div>

    Users list placeholder:
    <ul>
        <li v-for="u in users">
            {{ u.name }}
        </li>
    </ul>

</template>

<style scoped>
input {
    background-color: rgb(230, 230, 230);
    margin: 1px;
}
</style>
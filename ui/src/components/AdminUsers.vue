<script setup>
import { ref } from 'vue';
import Toast, { ToastTypes } from './Toast.vue';
// placeholder user data 
let users =
    [
        { name: "Alistair", stuff: "hello" },
        { name: "Sebastian", stuff: "hello" }
    ]
var newUsername = ref("");
var newPassword = ref("");

var addModal = ref(false);

const toastRef = ref(null)

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
                toastRef.value?.showToast("There was an issue with creating the account", ToastTypes.FAIL)
                throw new Error('Something went wrong');
            }
            // let data = response.json();
            // console.log(data)
            toastRef.value?.showToast("The account was successfully created", ToastTypes.SUCCESS)
        })
        .catch(error => {
            console.error('Error:', error);
        });
    addModal.value = false
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

    <Toast ref="toastRef" />

    <p>Maybe some flavor text</p>

    <!-- list header -->
    <div class="flex flex-row justify-between items-center mb-2">
        <h1 class="font-bold">Manage Users</h1>


        <div class="flex flex-row pb-2 items-center *:mx-3">
            <p>{{ users.length }} users</p>
            <input placeholder="search" class="p-2 rounded-lg">
            <button @click="addModal = true" class="bg-amber-200 p-2 px-4 rounded-lg cursor-pointer">Add User</button>
        </div>

    </div>

    <!-- list body -->
    <ul>
        <li class="test border-t-2 border-neutral-200 py-3 flex flex-row justify-between px-5" v-for="u in users">
            {{ u.name }}
            <div class="flex items-center">
                <p class="mr-25">User Type</p>
                <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 0 20 20">
                    <path fill="currentColor"
                        d="M5 6a4 4 0 1 1 8 0a4 4 0 0 1-8 0m-3 7c0-1.113.903-2 2.009-2h6.248A5.48 5.48 0 0 0 9 14.5c0 1.303.453 2.5 1.21 3.443Q9.617 18 9 18c-1.855 0-3.583-.386-4.865-1.203C2.833 15.967 2 14.69 2 13m17 1.5a4.5 4.5 0 1 1-9 0a4.5 4.5 0 0 1 9 0m-2.646-1.146a.5.5 0 0 0-.708-.708L14.5 13.793l-1.146-1.147a.5.5 0 0 0-.708.708l1.147 1.146l-1.147 1.146a.5.5 0 0 0 .708.708l1.146-1.147l1.146 1.147a.5.5 0 0 0 .708-.708L15.207 14.5z" />
                </svg>
            </div>
        </li>
    </ul>

</template>

<style scoped>
input {
    background-color: rgb(230, 230, 230);
    margin: 1px;
}
</style>
<script setup>
import { ref } from 'vue';
import Badge from './Badge.vue';
import Card from './Card.vue';
import CardDotted from './CardDotted.vue';
import Modal from './Modal.vue';

const roles = [
    { "ID": 1, "name": "Admin", "color": "#FF0000" },
    { "ID": 2, "name": "Moderator", "color": "#00FF00" },
    { "ID": 3, "name": "User", "color": "#0000FF" }
]

const editModal = ref(false)
const createModal = ref(false)
const tempRole = ref(null)

function showEdit(id) {

    tempRole.value = JSON.parse(JSON.stringify(roles.find(item => { return item.ID == id })))
    editModal.value = true;
}

function submitEdit() {
    // TODO: post to database in here
}

function showCreate() {
    createModal.value = true;
}

function submitCreate() {
    // TODO: post to database in here
}
</script>

<template>

    <Modal :show="editModal" @close="editModal = false">
        <div class="flex flex-col space-y-6">
            <div class="flex flex-col justify-between">
                <input class="text-2xl font-semibold bg-neutral-200/60 p-2 rounded-md text-neutral-800"
                    :value="tempRole.name">
                <p>Color select</p>
            </div>
            <div class="space-y-2">
                <h2 class="text-lg font-medium text-neutral-700">Permissions</h2>
                <div class="flex items-center justify-between py-2">
                    <span class="text-sm text-neutral-700">whatever</span>
                    <input type="checkbox" class="form-checkbox w-4 h-4" />
                </div>
            </div>
            <div class="flex justify-between space-x-4 pt-4 mt-4">
                <button @click="editModal = false"
                    class="cursor-pointer px-4 py-2 rounded border border-neutral-400 text-neutral-600 hover:bg-neutral-100 transition">
                    Cancel
                </button>
                <button @click="submitEdit"
                    class="cursor-pointer px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 transition">
                    Confirm
                </button>

            </div>
        </div>
    </Modal>

    <Modal :show="createModal" @close="createModal = false">
        <div class="flex flex-col space-y-6">
            <div class="flex flex-col justify-between">
                <input class="text-2xl font-semibold bg-neutral-200/60 p-2 rounded-md text-neutral-800">
                <p>Color select</p>
            </div>
            <div class="space-y-2">
                <h2 class="text-lg font-medium text-neutral-700">Permissions</h2>
                <div class="flex items-center justify-between py-2">
                    <span class="text-sm text-neutral-700">whatever</span>
                    <input type="checkbox" class="form-checkbox w-4 h-4" />
                </div>
            </div>
            <div class="flex justify-between space-x-4 pt-4 mt-4">
                <button @click="createModal = false"
                    class="cursor-pointer px-4 py-2 rounded border border-neutral-400 text-neutral-600 hover:bg-neutral-100 transition">
                    Cancel
                </button>
                <button @click="submitCreate"
                    class="cursor-pointer px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700 transition">
                    Confirm
                </button>

            </div>
        </div>
    </Modal>

    <!-- list header -->
    <div class="flex flex-row justify-between items-center mb-6">
        <h1 class="text-2xl font-semibold text-gray-800">Roles</h1>
        <div class="flex flex-row pb-2 items-center space-x-3">
            <p class="text-gray-600">{{ roles.length }} Roles</p>
            <input placeholder="Search"
                class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
                type="text">
            <!-- <button @click="createModal = true"
                class="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200 cursor-pointer">
                Add Role
            </button> -->
        </div>
    </div>

    <ul class="grid grid-cols-1 lg:grid-cols-2 auto-rows-fr gap-4 items-start">
        <li v-for="r in roles" class="h-full">
            <Card class="h-full">
                <div class="flex items-center justify-between">
                    <h1 class="text-lg">{{ r.name }}</h1>
                    <Badge :color="r.color">{{ r.color }}</Badge>
                </div>
                <div class="flex justify-end pt-5">
                    <button class="outline-1 rounded-sm px-2 py-1 outline-purple-800 text-purple-800 cursor-pointer"
                        @click="showEdit(r.ID)">Edit Role</button>
                </div>
            </Card>
        </li>
        <li class="h-full cursor-pointer" @click="showCreate">
            <CardDotted
                class="flex flex-col items-center justify-center text-neutral-500 h-full hover:text-purple-800 hover:outline-purple-800 transition-colors duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 50 50">
                    <path fill="currentColor"
                        d="M25 42c-9.4 0-17-7.6-17-17S15.6 8 25 8s17 7.6 17 17s-7.6 17-17 17m0-32c-8.3 0-15 6.7-15 15s6.7 15 15 15s15-6.7 15-15s-6.7-15-15-15" />
                    <path fill="currentColor" d="M16 24h18v2H16z" />
                    <path fill="currentColor" d="M24 16h2v18h-2z" />
                </svg>
                <h2>Create a New Role</h2>
            </CardDotted>
        </li>
    </ul>
</template>
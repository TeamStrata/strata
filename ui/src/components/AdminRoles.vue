<script setup>
import Badge from './Badge.vue';
import Card from './Card.vue';
import CardDotted from './CardDotted.vue';

const roles = [
    { "ID": 1, "name": "Admin", "color": "#FF0000" },
    { "ID": 2, "name": "Moderator", "color": "#00FF00" },
    { "ID": 3, "name": "User", "color": "#0000FF" }
]


</script>

<template>
    <!-- list header -->
    <div class="flex flex-row justify-between items-center mb-6">
        <h1 class="text-2xl font-semibold text-gray-800">Roles</h1>
        <div class="flex flex-row pb-2 items-center space-x-3">
            <p class="text-gray-600">{{ roles.length }} Roles</p>
            <input placeholder="Search"
                class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
                type="text">
            <button @click="addModal = true"
                class="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200 cursor-pointer">
                Add Role
            </button>
        </div>
    </div>

    <ul class="grid grid-cols-1 lg:grid-cols-2 auto-rows-fr gap-4 items-start">
        <li v-for="r in roles" class="h-full">
            <Card class="h-full" @click="$router.push(`roleEdit/${r.ID}`)">
                <div class="flex items-center justify-between">
                    <h1 class="text-lg">{{ r.name }}</h1>
                    <Badge>{{ r.color }}</Badge>
                </div>
                <div class="flex justify-end pt-5">
                    <button class="outline-1 rounded-sm px-2 py-1 outline-purple-800 text-purple-800 cursor-pointer"
                        @click="">Edit Role</button>
                </div>
            </Card>
        </li>
        <li class="h-full cursor-pointer">
            <CardDotted class="flex flex-col items-center justify-center text-neutral-500 h-full hover:text-purple-800 hover:outline-purple-800 transition-colors duration-300">
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
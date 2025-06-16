<script setup>
import { ref } from "vue";
import Toast, { ToastTypes } from "./Toast.vue";
import { apiFetch } from "@/api/request";
import Modal from "./Modal.vue";
import RoleMultiSelect from "./RoleMultiSelect.vue";

// If we can't connect to the database, or while the fetch is in progress, this displays.
var users = ref([]);
var newUsername = ref("");
var newPassword = ref("");

var addModal = ref(false);
const toastRef = ref(null);
const roleList = ref([])

function addUser() {
	const body = {
		username: newUsername.value,
		password: newPassword.value,
	}

	// submit user and pass to form
	apiFetch('/signup', 'POST', JSON.stringify(body))
		.then((response) => {
			if (!response.ok) {
				toastRef.value?.showToast(
					"There was an issue with creating the account",
					ToastTypes.FAIL,
				);
				throw new Error("Something went wrong");
			}

			toastRef.value?.showToast(
				"The account was successfully created",
				ToastTypes.SUCCESS,
			);

			// Push new user to the list
			users.value.push({ username: newUsername.value });
		})
		.catch((error) => {
			console.error("Error:", error);
		});
	addModal.value = false;
}

function deleteUser(username) {

	apiFetch(`/user/${username}`, "DELETE")
		.then((res) => {
			if (!res.ok) {
				toastRef.value?.showToast(
					"There was an issue deleting the user",
					ToastTypes.FAIL,
				);
				throw new Error("Error deleting user");
			}

			// Handle success
			toastRef.value?.showToast(
				"The user was succesfully deleted",
				ToastTypes.SUCCESS,
			);

			// Reload the user list
			loadUsers();
		})
		.catch((err) => {
			console.error(err);
		});
}

function loadUsers() {
	apiFetch('/users', 'GET')
		.then(async (response) => {
			// Handle error
			if (!response.ok) {
				toasRef.value?.showToast(
					"There was an error when loading users",
					ToastTypes.FAIL,
				);
				throw new Error("Error loading users");
			}

			users.value = await response.json();
		})
		.catch((error) => {
			console.error(error);
		});
}

function loadRoles() {
	apiFetch('/roles', 'GET').then(async res => {
		let val = await res.json();
		roleList.value = val;

	}).catch(e => {
		console.error(e)
	}
	)
}

function addUserRole() {

}

function removeUserRole() {
	
}

const isLoaded = ref(false);

// Load both users and roles, then set isLoaded to true
Promise.all([loadUsers(), loadRoles()]).then(() => {
	isLoaded.value = true;
});


</script>

<template>
	<!-- add user modal -->
	<Modal :show="addModal" @close="addModal = false">
		<h1 class="text-2xl font-semibold text-gray-800 mb-6">Add User</h1>
		<div>
			<label for="newUser" class="text-gray-600">Username</label>
			<input id="newUser" v-model="newUsername"
				class="mt-1 border border-gray-300 rounded-md p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="Enter username" />
		</div>
		<div class="mt-4">
			<label for="newPass" class="text-gray-600">Password</label>
			<input id="newPass" type="password" v-model="newPassword"
				class="mt-1 border border-gray-300 rounded-md p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="Enter password" />
		</div>
		<div class="flex flex-row justify-between mt-6">
			<button
				class="bg-neutral-200 px-4 py-2 rounded-md cursor-pointer hover:bg-neutral-300 focus:outline-none focus:ring-2 focus:ring-gray-300"
				@click="addModal = false">
				Cancel
			</button>
			<button
				class="bg-blue-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200"
				@click="addUser">
				Add User
			</button>
		</div>
	</Modal>

	<Toast ref="toastRef" />

	<p class="text-gray-700 mb-4">Manage your organizations members</p>

	<!-- list header -->
	<div class="flex flex-row justify-between items-center mb-6">
		<h1 class="text-2xl font-semibold text-gray-800">Members</h1>

		<div class="flex flex-row pb-2 items-center space-x-3">
			<p class="text-gray-600">{{ users.length }} members</p>
			<input placeholder="Search"
				class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
				type="text" />
			<button @click="addModal = true"
				class="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200 cursor-pointer">
				Add Member
			</button>
		</div>
	</div>

	<!-- list body -->
	<table class="w-full" v-if="isLoaded">
		<tr class="border-t-2 border-neutral-200 py-4 flex flex-row justify-between items-center px-5 hover:bg-gray-50 transition-colors"
			v-for="(u, index) in users" :key="index">
			<td class="text-gray-800 font-medium">{{ u.username }}</td>
			<td>
				<RoleMultiSelect :activeRoles="u.role" :allVals="roleList" @roleAdd="(role) => {console.log(role)}"></RoleMultiSelect>
			</td>
			<td class="flex items-center">
				Hi
				<!-- <button
					class="bg-red-500 text-white py-2 px-4 rounded-md cursor-pointer hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-200"
					@click="deleteUser(u.username)" v-if="u.username != 'admin'">
					<svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 0 20 20"
						class="text-gray-600">
						<path fill="#fff"
							d="M5 6a4 4 0 1 1 8 0a4 4 0 0 1-8 0m-3 7c0-1.113.903-2 2.009-2h6.248A5.48 5.48 0 0 0 9 14.5c0 1.303.453 2.5 1.21 3.443Q9.617 18 9 18c-1.855 0-3.583-.386-4.865-1.203C2.833 15.967 2 14.69 2 13m17 1.5a4.5 4.5 0 1 1-9 0a4.5 4.5 0 0 1 9 0m-2.646-1.146a.5.5 0 0 0-.708-.708L14.5 13.793l-1.146-1.147a.5.5 0 0 0-.708.708l1.147 1.146l-1.147 1.146a.5.5 0 0 0 .708.708l1.146-1.147l1.146 1.147a.5.5 0 0 0 .708-.708L15.207 14.5z" />
					</svg>
				</button> -->
			</td>
		</tr>
	</table>
</template>

<style scoped>
input {
	background-color: rgb(230, 230, 230);
	margin: 1px;
}
</style>

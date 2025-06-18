<script setup>
import { ref } from "vue";
import Toast, { ToastTypes } from "./Toast.vue";
import { apiFetch } from "@/api/request";
import Modal from "./Modal.vue";
import RoleMultiSelect from "./RoleMultiSelect.vue";
import { Button } from '@/components/ui/button';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuPortal,
	DropdownMenuSeparator,
	DropdownMenuShortcut,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogClose,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';


// If we can't connect to the database, or while the fetch is in progress, this displays.
var users = ref([]);
var newUsername = ref("");
var newPassword = ref("");

var addModal = ref(false);
const toastRef = ref(null);
const roleList = ref([])

let tempUser = ref({});

let showDelete = ref(false);
let userToDelete = ref(null);

function addUser(user) {
	const body = {
		username: user.username,
		password: user.password,
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

function addUserRole(uid, rid) {
	apiFetch(`/user/${uid}/role/${rid}`, 'POST').then((res) => {
		if (res.ok) {
			loadUsers()
		} else {
			console.error("Something went wrong");
		}
	})
}

function removeUserRole(uid, rid) {
	apiFetch(`/user/${uid}/role/${rid}`, 'DELETE').then((res) => {

		if (res.ok) {
			loadUsers()
		} else {
			console.error("Something went wrong");
		}
	})
}

function showCreate() {
	tempUser.value = {};
	addModal.value = true;
}

// function showEdit(id) {
//     console.log(id)
//     const userToEdit = users.value.find(item => { return user.id == id });
//     tempUser.value = JSON.parse(JSON.stringify(userToEdit));
//     originalRole = JSON.parse(JSON.stringify(userToEdit));
//     editModal.value = true;
// }

const isLoaded = ref(false);

// Load both users and roles, then set isLoaded to true
Promise.all([loadUsers(), loadRoles()]).then(() => {
	isLoaded.value = true;
});


</script>

<template>
	<!-- add user modal -->
	<Dialog :open="addModal" @update:open="addModal = $event">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Create New User</DialogTitle>
				<DialogDescription>
					Fill in the details for the new user.
				</DialogDescription>
			</DialogHeader>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="create-username-input" class="text-right">Name</Label>
					<Input id="create-username-input" type="text" v-model="tempUser.username" class="col-span-3" />
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="create-password-input" class="text-right">Password</Label>
					<Input id="create-password-input" type="text" v-model="tempUser.password"
						class="col-span-3 h-10 w-full" />
				</div>
			</div>
			<DialogFooter>
				<DialogClose as-child>
					<Button type="button" variant="outline">Cancel</Button>
				</DialogClose>
				<Button type="submit" @click="addUser(tempUser)">Confirm</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>

	<Dialog :open="showDelete" @update:open="showDelete = $event">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Delete User?</DialogTitle>
				<DialogDescription>
					This action cannot be undone. This will permanently
					delete the user.
				</DialogDescription>
			</DialogHeader>
			<DialogFooter class="flex justify-between space-x-4">
				<DialogClose as-child>
					<Button type="button" variant="outline">Cancel</Button>
				</DialogClose>
				<Button type="submit" @click="deleteUser(userToDelete)">Confirm</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>


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
			<button @click="showCreate"
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
				<RoleMultiSelect :activeRoles="u.role" :allVals="roleList"
					@roleAdd="(role) => { addUserRole(u.id, role) }" @roleRemove="(role) => { removeUserRole(u.id, role) }">
				</RoleMultiSelect>
			</td>
			<td class="flex items-center">
				<DropdownMenu>
					<DropdownMenuTrigger as-child>
						<Button variant="ghost">
							<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20">
								<path fill="currentColor" fill-rule="evenodd"
									d="M2.5 7.5a2.5 2.5 0 1 1 0 5a2.5 2.5 0 0 1 0-5m15 0a2.5 2.5 0 1 1 0 5a2.5 2.5 0 0 1 0-5m-7.274 0a2.5 2.5 0 1 1 0 5a2.5 2.5 0 0 1 0-5" />
							</svg>
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent class="w-56">
						<DropdownMenuGroup>
							<DropdownMenuItem>
								<span>Edit</span>
							</DropdownMenuItem>
							<DropdownMenuItem class="text-destructive focus:text-destructive"
								:onclick="() => { userToDelete = u.username; showDelete = true; }">
								<span>Delete</span>
							</DropdownMenuItem>
						</DropdownMenuGroup>
					</DropdownMenuContent>
				</DropdownMenu>
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

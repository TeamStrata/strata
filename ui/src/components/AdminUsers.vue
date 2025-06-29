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
import { computed } from "vue";
import Card from "./ui/card/Card.vue";


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

const editModal = ref(false);

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

// Send PATCH request to backend API
function submitEdit() {
	let body = {};

	body.id = tempUser.value.id;

	if (tempUser.value.username != "") {
		body.username = tempUser.value.username;
	}

	if (tempUser.value.password != "") {
		body.password = tempUser.value.password
	}

	let route = "/user/" + body.id;
	apiFetch(route, "PATCH", JSON.stringify(body))
		.then((response) => {
			if (!response.ok) {
				toastRef.value?.showToast(
					"There was an issue updating the role",
					ToastTypes.FAIL,
				);
				throw new Error("Unable to update role")
			} else {
				editModal.value = false;
				toastRef.value?.showToast(
					"User updated successfully",
					ToastTypes.SUCCESS,
				);
				loadUsers();
			}
		})
		.catch((error) => {
			console.error("Error:", error);
		});
}

function showEdit(user) {
	console.log(user)
	const userToEdit = users.value.find(item => { return item.id == user });;
	tempUser.value = JSON.parse(JSON.stringify(userToEdit));
	console.log(tempUser.value)
	// originalRole = JSON.parse(JSON.stringify(userToEdit));
	editModal.value = true;
}

const nameFilter = ref("");

const filteredUsers = computed(() => {
	if (!nameFilter.value) return users.value;
	return users.value.filter(u =>
		u.username?.toLowerCase().includes(nameFilter.value.toLowerCase())
	);
});

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

	<!-- edit user modal -->
	<Dialog :open="editModal" @update:open="editModal = $event">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Edit User</DialogTitle>
				<DialogDescription>
					Click 'Confirm' to save changes.
				</DialogDescription>
			</DialogHeader>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="edit-name-input" class="text-right">Name</Label>
					<Input id="edit-name-input" type="text" v-model="tempUser.username" class="col-span-3" />
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="edi-password-input" class="text-right">Password</Label>
					<Input id="edit-password-input" type="text" v-model="tempUser.password"
						class="col-span-3 h-10 w-full" />
				</div>
			</div>
			<DialogFooter>
				<DialogClose as-child>
					<Button type="button" variant="outline">
						Cancel
					</Button>
				</DialogClose>
				<Button type="submit" @click="submitEdit">Confirm</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>

	<!-- delete user modal -->
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
	<Card class="p-5">
		<!-- list header -->
		<div class="flex flex-row justify-between items-center mb-2">
			<h1 class="text-2xl font-semibold text-gray-800">All Members</h1>

			<div class="flex flex-row items-center gap-3">
				<p class="text-gray-600 text-nowrap">{{ users.length }} members</p>
				<Input id="edit-name-input" type="text" v-model="nameFilter" placeholder="Search..." />
				<Button type="submit" @click="showCreate"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
						viewBox="0 0 24 24">
						<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
							stroke-width="2" d="M5 12h14m-7-7v14" />
					</svg>Add Member</Button>
			</div>
		</div>

		<!-- list body -->
		<table class="w-full" v-if="isLoaded">
			<colgroup>
				<col span="1" style="width: 20%;">
				<col span="1" style="width: 75%;">
				<col span="1" style="width: 5%;">
			</colgroup>
			<tbody class="*:border-b-2 *:border-neutral-200 text-left">
				<tr>
					<th>Name</th>
					<th>Roles</th>
					<th>Actions</th>
				</tr>
				<tr class="transition-colors hover:bg-gray-50 align-cente" v-for="(u, index) in filteredUsers"
					:key="index">
					<td class="text-gray-800 font-medium py-3">{{ u.username }}</td>
					<td>
						<RoleMultiSelect :activeRoles="u.role" :allVals="roleList"
							@roleAdd="(role) => { addUserRole(u.id, role) }"
							@roleRemove="(role) => { removeUserRole(u.id, role) }">
						</RoleMultiSelect>
					</td>
					<td class="text-right">
						<DropdownMenu>
							<DropdownMenuTrigger as-child>
								<Button variant="ghost">
									<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20">
										<path fill="currentColor" fill-rule="evenodd"
											d="M2.5 7.5a2.5 2.5 0 1 1 0 5a2.5 2.5 0 0 1 0-5m15 0a2.5 2.5 0 1 1 0 5a2.5 2.5 0 0 1 0-5m-7.274 0a2.5 2.5 0 1 1 0 5a2.5 2.5 0 0 1 0-5" />
									</svg>
								</Button>
							</DropdownMenuTrigger>
							<DropdownMenuContent class="w-44">
								<DropdownMenuGroup>
									<DropdownMenuItem :onclick="() => { showEdit(u.id) }">
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
			</tbody>

		</table>
	</Card>
</template>

<style scoped>
input {
	background-color: rgb(230, 230, 230);
	margin: 1px;
}
</style>

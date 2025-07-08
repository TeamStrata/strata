<script setup>
import { ref } from 'vue';
import { apiFetch } from '@/api/request';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';

import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

import CardDotted from './CardDotted.vue';
import { Checkbox } from './ui/checkbox';
import Toast, { ToastTypes } from './Toast.vue';

// Component refs
const toastRef = ref(null);
const editModal = ref(false);
const createModal = ref(false);
const deleteModal = ref(false);

// Roles refs
const roles = ref([]);
let tempRole = ref({});
let originalRole = null;

// Permissions refs
const permissions = ref([]);
const selectedPermissions = ref([]);

// Show the edit role dialog
function showEdit(id) {
    console.log(id)
    const roleToEdit = roles.value.find(item => { return item.id == id });
    tempRole.value = JSON.parse(JSON.stringify(roleToEdit));
    originalRole = JSON.parse(JSON.stringify(roleToEdit));
    selectedPermissions.value = roleToEdit.permissions ? 
        roleToEdit.permissions.map(p => p.id) : [];
    editModal.value = true;
}

// Send PATCH request to backend API
function submitEdit() {
    let body = {};
    body.id = originalRole.id;

    if (originalRole.name != tempRole.value.name) {
        body.name = tempRole.value.name;
    }

    if (originalRole.color != tempRole.value.color) {
        body.color = tempRole.value.color;
    }

    const originalPermissionIds = (originalRole.permissions || []).map(p => p.id).sort();
    const currentPermissionIds = selectedPermissions.value.sort();
    if (JSON.stringify(originalPermissionIds) !== JSON.stringify(currentPermissionIds)) {
        body.permissions = selectedPermissions.value;
    }

    let route = "/role/" + body.id; 
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
                    "Role updated successfully",
                    ToastTypes.SUCCESS,
                );
                loadRoles();
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

// Reset component values
function showCreate() {
    tempRole.value = {};
    selectedPermissions.value = [];
    createModal.value = true;
}

// Post new role to backend
function submitCreate() {
    let body = {};
    if (tempRole.value.name != undefined) {
        body.name = tempRole.value.name;
    }

    if (tempRole.value.color != undefined) {
        body.color = tempRole.value.color;
    }

    if (tempRole.value.permissions != undefined) {
        body.permissions = tempRole.value.permissions;
    }

    if (selectedPermissions.value.length > 0) {
        body.permissions = selectedPermissions.value;
    }

    apiFetch("/role", "POST", JSON.stringify(body))
        .then((response) => {
            if (!response.ok) {
                toastRef.value?.showToast(
                    "There was an issue creating the new role",
                    ToastTypes.FAIL,
                );
                throw new Error("Failed to create role");
            } else {
                createModal.value = false;
                toastRef.value?.showToast(
                    "Role created  successfully",
                    ToastTypes.SUCCESS,
                )
                loadRoles();
            }
        })
        .catch((error) => {
            console.error(error)
        });
}

// Show delete role dialog
function showDelete(id) {
    const roleToEdit = roles.value.find(item => { return item.id == id });
    tempRole.value = JSON.parse(JSON.stringify(roleToEdit));
    deleteModal.value = true;
}

// Send DELETE role request to backend API
function submitDelete() {
    const route = "/role/" + tempRole.value.id;
    apiFetch(route, "DELETE")
        .then((response) => {
            if (!response.ok) {
                toastRef.value?.showToast(
                    "There was an issue deleting the role",
                    ToastTypes.FAIL,
                );
                throw new Error("Failed to delete role");
            } else {
                deleteModal.value = false;
                toastRef.value?.showToast(
                    "Role deleted successfully",
                    ToastTypes.SUCCESS,
                )
                loadRoles();
            }
        })
        .catch((error) => {
            console.error(error)
        });
}

// Get all global permissions
function loadPermissions() {
    const route = "/permissions/global";
    apiFetch(route)
    .then(async (response) => {
        if (!response.ok) {
            toastRef.value?.showToast(
                "There was an issue fetching permissions",
                ToastTypes.FAIL,
            )
            throw new Error("Failed to fetch permissions");
        } else {
            permissions.value = formatPermissions(await response.json());
        }
    })
    .catch((error) => {
        console.error(error);
    });
}

// Format permissions names
// (e.g. 'test_permission' -> 'Test Permission')
function formatPermissions(p) {
    let formatted = [];
    for (let i = 0; i < p.length; i++) {
        let oldName = p[i].name;
        let newName = oldName[0].toUpperCase() + oldName.slice(1, oldName.search('_')) + " " + oldName[oldName.search('_') + 1].toUpperCase() + oldName.slice(oldName.search('_') + 2, oldName.length);
        let tempPermission = {id: p[i].id, label: newName};
        formatted.push(tempPermission);
    }
    return formatted;
}

// Check if a permission is selected
function isPermissionSelected(permissionId) {
    return selectedPermissions.value.includes(permissionId);
}

// Toggle permission selection
function togglePermission(permissionId, checked) {
    if (checked) {
        if (!selectedPermissions.value.includes(permissionId)) {
            selectedPermissions.value.push(permissionId);
        }
    } else {
        selectedPermissions.value = selectedPermissions.value.filter(
            id => id !== permissionId
        );
    }
}

// Fetch roles from backend API
function loadRoles() {
    apiFetch("/roles")
        .then(async (response) => {
            if (!response.ok) {
                toastRef.value?.showToast(
                    "There was an issue getting all roles",
                    ToastTypes.FAIL,
                );
                throw new Error("Failed to fetch all roles");
            }
            roles.value = await response.json();
            // Format permissions
            for (let i = 0; i < roles.value.length; i++) {
                if (roles.value[i].permissions != undefined) {
                    roles.value[i].permissions = formatPermissions(roles.value[i].permissions);
                }
            }
        })
        .catch((error) => {
            console.error(error);
        });
}

loadPermissions();
loadRoles();
</script>

<template>
    <Toast ref="toastRef" />

    <Dialog :open="editModal" @update:open="editModal = $event">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>Edit Role</DialogTitle>
                <DialogDescription>
                    Click 'Confirm' to save changes.
                </DialogDescription>
            </DialogHeader>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="name-input" class="text-right">Name</Label>
                    <Input id="name-input" type="text" v-model="tempRole.name" class="col-span-3" />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="color-input" class="text-right">Color</Label>
                    <Input id="color-input" type="text" v-model="tempRole.color" class="col-span-3 h-10 w-full" />
                </div>
                <div class="space-y-2">
                    <Label class="text-right pt-1">Permissions</Label>
                    <div class="space-y-2 max-h-40">
                        <div 
                            v-for="permission in permissions" 
                            :key="permission.id"
                            class="flex items-center justify-between space-x-2 py-1 pl-2"
                        >
                            <Label 
                                :for="`permission-${permission.id}`" 
                                class="text-sm font-normal cursor-pointer"
                            >
                                {{ permission.label }}
                            </Label>
                            <Checkbox 
                                :id="`permission-${permission.id}`"
                                :model-value="isPermissionSelected(permission.id)"
                                @update:model-value="togglePermission(permission.id, $event)"
                            />
                        </div>
                    </div>
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

    <Dialog :open="createModal" @update:open="createModal = $event">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>Create New Role</DialogTitle>
                <DialogDescription>
                    Fill in the details for the new role.
                </DialogDescription>
            </DialogHeader>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="create-name-input" class="text-right">Name</Label>
                    <Input id="create-name-input" type="text" v-model="tempRole.name" class="col-span-3" />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="create-color-input" class="text-right">Color</Label>
                    <Input id="create-color-input" type="text" v-model="tempRole.color" class="col-span-3 h-10 w-full" />
                </div>
            </div>
            <DialogFooter>
                <DialogClose as-child>
                    <Button type="button" variant="outline">Cancel</Button>
                </DialogClose>
                <Button type="submit" @click="submitCreate">Confirm</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog :open="deleteModal" @update:open="deleteModal = $event">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>Delete Role?</DialogTitle>
                <DialogDescription>
                    This action cannot be undone. This will permanently
                    delete the role.
                </DialogDescription>
            </DialogHeader>
            <DialogFooter class="flex justify-between space-x-4">
                <DialogClose as-child>
                    <Button type="button" variant="outline">Cancel</Button>
                </DialogClose>
                <Button type="submit" @click="submitDelete">Confirm</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <!-- list header -->
    <div class="flex flex-row justify-between items-center mb-6">
        <h1 class="text-2xl font-semibold text-gray-800">Roles</h1>
        <div class="flex flex-row pb-2 items-center space-x-3">
            <p class="text-gray-600">{{ roles.length }} Roles</p>
            <input placeholder="Search"
                class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
                type="text">
        </div>
    </div>

    <ul class="grid grid-cols-1 lg:grid-cols-2 auto-rows-fr gap-4 items-start">
        <li v-for="r in roles" :key="r.id" class="h-full">
            <Card class="h-full">
                <CardHeader class="flex items-center justify-between">
                    <CardTitle class="text-lg">{{ r.name }}</CardTitle>
                    <Badge :style="{ backgroundColor: '#' + r.color }">{{ '#' + r.color }}</Badge>
                    <!-- <Badge>{{ '#' + r.color }}</Badge> -->
                </CardHeader>
                <CardFooter class="flex justify-between gap-x-4">
                    <Button class="flex-1" variant="outline" @click="showEdit(r.id)">Edit</Button>
                    <Button class="flex-1" @click="showDelete(r.id)">Delete</Button>
                </CardFooter>
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

<script setup>
import { ref } from 'vue';
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardFooter, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import CardDotted from './CardDotted.vue';
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
import { apiFetch } from '@/api/request';
import Toast, { ToastTypes } from './Toast.vue';

const settings = ref([]);

const toastRef = ref(null);
const editModal = ref(false);
const createModal = ref(false);
const deleteModal = ref(false);



// Show the edit role dialog
function showEdit(id) {
    console.log(id)
    const settingToEdit = settings.value.find(item => { return item.id == id });
    tempRole.value = JSON.parse(JSON.stringify(settingToEdit));
    originalRole = JSON.parse(JSON.stringify(settingToEdit));
    editModal.value = true;
}

// Send PATCH request to backend API
function submitEdit(key) {
    let value = document.getElementById(`setting-${key}`).value
    console.log(value);
 
    apiFetch(`/settings/${key}`, "PATCH", value)
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
                loadSettings();
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

// Fetch roles from backend API
function loadSettings() {
    apiFetch("/settings")
        .then(async (response) => {
            
            if (!response.ok) {
                toastRef.value?.showToast(
                    "There was an issue getting all settings",
                    ToastTypes.FAIL,
                );
                throw new Error("Failed to fetch all settings");
            }

            settings.value = await response.json();
            console.log(settings.value)
        })
        .catch((error) => {
            console.error(error);
        });
}

loadSettings();
</script>

<template>
    <Toast ref="toastRef" />


    <!-- list header -->
    <div class="flex flex-row justify-between items-center mb-6">
        <h1 class="text-2xl font-semibold text-gray-800">Settings</h1>
        <div class="flex flex-row pb-2 items-center space-x-3">
            <p class="text-gray-600">{{ settings.length }} Setting(s)</p>
            <input placeholder="Search"
                class="p-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent w-64"
                type="text">
        </div>
    </div>

    <ul class="grid grid-cols-1 lg:grid-cols-1 auto-rows-fr gap-4 items-start">
        <li v-for="s in settings" :key="s.skey" class="h-full">
            <Card class="h-full">
                <CardHeader class="flex items-center justify-between">
                    <CardTitle class="text-lg">{{ s.skey }}</CardTitle>
                    <Input :id="`setting-${s.skey}`" :defaultValue="s.svalue"></Input>
                </CardHeader>
                <CardFooter class="flex justify-between gap-x-4">
                    <Button class="flex-1" variant="outline" @click="submitEdit(s.skey)">Push Change</Button>
                </CardFooter>
            </Card>
        </li>
    </ul> 
</template>

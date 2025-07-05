<script setup>
import { ref } from 'vue';

import { useUserStore } from '@/stores/user';
import { apiFetch } from '@/api/request';
import { Button } from '@/components/ui/button';
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

import ProfileArea from './ProfileArea.vue';
import Toast, { ToastTypes } from './Toast.vue';
import DropdownMenu from './ui/dropdown-menu/DropdownMenu.vue';
import DropdownMenuContent from './ui/dropdown-menu/DropdownMenuContent.vue';
import DropdownMenuTrigger from './ui/dropdown-menu/DropdownMenuTrigger.vue';
import Sidebar from './ui/sidebar/Sidebar.vue';
import SidebarContent from './ui/sidebar/SidebarContent.vue';
import SidebarFooter from './ui/sidebar/SidebarFooter.vue';
import SidebarGroup from './ui/sidebar/SidebarGroup.vue';
import SidebarGroupAction from './ui/sidebar/SidebarGroupAction.vue';
import SidebarGroupContent from './ui/sidebar/SidebarGroupContent.vue';
import SidebarGroupLabel from './ui/sidebar/SidebarGroupLabel.vue';
import SidebarHeader from './ui/sidebar/SidebarHeader.vue';
import SidebarMenu from './ui/sidebar/SidebarMenu.vue';
import SidebarMenuButton from './ui/sidebar/SidebarMenuButton.vue';
import SidebarMenuItem from './ui/sidebar/SidebarMenuItem.vue';

const user = useUserStore()
const toastRef = ref(null);
const dashboards = ref([]);
const createDashboardDialog = ref(false);
const newDashboard = ref({});

function addDashboard() {
    const route = '/dashboard';
    newDashboard.value.content = 'placeholder content';
    apiFetch(route, 'POST', JSON.stringify(newDashboard.value))
    .then(async (response) => {
        if (!response.ok) {
            toastRef.value?.showToast(
                "There was an issue creating the dashboard",
                ToastTypes.FAIL,
            );
            throw new Error("Unable to create dashboard")
        } else {
            createDashboardDialog.value = false;
            toastRef.value?.showToast(
                "Dashboard created successfully",
                ToastTypes.SUCCESS,
            )
            loadDashboards();
        }
    })
    .catch((error) => {
        console.log(error);
    })
}

function showCreateDashboardDialog() {
    createDashboardDialog.value = true;
}

function loadDashboards() {
    const route = '/dashboards'
    apiFetch(route)
    .then(async (response) => {
            if (!response.ok) {
                throw new Error('unable to fetch dashboards');
            } else {
                dashboards.value = await response.json();
            }
        }
    )
    .catch((error) => {
        console.error(error);
        throw new Error(error);
    });
}
loadDashboards();
</script>

<template>
    <Toast ref="toastRef" />
    
    <Sidebar variant="inset">
        <SidebarHeader>
            <img src="@/assets/StrataFullx256.png" alt="STRATA" class="mx-auto w-5/6 pt-2"></img>
        </SidebarHeader>
        <SidebarContent>
            <SidebarGroup>
                <SidebarGroupLabel>Administration</SidebarGroupLabel>
                <SidebarGroupContent>
                    <SidebarMenu>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <RouterLink to="/admin/users">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                        <g fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2">
                                            <path d="M18 21a8 8 0 0 0-16 0" />
                                            <circle cx="10" cy="8" r="5" />
                                            <path d="M22 20c0-3.37-2-6.5-4-8a5 5 0 0 0-.45-8.3" />
                                        </g>
                                    </svg>
                                    Members
                                </RouterLink>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <RouterLink to="/admin/roles">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                        <g fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2">
                                            <path
                                                d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z" />
                                            <path d="m9 12l2 2l4-4" />
                                        </g>
                                    </svg>
                                    Roles
                                </RouterLink>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <RouterLink to="/admin/roles">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                        <g fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2">
                                            <path
                                                d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2" />
                                            <circle cx="12" cy="12" r="3" />
                                        </g>
                                    </svg>
                                    Settings
                                </RouterLink>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
            <SidebarGroup>
                <SidebarGroupLabel>Workbench</SidebarGroupLabel>
                <SidebarGroupContent>
                    <SidebarMenu>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <RouterLink to="/query/list">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                        <path fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2"
                                            d="m18 16l4-4l-4-4M6 8l-4 4l4 4m8.5-12l-5 16" />
                                    </svg>
                                    Queries
                                </RouterLink>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <RouterLink to="/">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                                        <g fill="none" stroke="currentColor" stroke-linecap="round"
                                            stroke-linejoin="round" stroke-width="2">
                                            <path d="M3 3v16a2 2 0 0 0 2 2h16" />
                                            <path d="m19 9l-5 5l-4-4l-3 3" />
                                        </g>
                                    </svg>
                                    Graphs
                                </RouterLink>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
            <SidebarGroup>
                <SidebarGroupLabel>Dashboards</SidebarGroupLabel>
                <SidebarGroupAction @click="showCreateDashboardDialog" title="Add Dashboards">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                        <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                            stroke-width="2" d="M5 12h14m-7-7v14" />
                    </svg> <span class="sr-only">Add Dashboards</span>
                </SidebarGroupAction>
                <SidebarContent v-for="d in dashboards">
                    <SidebarMenuItem>
                        <SidebarMenuButton>
                            <RouterLink :to="`/dashboard/${d.id}`">
                                {{ d.title }}
                            </RouterLink>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarContent>
            </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
            <SidebarMenu>
                <SidebarMenuItem>
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <SidebarMenuButton class="h-fit hover:cursor-pointer">
                                <span
                                    class="flex items-center justify-center bg-neutral-200 text-black rounded-full w-8 h-8 p-2 text-lg">
                                    {{ user.username[0].toUpperCase() }}
                                </span> <span>{{ user.username }}</span>
                                <svg class="ml-auto rotate-90" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                    viewBox="0 0 24 24">
                                    <path fill="none" stroke="currentColor" stroke-linecap="round"
                                        stroke-linejoin="round" stroke-width="2" d="m16 18l6-6l-6-6M8 6l-6 6l6 6" />
                                </svg>
                            </SidebarMenuButton>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent side="top" class="w-60">
                            <ProfileArea></ProfileArea>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </SidebarMenuItem>
            </SidebarMenu>
        </SidebarFooter>
    </Sidebar>

    <Dialog :open="createDashboardDialog" @update:open="createDashboardDialog = $event">
        <DialogContent class="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Create New Dashboard</DialogTitle>
                    <DialogDescription>
                        Fill in the details for the new dashboard.
                    </DialogDescription>
                </DialogHeader>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="name-input" class="text-right">Name</Label>
                    <Input id="name-input" type="text" v-model="newDashboard.title" class="col-span-3" />
                </div>
                <DialogFooter>
                    <DialogClose as-child>
                        <Button type="button" variant="outline">Cancel</Button>
                    </DialogClose>
                    <Button type="submit" @click="addDashboard">Confirm</Button>
                </DialogFooter>
        </DialogContent>
    </Dialog> 
</template>
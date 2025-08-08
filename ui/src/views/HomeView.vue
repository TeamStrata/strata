<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import Navbar from '@/components/Navbar.vue';
import Separator from '@/components/ui/separator/Separator.vue';
import SidebarInset from '@/components/ui/sidebar/SidebarInset.vue';
import SidebarProvider from '@/components/ui/sidebar/SidebarProvider.vue';
import SidebarTrigger from '@/components/ui/sidebar/SidebarTrigger.vue';
import { usePageInfoStore } from '@/stores/pageInfo';
import { Edit, Settings } from 'lucide-vue-next';


const route = useRoute()

const pageInfoStore = usePageInfoStore();
const pageInfo = computed(() => {
  // return route.meta
  if (route.meta.title) {
    return {
      title: route.meta.title,
      description: route.meta.description,
      configurable: route.meta.configurable,
      ignorePadding: route.meta.ignorePadding || false,
    };
  } else {
    return {
      title: pageInfoStore.pageInfo.title.value,
      description: pageInfoStore.pageInfo.description.value,
      configurable: pageInfoStore.pageInfo.configurable.value,
      ignorePadding: route.meta.ignorePadding || false,
    };
  }
})

import { useEventBus } from '@/stores/eventBus'
const bus = useEventBus()



</script>

<template>
  <!-- main separator -->
  <div class="h-screen w-screen overflow-clip">
    <!-- NEW Sidebar -->
    <SidebarProvider>
      <Navbar></Navbar>
      <SidebarInset>
        <!-- page header -->
        <div class="flex justify-between items-center px-3 py-2">
          <div class="flex items-end gap-5">
            <SidebarTrigger class="md:hidden"></SidebarTrigger>
            <h1 class="scroll-m-20 text-2xl font-semibold tracking-tight">{{ pageInfo.title }}</h1>
            <p class="leading-7">{{ pageInfo.description }}</p>
          </div>
          <!-- TODO: make this a real dropdown menu and button interaction and stuff -->
          <!-- Hidden options button that only shows on pages with the "configurable" property (for dashboard options) -->
          <div class="flex items-center gap-3 *:hover:text-neutral-600 *:cursor-pointer" :class="pageInfo.configurable ? 'block' : 'hidden'">
            <Edit @click="bus.emit('goEdit')"></Edit>
            <Settings @click="bus.emit('goSettings')"></Settings>
          </div>
        </div>
        <Separator></Separator>
        <!-- main content area -->
        <div class="h-full w-full" :class="pageInfo.ignorePadding ? 'p-0' : 'p-5'">
          <RouterView :key="$route.path"></RouterView>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>



<!-- <template>
  <header>
    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
        <RouterLink to="/login">login</RouterLink>
        <RouterLink to="/admin">Admin</RouterLink>
      </nav>
    </div>
  </header>

  <p>Huge content placeholder</p>
</template> -->

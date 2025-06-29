<script setup>
import Navbar from '@/components/Navbar.vue';
import ProfileArea from '@/components/ProfileArea.vue';
import Separator from '@/components/ui/separator/Separator.vue';
import SidebarInset from '@/components/ui/sidebar/SidebarInset.vue';
import SidebarProvider from '@/components/ui/sidebar/SidebarProvider.vue';
import { watch, ref, computed } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

const route = useRoute()

//expand this in future to make it handle dynamic pages like
const pageInfo = computed(() => { return route.meta })

</script>

<template>
  <!-- main separator -->
  <div>
    <!-- NEW Sidebar -->
    <SidebarProvider>
      <Navbar></Navbar>
      <SidebarInset>
        <!-- page header -->
        <div class="flex justify-between items-center px-5 py-2">
          <div class="flex items-end gap-5">
            <h1 class="scroll-m-20 text-2xl font-semibold tracking-tight">{{ pageInfo.title }}</h1>
            <p class="leading-7">{{ pageInfo.description }}</p>
          </div>
          <!-- TODO: make this a real dropdown menu and button interaction and stuff -->
          <!-- Hidden options button that only shows on pages with the "configurable" property (for dashboard options) -->
          <div :class="pageInfo.configurable ? 'block' : 'hidden'">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
              <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 12h16M4 18h16M4 6h16" />
            </svg>
          </div>
        </div>
        <Separator></Separator>
        <!-- main content area -->
        <div class="p-5">
          <RouterView></RouterView>
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

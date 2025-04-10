<script setup>
import ProfileArea from '@/components/ProfileArea.vue';
import { RouterLink } from 'vue-router';

</script>

<template>
  <!-- main separator -->
  <div class="flex flex-row h-screen">
    <!-- sidebar -->
    <div class="w-xs h-full bg-neutral-100 flex flex-col p-7 border-r-1 border-r-gray-200">
      <!-- logo -->
      <div class="mx-auto p-4">
        Logo Placeholder
      </div>
      <!-- admin controls -->
      <div>
        <RouterLink to="/admin">Admin button</RouterLink>
      </div>
      <div>
        <RouterLink to="">Queries</RouterLink>
      </div>
      <!-- navigation stuff -->
      <div class="flex flex-col">
        Spaces
        <RouterLink to="/">Link</RouterLink>
      </div>
      <!-- header bar -->
      <div class="mt-auto">
        <ProfileArea></ProfileArea>
      </div>
    </div>
    <!-- right side -->
    <div class="flex flex-col p-10">
      <!-- main content area -->
      <div>
        THERE SHOULD BE CONTENT HERE BUT WE HAVE NO CONTENT TO PUT HERE YET
      </div>
    </div>
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

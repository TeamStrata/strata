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
        <img src="@/assets/StrataFullx256.png" alt="STRATA"></img>
      </div>
      <div class="w-full h-0.5 bg-neutral-200 my-3"></div>

      <!-- admin controls -->
      <div>
        <RouterLink class="flex flex-row *:mr-2 mx-4 py-2" to="/admin">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 512 512"><path fill="currentColor" d="m466.5 83.7l-192-80a48.15 48.15 0 0 0-36.9 0l-192 80C27.7 91.1 16 108.6 16 128c0 198.5 114.5 335.7 221.5 380.3c11.8 4.9 25.1 4.9 36.9 0C360.1 472.6 496 349.3 496 128c0-19.4-11.7-36.9-29.5-44.3M256.1 446.3l-.1-381l175.9 73.3c-3.3 151.4-82.1 261.1-175.8 307.7"/></svg>
          Administration
        </RouterLink>
        <RouterLink class="flex flex-row *:mr-2 mx-4 py-2" to="/admin">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 48 48">
            <defs>
              <mask id="ipSLog0">
                <g fill="none" stroke-linejoin="round" stroke-width="4">
                  <path fill="#fff" stroke="#fff" d="M13 10h28v34H13z" />
                  <path stroke="#fff" stroke-linecap="round" d="M35 10V4H8a1 1 0 0 0-1 1v33h6" />
                  <path stroke="#000" stroke-linecap="round" d="M21 22h12m-12 8h12" />
                </g>
              </mask>
            </defs>
            <path fill="currentColor" d="M0 0h48v48H0z" mask="url(#ipSLog0)" />
          </svg>
          Audit Logs
        </RouterLink>

        <div class="w-full h-0.5 bg-neutral-200 my-3"></div>
      </div>
      <!-- navigation stuff -->
      <div>
        <RouterLink class="flex flex-row *:mr-2 mx-4 py-2" to="/">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <path fill="currentColor"
              d="M10.995 4.68v3.88A2.44 2.44 0 0 1 8.545 11h-3.86a2.38 2.38 0 0 1-1.72-.72a2.4 2.4 0 0 1-.71-1.72V4.69a2.44 2.44 0 0 1 2.43-2.44h3.87a2.42 2.42 0 0 1 1.72.72a2.4 2.4 0 0 1 .72 1.71m10.75.01v3.87a2.46 2.46 0 0 1-2.43 2.44h-3.88a2.5 2.5 0 0 1-1.73-.71a2.44 2.44 0 0 1-.71-1.73V4.69a2.4 2.4 0 0 1 .72-1.72a2.42 2.42 0 0 1 1.72-.72h3.87a2.46 2.46 0 0 1 2.44 2.44m0 10.75v3.87a2.46 2.46 0 0 1-2.43 2.44h-3.88a2.5 2.5 0 0 1-1.75-.69a2.42 2.42 0 0 1-.71-1.73v-3.87a2.4 2.4 0 0 1 .72-1.72a2.42 2.42 0 0 1 1.72-.72h3.87a2.46 2.46 0 0 1 2.44 2.44zm-10.75.01v3.87a2.46 2.46 0 0 1-2.45 2.43h-3.86a2.42 2.42 0 0 1-2.43-2.43v-3.87A2.46 2.46 0 0 1 4.685 13h3.87a2.5 2.5 0 0 1 1.73.72a2.45 2.45 0 0 1 .71 1.73" />
          </svg>
          Home
        </RouterLink>
        <RouterLink class="flex flex-row *:mr-2 mx-4 py-2" to="/query/list">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <path fill="currentColor"
              d="M18.68 12.32a4.49 4.49 0 0 0-6.36.01a4.49 4.49 0 0 0 0 6.36a4.51 4.51 0 0 0 5.57.63L21 22.39L22.39 21l-3.09-3.11c1.13-1.77.87-4.09-.62-5.57m-1.41 4.95c-.98.98-2.56.97-3.54 0c-.97-.98-.97-2.56.01-3.54c.97-.97 2.55-.97 3.53 0c.97.98.97 2.56 0 3.54M10.9 20.1a6.5 6.5 0 0 1-1.48-2.32C6.27 17.25 4 15.76 4 14v3c0 2.21 3.58 4 8 4c-.4-.26-.77-.56-1.1-.9M4 9v3c0 1.68 2.07 3.12 5 3.7v-.2c0-.93.2-1.85.58-2.69C6.34 12.3 4 10.79 4 9m8-6C7.58 3 4 4.79 4 7c0 2 3 3.68 6.85 4h.05c1.2-1.26 2.86-2 4.6-2c.91 0 1.81.19 2.64.56A3.22 3.22 0 0 0 20 7c0-2.21-3.58-4-8-4" />
          </svg>
          Queries
        </RouterLink>
      </div>
      <div class="flex flex-col">
        Spaces
        <RouterLink to="/dashboard">Link</RouterLink>
      </div>
      <!-- header bar -->
      <div class="mt-auto">
        <ProfileArea></ProfileArea>
      </div>
    </div>
    <!-- right side -->
    <div class="w-full h-full p-5">
      <!-- main content area -->
      <RouterView></RouterView>
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

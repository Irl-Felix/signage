<!-- App.vue -->
<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppNavbar from './components/AppNavbar.vue'
</script>

<template>
  <v-app>
    <AppNavbar>
      <v-container fluid>
        <RouterView />
      </v-container>
    </AppNavbar>
  </v-app>
</template>

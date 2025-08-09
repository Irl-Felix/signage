<!-- AppNavbar.vue -->
<template>
  <!-- Navigation Drawer -->
  <v-navigation-drawer app permanent v-model:rail="isRail" :width="240">
    <!-- Drawer Header -->
    <div style="height: 56px; display: flex; align-items: center; justify-content: center">
      <v-img src="https://cdn.vuetifyjs.com/images/logos/logo.svg" max-height="40" contain />
    </div>

    <v-divider />

    <!-- Main Navigation -->
    <v-list density="compact" nav>
      <v-list-item
        v-for="link in navLinks"
        :key="link.to"
        :prepend-icon="link.icon"
        :title="link.title"
        :to="link.to"
        link
      />
    </v-list>

    <!-- Push email section to bottom -->
    <template v-slot:append>
      <v-divider />
      <div
        v-if="!isRail"
        style="padding: 12px; font-size: 0.8rem; text-align: center; color: rgba(0, 0, 0, 0.6)"
      >
        sandra_a88@gmail.com
      </div>
      <v-tooltip v-else text="sandra_a88@gmail.com" location="top">
        <template v-slot:activator="{ props }">
          <div v-bind="props" style="padding: 12px; text-align: center">
            <v-icon size="20" color="grey">mdi-email-outline</v-icon>
          </div>
        </template>
      </v-tooltip>
    </template>
  </v-navigation-drawer>

  <!-- App Bar -->
  <v-app-bar app elevation="2" density="comfortable">
    <v-btn icon @click="toggleDrawer">
      <v-icon>{{ isRail ? 'mdi-menu-open' : 'mdi-menu' }}</v-icon>
    </v-btn>

    <v-spacer />

    <v-text-field
      placeholder="Search..."
      prepend-inner-icon="mdi-magnify"
      variant="outlined"
      density="compact"
      hide-details
      style="max-width: 250px"
    />

    <v-btn icon>
      <v-icon>mdi-bell-outline</v-icon>
    </v-btn>

    <v-menu location="bottom end" offset-y>
      <template v-slot:activator="{ props }">
        <v-btn v-bind="props" icon>
          <v-avatar size="36">
            <img src="https://randomuser.me/api/portraits/women/85.jpg" alt="Avatar" />
          </v-avatar>
        </v-btn>
      </template>
      <v-card min-width="250" style="margin-top: 8px">
        <v-list-item title="Sandra Adams" subtitle="Admin" />
        <v-divider />
        <v-list density="compact" nav>
          <v-list-item
            v-for="link in avatarMenuLinks"
            :key="link.to"
            :prepend-icon="link.icon"
            :title="link.title"
            :to="link.to"
            link
          />
        </v-list>
      </v-card>
    </v-menu>
  </v-app-bar>

  <!-- Main Content (slot for pages) -->
  <v-main>
    <slot />
  </v-main>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isRail = ref(true)

const toggleDrawer = () => {
  isRail.value = !isRail.value
}

const navLinks = [
  { icon: 'mdi-account-group', title: 'Users', to: '/users' },
  { icon: 'mdi-shield-account', title: 'Roles & Permissions', to: '/roleperm' },
  { icon: 'mdi-domain', title: 'Tenants', to: '/tenants' },
  { icon: 'mdi-clipboard-list', title: 'Audit Logs', to: '/logs' },
]

const avatarMenuLinks = [
  { icon: 'mdi-account', title: 'Profile', to: '/profile' },
  { icon: 'mdi-cog', title: 'Settings', to: '/settings' },
  { icon: 'mdi-logout', title: 'Logout', to: '/logout' },
]
</script>

<style scoped>
.v-list-item-title {
  font-weight: 500;
}
</style>

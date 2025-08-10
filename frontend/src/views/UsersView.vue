<template>
  <v-row class="mb-4" dense>
    <!-- Total Users -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="secondary" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-account-group</v-icon>
          Total Users
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.totalUsers }}
        </v-card-text>
      </v-card>
    </v-col>
    <!-- Total Active Users -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="primary" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-account-check-outline</v-icon>
          Active Users
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.activeUsers }}
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Pending Users -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="warning" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-account-clock-outline</v-icon>
          Pending Users
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.pendingUsers }}
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Suspended Users -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="error" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-account-cancel-outline</v-icon>
          Suspended Users
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.suspendedUsers }}
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <!-- Filters -->
  <v-row class="mb-2" dense>
    <v-col cols="12" sm="4">
      <v-text-field
        v-model="filters.search"
        label="Search users"
        prepend-inner-icon="mdi-magnify"
        variant="outlined"
        clearable
        hide-details="auto"
        dense
        class="mr-2"
        @keyup.enter="fetchUsers"
        @click:clear="fetchUsers"
      />
    </v-col>
    <v-col cols="12" sm="4">
      <v-select
        v-model="filters.status"
        :items="statusOptions"
        label="Status"
        variant="outlined"
        clearable
        hide-details="auto"
        dense
        class="mr-2"
      />
    </v-col>
    <v-col cols="12" sm="4">
      <v-select
        v-model="filters.role"
        :items="roleOptions"
        label="Role"
        variant="outlined"
        clearable
        hide-details="auto"
        dense
      />
    </v-col>
  </v-row>

  <!-- Data Table -->
  <v-data-table-server
    :headers="headers"
    :items="displayedUsers"
    :loading="loading"
    :items-length="displayedUsers.length"
    density="comfortable"
    class="elevation-1"
  >
    <!-- Avatar -->
    <template v-slot:[`item.avatar`]="{ item }">
      <v-avatar size="32" class="mr-2">
        <img v-if="item.avatar" :src="item.avatar" alt="User Avatar" />
        <span v-else>{{ getInitials(item.name) }}</span>
      </v-avatar>
    </template>

    <!-- User column: Name (bold) + Email (small) -->

    <template v-slot:[`item.name`]="{ item }">
      <div>
        <div class="font-weight-bold text-body-1">{{ item.name }}</div>
        <div class="text-body-2 text-medium-emphasis">{{ item.email }}</div>
      </div>
    </template>

    <!-- Session column -->

    <template v-slot:[`item.session`]="{ item }">
      <v-icon :color="item.session ? 'success' : 'error'" size="small">
        {{ item.session ? "mdi-check-circle" : "mdi-close-circle" }}
      </v-icon>
    </template>

    <!-- Actions: icons displayed horizontally -->
    <template v-slot:[`item.actions`]="{ item }">
      <v-btn
        icon
        variant="text"
        size="small"
        color="primary"
        @click="viewUser(item)"
      >
        <v-icon>mdi-eye</v-icon>
      </v-btn>

      <v-btn
        icon
        variant="text"
        size="small"
        color="warning"
        @click="editRoles(item)"
      >
        <v-icon>mdi-account-cog</v-icon>
      </v-btn>

      <v-btn
        icon
        variant="text"
        size="small"
        color="error"
        @click="deleteUser(item)"
      >
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </template>
  </v-data-table-server>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";

const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  pendingUsers: 0,
  suspendedUsers: 0,
});

const loading = ref(false);
interface User {
  avatar?: string; // derived from profile_url
  name: string;
  email: string;
  business: string;
  global_role?: string;
  status?: string;
  session?: boolean;
  created_at?: string;
  profile_url?: string;
}

const users = ref<User[]>([]);

// Filter state
const filters = ref({
  search: "",
  status: null as string | null,
  role: null as string | null,
});

const statusOptions = [
  { title: "All Statuses", value: null },
  { title: "Active", value: "active" },
  { title: "Pending", value: "pending" },
  { title: "Suspended", value: "suspended" },
];
// Update roleOptions to match backend role codes
const roleOptions = [
  { title: "All Roles", value: null },
  { title: "System Admin", value: "SYS_ADMIN" },
  { title: "Business Owner", value: "BUSINESS_OWNER" },
  { title: "Branch Manager", value: "BRANCH_MANAGER" },
  { title: "Content Editor", value: "CONTENT_EDITOR" },
];

const headers = [
  { title: "", key: "avatar", sortable: false },
  { title: "User", key: "name", sortable: true },
  { title: "Session", key: "session", sortable: false },
  { title: "Business", key: "business", sortable: true },
  { title: "Global Role", key: "global_role", sortable: true },
  { title: "Status", key: "status", sortable: true },
  { title: "Created", key: "created_at", sortable: true },
  { title: "Actions", key: "actions", sortable: false },
];

// API base uses relative path; Vite dev server and Nginx will proxy /api to backend
function joinUrl(base: string, path: string) {
  if (!base.endsWith("/") && !path.startsWith("/")) return `${base}/${path}`;
  if (base.endsWith("/") && path.startsWith("/"))
    return `${base}${path.slice(1)}`;
  return `${base}${path}`;
}
const API_BASE = "/api";

const API_URL = joinUrl(API_BASE, "/admin/users/test");
const STATS_URL = joinUrl(API_BASE, "/admin/users/stats");
const abortController = ref<AbortController | null>(null);
// Fetch user stats for dashboard cards
async function fetchUserStats() {
  try {
    const res = await fetch(STATS_URL, {
      headers: { Accept: "application/json" },
    });
    if (!res.ok) throw new Error(`Stats request failed: ${res.status}`);
    const json = await res.json();
    // Defensive: check for expected keys
    stats.value = {
      totalUsers: json.total_users ?? 0,
      activeUsers: json.active_users ?? 0,
      pendingUsers: json.pending_users ?? 0,
      suspendedUsers: json.suspended_users ?? 0,
    };
  } catch (err) {
    console.error("Failed to load user stats", err);
    // Optionally: keep stats as zeroes
  }
}

function buildQueryParams() {
  const params = new URLSearchParams();
  if (filters.value.search) params.append("search", filters.value.search);
  if (filters.value.status) params.append("status", filters.value.status);
  if (filters.value.role) params.append("role", filters.value.role);
  return params.toString();
}

// Debounce helper typed without any
function debounce<F extends (...args: unknown[]) => void>(fn: F, delay = 300) {
  let t: ReturnType<typeof setTimeout> | undefined;
  return (...args: Parameters<F>) => {
    if (t) clearTimeout(t);
    t = setTimeout(() => fn(...args), delay);
  };
}

const debouncedFetchUsers = debounce(() => fetchUsers(), 300);

async function fetchUsers() {
  try {
    loading.value = true;
    if (abortController.value) abortController.value.abort();
    abortController.value = new AbortController();

    let url = API_URL;
    const query = buildQueryParams();
    if (query) url += `?${query}`;

    console.log("[Users] Fetch:", url);

    const res = await fetch(url, {
      method: "GET",
      headers: { Accept: "application/json" },
      signal: abortController.value.signal,
    });
    if (!res.ok) throw new Error(`Request failed: ${res.status}`);

    const json = (await res.json()) as unknown;

    // Expecting { users: User[] }
    type UsersResponse = { users: User[] };
    const data = json as UsersResponse;
    if (data && Array.isArray(data.users)) {
      users.value = data.users.map((u) => ({
        ...u,
        avatar: u.profile_url || undefined,
      }));
    } else {
      users.value = [];
      console.warn("Unexpected API response shape", json);
    }
  } catch (err: unknown) {
    const isAbort = err instanceof DOMException && err.name === "AbortError";
    if (!isAbort) {
      console.error("Failed to load users", err);
    }
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchUsers();
  fetchUserStats();
});

// Watchers
watch(
  () => filters.value.search,
  () => debouncedFetchUsers()
);
watch(() => [filters.value.status, filters.value.role], fetchUsers);

const displayedUsers = computed(() => {
  const s = (filters.value.search || "").trim().toLowerCase();
  const st = (filters.value.status || "").toLowerCase();
  const rl = (filters.value.role || "").toLowerCase();

  return users.value.filter((u) => {
    // search across name, email, business
    if (s) {
      const hay =
        `${u.name || ""} ${u.email || ""} ${u.business || ""}`.toLowerCase();
      if (!hay.includes(s)) return false;
    }
    // status match
    if (st) {
      if ((u.status || "").toLowerCase() !== st) return false;
    }
    // role match (codes like SYS_ADMIN)
    if (rl) {
      if ((u.global_role || "").toLowerCase() !== rl) return false;
    }
    return true;
  });
});

function getInitials(name: string) {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase();
}

function viewUser(item: User) {
  console.log("View user:", item);
}

function editRoles(item: User) {
  console.log("Edit roles for:", item);
}

function deleteUser(item: User) {
  console.log("Delete user:", item);
}
</script>

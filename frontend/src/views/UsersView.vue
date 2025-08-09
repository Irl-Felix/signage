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

  <!-- Search & Filters -->
  <v-row class="mb-4" dense>
    <!-- Search Bar -->
    <v-col cols="12" md="4">
      <v-text-field
        v-model="filters.search"
        placeholder="Search users..."
        prepend-inner-icon="mdi-magnify"
        density="compact"
        variant="outlined"
        hide-details
        clearable
      />
    </v-col>

    <!-- Status Filter -->
    <v-col cols="12" sm="6" md="4">
      <v-select
        v-model="filters.status"
        :items="statusOptions"
        label="Status"
        clearable
        density="compact"
        variant="outlined"
        hide-details
      />
    </v-col>

    <!-- Role Filter -->
    <v-col cols="12" sm="6" md="4">
      <v-select
        v-model="filters.role"
        :items="roleOptions"
        label="Global Role"
        clearable
        density="compact"
        variant="outlined"
        hide-details
      />
    </v-col>
  </v-row>

  <!-- Data Table -->
  <v-data-table-server
    v-model:page="pagination.page"
    v-model:items-per-page="pagination.itemsPerPage"
    v-model:sort-by="pagination.sortBy"
    :headers="headers"
    :items="users"
    :loading="loading"
    :server-items-length="totalUsers"
    :items-length="totalUsers"
    :items-per-page-options="[5, 10, 25, 50]"
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
import { ref, onMounted, watch } from "vue";

const filters = ref({
  search: "",
  status: null,
  role: null,
});

const stats = ref({
  activeUsers: 120,
  pendingUsers: 15,
  suspendedUsers: 5,
  totalUsers: 140,
});

const statusOptions = [
  { title: "Active", value: "active" },
  { title: "Pending", value: "pending" },
  { title: "Suspended", value: "suspended" },
];
const roleOptions = [
  { title: "System Admin", value: "SYS_ADMIN" },
  { title: "Support Admin", value: "SUPPORT_ADMIN" },
  { title: "Viewer Global", value: "VIEWER_GLOBAL" },
];

const loading = ref(false);
interface User {
  avatar?: string;
  name: string;
  email: string;
  business: string;
  global_role?: string;
  status?: string;
  session?: boolean;
  created_at?: string;
}

const users = ref<User[]>([]);
const totalUsers = ref(0);

const pagination = ref({
  page: 1,
  itemsPerPage: 10,
  sortBy: [{ key: "created_at", order: "desc" as const }],
});

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

// API integration via Vite env variable (base), fixed path
function joinUrl(base: string, path: string) {
  if (!base.endsWith("/") && !path.startsWith("/")) return `${base}/${path}`;
  if (base.endsWith("/") && path.startsWith("/"))
    return `${base}${path.slice(1)}`;
  return `${base}${path}`;
}
const API_URL = joinUrl(import.meta.env.VITE_BACKEND_URL, "/admin/users");
const abortController = ref<AbortController | null>(null);

function debounce<F extends (...args: unknown[]) => void>(fn: F, wait = 300) {
  let t: ReturnType<typeof setTimeout> | undefined;
  return (...args: Parameters<F>) => {
    if (t) clearTimeout(t);
    t = setTimeout(() => fn(...args), wait);
  };
}

function buildQueryParams() {
  const params = new URLSearchParams();
  params.set("page", String(pagination.value.page));
  params.set("per_page", String(pagination.value.itemsPerPage));

  const sort = pagination.value.sortBy?.[0];
  if (sort?.key) {
    params.set("sort", String(sort.key));
    params.set("order", sort.order ?? "asc");
  }

  if (filters.value.search) params.set("search", filters.value.search);
  if (filters.value.status) params.set("status", String(filters.value.status));
  if (filters.value.role) params.set("role", String(filters.value.role));

  return params;
}

async function fetchUsers() {
  try {
    loading.value = true;
    // Cancel any in-flight request
    if (abortController.value) abortController.value.abort();
    abortController.value = new AbortController();

    const url = new URL(API_URL);
    url.search = buildQueryParams().toString();

    const res = await fetch(url.toString(), {
      method: "GET",
      headers: { Accept: "application/json" },
      signal: abortController.value.signal,
    });
    if (!res.ok) throw new Error(`Request failed: ${res.status}`);

    const json = (await res.json()) as unknown;

    function getTotalLike(obj: Record<string, unknown>, itemsLen: number) {
      const meta =
        typeof obj.meta === "object" && obj.meta !== null
          ? (obj.meta as Record<string, unknown>)
          : undefined;
      const candidates = [obj.total, obj.count, meta?.total];
      for (const c of candidates) {
        if (typeof c === "number") return c;
        if (
          typeof c === "string" &&
          c.trim() !== "" &&
          !Number.isNaN(Number(c))
        )
          return Number(c);
      }
      return itemsLen;
    }

    if (Array.isArray(json)) {
      users.value = json as User[];
      totalUsers.value = json.length;
    } else if (typeof json === "object" && json !== null) {
      const obj = json as Record<string, unknown>;

      if (Array.isArray(obj.data)) {
        users.value = obj.data as User[];
        totalUsers.value = getTotalLike(obj, (obj.data as unknown[]).length);
      } else if (Array.isArray(obj.items)) {
        users.value = obj.items as User[];
        totalUsers.value = getTotalLike(obj, (obj.items as unknown[]).length);
      } else if (Array.isArray(obj.users)) {
        users.value = obj.users as User[];
        totalUsers.value = getTotalLike(obj, (obj.users as unknown[]).length);
      } else {
        users.value = [];
        totalUsers.value = 0;
        console.warn("Unexpected API response shape", json);
      }
    } else {
      users.value = [];
      totalUsers.value = 0;
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

const debouncedFetchUsers = debounce(fetchUsers, 400);

// Watch pagination and sorting
watch(
  [
    () => pagination.value.page,
    () => pagination.value.itemsPerPage,
    () => pagination.value.sortBy,
  ],
  () => {
    fetchUsers();
  },
  { deep: true }
);

// Watch filters (debounce search, immediate for others)
watch(
  () => filters.value.search,
  () => {
    pagination.value.page = 1;
    debouncedFetchUsers();
  }
);

watch([() => filters.value.status, () => filters.value.role], () => {
  pagination.value.page = 1;
  fetchUsers();
});

onMounted(() => {
  fetchUsers();
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

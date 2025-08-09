<template>
  <!-- 📊 Stats Cards -->
  <v-row class="mb-4" dense>
    <!-- Total Businesses -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="secondary" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-domain</v-icon>
          Total Businesses
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.totalBusinesses }}
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Active Businesses -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="primary" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-domain-check</v-icon>
          Active Businesses
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.activeBusinesses }}
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Pending Businesses -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="warning" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-domain-clock</v-icon>
          Pending Approval
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.pendingBusinesses }}
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Suspended Businesses -->
    <v-col cols="12" sm="6" md="3">
      <v-card color="error" variant="flat" class="pa-4" elevation="2">
        <v-card-title class="text-white d-flex align-center">
          <v-icon start class="mr-2">mdi-domain-remove</v-icon>
          Suspended
        </v-card-title>
        <v-card-text class="text-h5 font-weight-bold text-white">
          {{ stats.suspendedBusinesses }}
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <!-- 🔍 Search & Filters -->
  <v-row class="mb-4" dense>
    <!-- Search Bar -->
    <v-col cols="12" md="4">
      <v-text-field
        v-model="filters.search"
        placeholder="Search businesses..."
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

    <!-- Plan Filter -->
    <v-col cols="12" sm="6" md="4">
      <v-select
        v-model="filters.plan"
        :items="planOptions"
        label="Plan"
        clearable
        density="compact"
        variant="outlined"
        hide-details
      />
    </v-col>
  </v-row>

  <!-- 📋 Data Table -->
  <v-data-table-server
    v-model:page="pagination.page"
    v-model:items-per-page="pagination.itemsPerPage"
    v-model:sort-by="pagination.sortBy"
    :headers="headers"
    :items="businesses"
    :loading="loading"
    :items-length="totalBusinesses"
    :items-per-page-options="[5, 10, 25, 50]"
    density="comfortable"
    class="elevation-1"
  >
    <!-- Owner Column -->
    <template v-slot:[`item.owner`]="{ item }">
      <div class="d-flex align-center">
        <v-avatar size="28" class="mr-2">
          <img v-if="item.owner?.avatar" :src="item.owner.avatar" alt="Owner" />
          <span v-else>{{ getInitials(item.owner?.name) }}</span>
        </v-avatar>
        <div>
          <div>{{ item.owner?.name }}</div>
          <small class="text-grey">{{ item.owner?.email }}</small>
        </div>
      </div>
    </template>

    <!-- Actions Column -->
    <template v-slot:[`item.actions`]="{ item }">
      <v-menu>
        <template #activator="{ props }">
          <v-btn icon v-bind="props">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>
        <v-list density="compact">
          <v-list-item @click="viewBusiness(item)">View Details</v-list-item>
          <v-list-item @click="changePlan(item)">Change Plan</v-list-item>
          <v-list-item @click="suspendBusiness(item)">Suspend</v-list-item>
          <v-list-item @click="deleteBusiness(item)">Delete</v-list-item>
        </v-list>
      </v-menu>
    </template>
  </v-data-table-server>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// Filters
const filters = ref({
  search: '',
  status: null,
  plan: null,
})

// Stats
const stats = ref({
  totalBusinesses: 120,
  activeBusinesses: 100,
  pendingBusinesses: 15,
  suspendedBusinesses: 5,
})

// Status & Plan filter options
const statusOptions = [
  { title: 'Active', value: 'active' },
  { title: 'Pending', value: 'pending' },
  { title: 'Suspended', value: 'suspended' },
]
const planOptions = [
  { title: 'Basic', value: 'basic' },
  { title: 'Pro', value: 'pro' },
  { title: 'Enterprise', value: 'enterprise' },
]

// Loading state
const loading = ref(false)

// Dummy business data
interface Business {
  name: string
  plan: string
  status: string
  owner: {
    avatar?: string
    name: string
    email: string
  }
  branches_count: number
  users_count: number
  created_at: string
}

const businesses = ref<Business[]>([
  {
    name: 'Acme Corporation',
    plan: 'enterprise',
    status: 'active',
    owner: {
      avatar: 'https://randomuser.me/api/portraits/men/45.jpg',
      name: 'Robert Johnson',
      email: 'robert.johnson@acme.com',
    },
    branches_count: 5,
    users_count: 42,
    created_at: '2023-01-15',
  },
  {
    name: 'Beta Solutions',
    plan: 'pro',
    status: 'pending',
    owner: {
      avatar: '',
      name: 'Emily Davis',
      email: 'emily.davis@betasolutions.com',
    },
    branches_count: 2,
    users_count: 15,
    created_at: '2023-03-10',
  },
  {
    name: 'Delta Retail',
    plan: 'basic',
    status: 'suspended',
    owner: {
      avatar: 'https://randomuser.me/api/portraits/women/65.jpg',
      name: 'Sophia Williams',
      email: 'sophia.williams@deltaretail.com',
    },
    branches_count: 1,
    users_count: 6,
    created_at: '2023-05-05',
  },
  {
    name: 'Gamma Tech',
    plan: 'pro',
    status: 'active',
    owner: {
      avatar: '',
      name: 'Michael Brown',
      email: 'michael.brown@gammatech.com',
    },
    branches_count: 3,
    users_count: 28,
    created_at: '2023-06-20',
  },
])

// Total
const totalBusinesses = ref(businesses.value.length)

// Pagination
const pagination = ref({
  page: 1,
  itemsPerPage: 10,
  sortBy: [{ key: 'created_at', order: 'desc' as const }],
})

// Table headers
const headers = [
  { title: 'Business Name', key: 'name', sortable: true },
  { title: 'Plan', key: 'plan', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Owner', key: 'owner', sortable: false },
  { title: 'Branches', key: 'branches_count', sortable: true },
  { title: 'Users', key: 'users_count', sortable: true },
  { title: 'Created', key: 'created_at', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false },
]

function getInitials(name?: string) {
  return name
    ? name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
    : ''
}

function viewBusiness(item: Business) {
  console.log('View business:', item)
}
function changePlan(item: Business) {
  console.log('Change plan:', item)
}
function suspendBusiness(item: Business) {
  console.log('Suspend business:', item)
}
function deleteBusiness(item: Business) {
  console.log('Delete business:', item)
}
</script>

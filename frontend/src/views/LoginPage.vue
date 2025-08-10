<template>
  <div>
    <v-img
      class="mx-auto my-6"
      max-width="228"
      src="https://cdn.vuetifyjs.com/docs/images/logos/vuetify-logo-v3-slim-text-light.svg"
    ></v-img>

    <v-card
      variant="flat"
      class="mx-auto pa-12 pb-8"
      max-width="448"
      rounded="lg"
    >
      <div class="text-subtitle-1 text-medium-emphasis">Account</div>

      <v-text-field
        v-model="email"
        density="compact"
        placeholder="Email address"
        prepend-inner-icon="mdi-email-outline"
        variant="outlined"
        type="email"
        autocomplete="username"
      ></v-text-field>

      <div
        class="text-subtitle-1 text-medium-emphasis d-flex align-center justify-space-between"
      >
        Password

        <a
          class="text-caption text-decoration-none text-blue"
          href="#"
          rel="noopener noreferrer"
          target="_blank"
        >
          Forgot login password?</a
        >
      </div>

      <v-text-field
        v-model="password"
        :append-inner-icon="visible ? 'mdi-eye-off' : 'mdi-eye'"
        :type="visible ? 'text' : 'password'"
        density="compact"
        placeholder="Enter your password"
        prepend-inner-icon="mdi-lock-outline"
        variant="outlined"
        autocomplete="current-password"
        @click:append-inner="visible = !visible"
      ></v-text-field>

      <v-card class="mb-12" color="surface-variant" variant="tonal">
        <v-card-text class="text-medium-emphasis text-caption">
          Warning: After 3 consecutive failed login attempts, you account will
          be temporarily locked for three hours. If you must login now, you can
          also click "Forgot login password?" below to reset the login password.
        </v-card-text>
      </v-card>

      <v-btn
        class="mb-8"
        color="blue"
        size="large"
        variant="tonal"
        block
        :loading="loading"
        @click="onLogin"
      >
        Log In
      </v-btn>
      <div v-if="error" class="text-error text-caption mb-2">{{ error }}</div>
      <div v-if="success" class="text-success text-caption mb-2">
        {{ success }}
      </div>

      <v-card-text class="text-center">
        <router-link class="text-blue text-decoration-none" to="/register">
          Sign up now <v-icon icon="mdi-chevron-right"></v-icon>
        </router-link>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const visible = ref(false);
const email = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");
const success = ref("");
const router = useRouter();

async function onLogin() {
  error.value = "";
  loading.value = true;
  try {
    const res = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: email.value, password: password.value }),
    });
    if (!res.ok) {
      const msg = (await res.json())?.message || `Login failed (${res.status})`;
      error.value = msg;
      return;
    }
    error.value = "";
    success.value = "Login successful! Redirecting...";
    setTimeout(() => {
      router.push("/");
    }, 1200);
  } catch {
    error.value = "Network error";
  } finally {
    loading.value = false;
  }
}
</script>

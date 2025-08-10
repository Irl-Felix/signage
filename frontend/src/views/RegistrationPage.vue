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
      <div class="text-subtitle-1 text-medium-emphasis">
        Create your account
      </div>

      <v-text-field
        v-model="firstName"
        density="compact"
        placeholder="First name"
        prepend-inner-icon="mdi-account-outline"
        variant="outlined"
        autocomplete="given-name"
      ></v-text-field>

      <v-text-field
        v-model="lastName"
        density="compact"
        placeholder="Last name"
        prepend-inner-icon="mdi-account-outline"
        variant="outlined"
        autocomplete="family-name"
      ></v-text-field>

      <v-text-field
        v-model="email"
        density="compact"
        placeholder="Email address"
        prepend-inner-icon="mdi-email-outline"
        variant="outlined"
        type="email"
        autocomplete="email"
      ></v-text-field>

      <v-text-field
        v-model="password"
        :append-inner-icon="visible ? 'mdi-eye-off' : 'mdi-eye'"
        :type="visible ? 'text' : 'password'"
        density="compact"
        placeholder="Password"
        prepend-inner-icon="mdi-lock-outline"
        variant="outlined"
        autocomplete="new-password"
        @click:append-inner="visible = !visible"
      ></v-text-field>

      <v-text-field
        v-model="confirmPassword"
        :append-inner-icon="confirmVisible ? 'mdi-eye-off' : 'mdi-eye'"
        :type="confirmVisible ? 'text' : 'password'"
        density="compact"
        placeholder="Confirm password"
        prepend-inner-icon="mdi-lock-check-outline"
        variant="outlined"
        autocomplete="new-password"
        @click:append-inner="confirmVisible = !confirmVisible"
      ></v-text-field>

      <v-btn
        class="mb-8 mt-4"
        color="blue"
        size="large"
        variant="tonal"
        block
        :loading="loading"
        @click="onRegister"
      >
        Sign Up
      </v-btn>
      <div v-if="error" class="text-error text-caption mb-2">{{ error }}</div>
      <div v-if="success" class="text-success text-caption mb-2">
        {{ success }}
      </div>

      <v-card-text class="text-center">
        <span class="text-caption text-medium-emphasis"
          >Already have an account?</span
        >
        <router-link class="text-blue text-decoration-none ml-1" to="/login">
          Log in <v-icon icon="mdi-chevron-right"></v-icon>
        </router-link>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const visible = ref(false);
const confirmVisible = ref(false);
const firstName = ref("");
const lastName = ref("");
const email = ref("");
const password = ref("");
const confirmPassword = ref("");
const loading = ref(false);
const error = ref("");
const success = ref("");
const router = useRouter();

async function onRegister() {
  error.value = "";
  if (
    !firstName.value ||
    !lastName.value ||
    !email.value ||
    !password.value ||
    !confirmPassword.value
  ) {
    error.value = "All fields are required.";
    return;
  }
  if (password.value !== confirmPassword.value) {
    error.value = "Passwords do not match.";
    return;
  }
  loading.value = true;
  try {
    const res = await fetch("/api/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: email.value,
        password: password.value,
        first_name: firstName.value,
        last_name: lastName.value,
      }),
    });
    if (!res.ok) {
      const msg =
        (await res.json())?.message || `Registration failed (${res.status})`;
      error.value = msg;
      success.value = "";
      return;
    }
    error.value = "";
    success.value = "Registration successful! Redirecting to login...";
    setTimeout(() => {
      router.push("/login");
    }, 1500);
  } catch {
    error.value = "Network error";
    success.value = "";
  } finally {
    loading.value = false;
  }
}
</script>

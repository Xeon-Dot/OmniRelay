<template>
  <AuthShell :title="$t('auth.registerTitle')" :subtitle="$t('auth.registerSubtitle')">
    <form class="auth-form" @submit.prevent="handleRegister">
      <OrTextField
        v-model="username"
        :label="$t('auth.username')"
        leading-icon="person"
        placeholder="admin"
        autocomplete="username"
        required
      />

      <OrTextField
        v-model="email"
        type="email"
        :label="$t('auth.email')"
        leading-icon="email"
        placeholder="admin@example.com"
        autocomplete="email"
        required
      />

      <OrTextField
        v-model="password"
        :type="showPw ? 'text' : 'password'"
        :label="$t('auth.password')"
        leading-icon="lock"
        :hint="$t('auth.passwordRequirements')"
        placeholder="••••••••"
        autocomplete="new-password"
        required
      >
        <template #trailing>
          <OrIconButton
            :icon="showPw ? 'visibility_off' : 'visibility'"
            :label="$t('auth.togglePassword')"
            size="sm"
            tabindex="-1"
            @click="showPw = !showPw"
          />
        </template>
      </OrTextField>

      <OrTextField
        v-model="confirmPassword"
        :type="showPw ? 'text' : 'password'"
        :label="$t('auth.confirmPassword')"
        leading-icon="lock"
        :error="mismatch ? $t('auth.passwordsDontMatch') : ''"
        placeholder="••••••••"
        autocomplete="new-password"
        required
      />

      <div v-if="error" class="auth-error">
        <OrIcon name="error" :size="14" />
        {{ error }}
      </div>

      <OrButton type="submit" :loading="loading" :disabled="mismatch" block>
        {{ $t("auth.createAccount") }}
      </OrButton>
    </form>

    <template #footer>
      {{ $t("auth.haveAccount") }}
      <router-link to="/login">{{ $t("auth.signInLink") }}</router-link>
    </template>
  </AuthShell>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "../stores/auth";
import AuthShell from "../components/AuthShell.vue";

const { t } = useI18n();
const auth = useAuthStore();
const router = useRouter();
const username = ref("");
const email = ref("");
const password = ref("");
const confirmPassword = ref("");
const error = ref("");
const loading = ref(false);
const showPw = ref(false);

const mismatch = computed(
  () =>
    Boolean(password.value) &&
    Boolean(confirmPassword.value) &&
    password.value !== confirmPassword.value,
);

async function handleRegister() {
  if (password.value !== confirmPassword.value) return;
  loading.value = true;
  error.value = "";
  try {
    await auth.register(username.value, email.value, password.value);
    router.push("/");
  } catch (e: any) {
    error.value = e.response?.data?.error || t("auth.registrationFailed");
  } finally {
    loading.value = false;
  }
}
</script>

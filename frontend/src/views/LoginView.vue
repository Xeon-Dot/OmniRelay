<template>
  <AuthShell :title="$t('auth.loginTitle')" :subtitle="$t('auth.loginSubtitle')">
    <form v-if="!twoFactorToken" class="auth-form" @submit.prevent="handleLogin">
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
        placeholder="••••••••"
        autocomplete="current-password"
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

      <div v-if="error" class="auth-error">
        <OrIcon name="error" :size="14" />
        {{ error }}
      </div>

      <OrButton type="submit" :loading="loading" block>
        {{ $t("auth.signIn") }}
      </OrButton>
    </form>

    <form v-else class="auth-form" @submit.prevent="handleTwoFactor">
      <OrTextField
        v-model="twoFactorCode"
        type="text"
        :label="$t('auth.twoFactorCode')"
        leading-icon="security"
        placeholder="123456"
        autocomplete="one-time-code"
        required
      />
      <div v-if="error" class="auth-error">
        <OrIcon name="error" :size="14" />
        {{ error }}
      </div>
      <OrButton type="submit" :loading="loading" block>
        {{ $t("auth.verify") }}
      </OrButton>
    </form>

    <template #footer>
      <template v-if="!twoFactorToken">
        {{ $t("auth.noAccount") }}
        <router-link to="/register">{{ $t("auth.createOne") }}</router-link>
      </template>
      <a v-else href="#" @click.prevent="twoFactorToken = ''">{{ $t("auth.back") }}</a>
    </template>
  </AuthShell>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "../stores/auth";
import AuthShell from "../components/AuthShell.vue";

const { t } = useI18n();
const auth = useAuthStore();
const router = useRouter();
const email = ref("");
const password = ref("");
const error = ref("");
const loading = ref(false);
const showPw = ref(false);
const twoFactorToken = ref("");
const twoFactorCode = ref("");

async function handleLogin() {
  loading.value = true;
  error.value = "";
  try {
    const res = await auth.login(email.value, password.value);
    if (res.requires2fa) {
      twoFactorToken.value = res.twoFactorToken;
    } else {
      router.push("/");
    }
  } catch (e: any) {
    error.value = e.response?.data?.error || t("auth.invalidCredentials");
  } finally {
    loading.value = false;
  }
}

async function handleTwoFactor() {
  loading.value = true;
  error.value = "";
  try {
    await auth.verifyTwoFactor(twoFactorToken.value, twoFactorCode.value.trim());
    router.push("/");
  } catch (e: any) {
    error.value = e.response?.data?.error || t("auth.invalidCredentials");
  } finally {
    loading.value = false;
  }
}
</script>

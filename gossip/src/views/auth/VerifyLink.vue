<script setup lang="ts">
import { useVerifyMagicLink } from "@/composables/auth";
import { onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const { success, verifyMagicLink, asyncStatus, data, error } =
    useVerifyMagicLink();

onMounted(async () => {
    await verifyMagicLink({ token: route.query.token as string });
    if (success) {
        if (data.value.onboard) {
            router.push(`/user/${data.value.userId}/onboarding`);
        } else {
            router.push("/chat");
        }
    }
});
</script>

<template>
    <div v-if="asyncStatus === 'loading'">Logging in...</div>
    <div v-else-if="error">
        <p>{{ error.message }}</p>
    </div>
</template>

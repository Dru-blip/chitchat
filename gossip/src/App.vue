<script setup lang="ts">
import { useKeys } from "./composables/crypto";
import { PiniaColadaDevtools } from "@pinia/colada-devtools";

const { ready, error } = useKeys();
</script>

<template>
    <main>
        <p v-if="error">{{ error.message }}</p>
        <RouterView v-else-if="ready" />
        <p v-else>Loading..</p>
    </main>
    <PiniaColadaDevtools />
</template>

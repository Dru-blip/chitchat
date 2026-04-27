<script setup lang="ts">
import { ref, type HTMLAttributes } from "vue";

import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { getSerializedIdentityKeys } from "@/lib/local-stores";
import { cn } from "@/lib/utils";
import { useSendMagicLink } from "@/composables/auth";

const props = defineProps<{
    class?: HTMLAttributes["class"];
}>();

const email = ref("");
const { error, success, asyncStatus, sendMagicLink } = useSendMagicLink();

const sendLink = async (e: any) => {
    e.preventDefault();

    const identityKeyPair = await getSerializedIdentityKeys();
    await sendMagicLink({ email: email.value, pubkey: identityKeyPair.pubKey });
};
</script>

<template>
    <div :class="cn('flex flex-col gap-6', props.class)">
        <div v-if="error">
            <p>{{ error.message }}</p>
        </div>
        <div v-if="success">
            <p>Check your email for the login link.</p>
        </div>
        <form v-else @submit="sendLink">
            <FieldGroup>
                <!--<!-- App Name -->
                <!-- <h1 class="text-xl font-bold text-center">Chitchat</h1> -->

                <Field>
                    <FieldLabel for="email"> Email </FieldLabel>
                    <Input
                        id="email"
                        type="email"
                        v-model="email"
                        placeholder="m@example.com"
                        required
                    />
                </Field>

                <Field>
                    <Button type="submit" :disabled="asyncStatus === 'loading'">
                        Login
                    </Button>
                    <span class="text-gray-500">We'll send you login link</span>
                </Field>
            </FieldGroup>
        </form>
    </div>
</template>

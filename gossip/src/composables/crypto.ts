import { initializeKeys } from "@/lib/key-generators";
import { onMounted, ref } from "vue";

export const useKeys = () => {
  const ready = ref(false);
  const error = ref<Error | null>(null);

  onMounted(async () => {
    try {
      await initializeKeys();
    } catch (e) {
      error.value = e as Error;
    } finally {
      ready.value = true;
    }
  });

  return { ready, error };
};

import { defineMutation, useMutation } from "@pinia/colada";
import { ref } from "vue";

interface SendMagicLinkPayload {
  email: string;
  pubkey: string;
}

interface VerifyMagicLinkPayload {
  token: string;
}

interface ErrorResponse {
  code: string;
  message: string;
  details: Record<string, any>;
}

export const useSendMagicLink = defineMutation(() => {
  const success = ref(false);
  const { mutateAsync, ...mutation } = useMutation({
    mutation: async (payload: SendMagicLinkPayload) => {
      const res = await fetch(
        import.meta.env.VITE_API_URL + "auth/send-magic-link",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ ...payload }),
        },
      );

      const data = await res.json();

      if (!res.ok) {
        throw {
          code: data.code,
          message: data.message,
          details: data.details,
        } satisfies ErrorResponse;
      }

      return data;
    },

    onSuccess: () => {
      success.value = true;
    },
  });

  return {
    ...mutation,
    sendMagicLink: (payload: SendMagicLinkPayload) => mutateAsync(payload),
    success,
  };
});

export const useVerifyMagicLink = defineMutation(() => {
  const success = ref(false);
  const { mutateAsync, ...mutation } = useMutation({
    mutation: async (payload: VerifyMagicLinkPayload) => {
      const res = await fetch(
        import.meta.env.VITE_API_URL + "auth/verify-magic-link",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ ...payload }),
        },
      );

      const data = await res.json();

      if (!res.ok) {
        throw {
          code: data.code,
          message: data.message,
          details: data.details,
        } satisfies ErrorResponse;
      }

      return data;
    },

    onSuccess: () => {
      success.value = true;
    },
  });

  return {
    ...mutation,
    verifyMagicLink: (payload: VerifyMagicLinkPayload) => mutateAsync(payload),
    success,
  };
});

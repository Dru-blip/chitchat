"use client";

import { ErrorResponse } from "@/types";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function Page() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const token = searchParams.get("token");

  const verifyLink = async () => {
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + "auth/verify-magic-link",
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token }),
        credentials: "include",
      },
    );

    const data = await res.json();

    if (!res.ok) {
      //TODO: error handling
      throw {
        code: data.code,
        message: data.message,
        details: data.details,
      } satisfies ErrorResponse;
    }
    if (data.onboard) {
      router.push("/onboarding");
    } else {
      router.push("/chat");
    }
  };
  useEffect(() => {
    verifyLink();
  });
  return <div>Verifying link...</div>;
}

"use client";

import { useState, useEffect } from "react";
import { apiFetch } from "@/lib/utils";
import type { User, ErrorResponse } from "@/types";

export function useUser() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiFetch<User, ErrorResponse>("users/me")
      .then(({ data, error }) => {
        if (error) {
          setError(error.message || "Failed to fetch user");
        } else {
          setUser(data);
        }
      })
      .finally(() => setLoading(false));
  }, []);

  return { user, loading, error };
}

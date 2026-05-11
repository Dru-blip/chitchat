"use client";

import { initializeKeys } from "@/lib/key-generator";
import { useEffect } from "react";

export function PreKeyInitializer() {
  useEffect(() => {
    initializeKeys()
      .then(() => {
        console.log("keys initialized and loaded successfully");
      })
      .catch((error) => {
        console.error("failed to initialize keys", error);
      });
  }, []);
  return null;
}

"use client";

import { uploadPrekeysToServer } from "@/lib/prekey-upload";
import { useEffect } from "react";

export const UploadPrekeys = () => {
  useEffect(() => {
    uploadPrekeysToServer({ skipSessionCheck: true }).catch((error) => {
      console.error("failed to upload prekeys", error);
    });
  }, []);

  return null;
};

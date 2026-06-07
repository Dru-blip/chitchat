"use client";

import { createContext, use, useContext } from "react";
import { useUser } from "@/hooks/useUser";
import type { User } from "@/types";

interface UserContextType {
  user: User | null;
  loading: boolean;
  error: string | null;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const { user, loading, error } = useUser();

  return <UserContext value={{ user, loading, error }}>{children}</UserContext>;
}

export function useUserContext() {
  const context = use(UserContext);
  if (context === undefined) {
    throw new Error("useUserContext must be used within a UserProvider");
  }
  return context;
}

"use client";

import { ThemeProvider } from "@/components/theme-provider";
import { UserProvider } from "@/context/user";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <UserProvider>{children}</UserProvider>
    </ThemeProvider>
  );
}

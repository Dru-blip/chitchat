import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export async function apiFetch<T, E = Record<string, unknown>>(
  url: string,
  options?: RequestInit,
): Promise<{ data: T; error: null } | { data: null; error: E }> {
  try {
    const res = await fetch(process.env.NEXT_PUBLIC_API_URL + url, {
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      ...options,
    });

    const data = await res.json();

    if (!res.ok) {
      return { data: null, error: data };
    }

    return { data: data as T, error: null };
  } catch (err) {
    console.error(err);
    return {
      data: null,
      error: { message: "Something went wrong, please try again." } as E,
    };
  }
}

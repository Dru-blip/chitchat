import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { format, isToday, isYesterday, isThisYear } from "date-fns";

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

export function getInitials(name: string | null | undefined): string {
  if (!name) return "?";
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase();
}

export function formatRelativeTime(dateStr: string | undefined): string | null {
  if (!dateStr) return null;
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return null;

  if (isToday(date)) return format(date, "h:mm a");
  if (isYesterday(date)) return "Yesterday";
  if (isThisYear(date)) return format(date, "MMM d");
  return format(date, "MM/dd/yy");
}

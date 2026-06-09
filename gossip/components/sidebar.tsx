"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useUserContext } from "@/context/user";
import { cn, getInitials } from "@/lib/utils";
import {
  BubbleChatIcon,
  Logout03Icon,
  Settings01Icon,
  UserIcon,
} from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { href: "/chats", label: "Chats", icon: BubbleChatIcon, badge: 4 },
];

function UserSectionSkeleton() {
  return (
    <div className="flex items-center gap-3 w-full rounded-lg p-2 animate-pulse">
      <div className="size-8 rounded-full bg-muted" />
      <div className="flex flex-col gap-1.5 flex-1">
        <div className="h-3.5 w-24 bg-muted rounded" />
        <div className="h-3 w-32 bg-muted rounded" />
      </div>
    </div>
  );
}

export const Sidebar = () => {
  const pathname = usePathname();
  const { user, loading, error } = useUserContext();

  const dropdownUserSection = loading ? (
    <UserSectionSkeleton />
  ) : error ? (
    <div className="flex items-center gap-3 w-full rounded-lg p-2">
      <div className="size-8 shrink-0 rounded-full bg-destructive/10 flex items-center justify-center">
        <span className="text-destructive text-xs font-bold">!</span>
      </div>
      <div className="flex flex-col items-start min-w-0 flex-1">
        <span className="text-sm font-medium text-destructive leading-tight truncate w-full">
          Error
        </span>
        <span className="text-xs text-muted-foreground truncate w-full">
          {error}
        </span>
      </div>
    </div>
  ) : (
    <div className="flex items-center gap-3 w-full rounded-lg p-2">
      <Avatar className="size-8 shrink-0 ring-1 ring-border">
        <AvatarImage
          src={user?.image ?? user?.image ?? undefined}
          alt={user?.name ?? "User"}
        />
        <AvatarFallback className="text-[11px] font-semibold bg-muted text-muted-foreground">
          {getInitials(user?.name ?? null)}
        </AvatarFallback>
      </Avatar>
      <div className="flex flex-col items-start min-w-0 flex-1">
        <span className="text-sm font-medium text-sidebar-foreground leading-tight truncate w-full">
          {user?.name ?? "Unknown"}
        </span>
        <span className="text-xs text-muted-foreground truncate w-full">
          {user?.email ?? ""}
        </span>
      </div>
    </div>
  );

  return (
    <>
      <aside className="hidden md:flex flex-col h-screen w-14 border-r border-border bg-sidebar shrink-0">
        <nav className="flex-1 flex flex-col items-center gap-0.5 p-2 pt-3">
          {navItems.map((item) => {
            const isActive =
              pathname === item.href || pathname.startsWith(item.href + "/");
            return (
              <Link
                key={item.href}
                href={item.href}
                title={item.label}
                className={cn(
                  "relative flex items-center justify-center size-10 rounded-lg text-sm font-medium transition-colors duration-150 outline-none",
                  isActive
                    ? "bg-primary text-primary-foreground"
                    : "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
                )}
              >
                <div className="relative">
                  <HugeiconsIcon
                    icon={item.icon}
                    size={20}
                    color="currentColor"
                    strokeWidth={isActive ? 2 : 1.5}
                  />
                  {item.badge > 0 && (
                    <Badge
                      className={cn(
                        "absolute -top-2 -right-2 min-w-4 h-4 px-1 flex items-center justify-center text-[9px] leading-none rounded-full border-0",
                        isActive
                          ? "bg-primary-foreground text-primary"
                          : "bg-primary text-primary-foreground",
                      )}
                    >
                      {item.badge}
                    </Badge>
                  )}
                </div>
              </Link>
            );
          })}
        </nav>

        {/* User */}
        <div className="border-t border-border p-2 shrink-0 flex justify-center">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                className={cn(
                  "flex items-center justify-center size-10 rounded-lg",
                  "hover:bg-sidebar-accent transition-colors duration-150 outline-none",
                )}
              >
                {loading ? (
                  <div className="size-8 rounded-full bg-muted animate-pulse" />
                ) : error ? (
                  <div className="size-8 rounded-full bg-destructive/10 flex items-center justify-center">
                    <span className="text-destructive text-xs font-bold">
                      !
                    </span>
                  </div>
                ) : (
                  <Avatar className="size-8 ring-1 ring-border">
                    <AvatarImage
                      src={user?.image ?? undefined}
                      alt={user?.name ?? "User"}
                    />
                    <AvatarFallback className="text-[11px] font-semibold bg-muted text-muted-foreground">
                      {getInitials(user?.name ?? null)}
                    </AvatarFallback>
                  </Avatar>
                )}
              </button>
            </DropdownMenuTrigger>

            <DropdownMenuContent
              side="top"
              align="start"
              sideOffset={6}
              className="w-52"
            >
              <DropdownMenuLabel className="font-normal py-2">
                <p className="text-[11px] text-muted-foreground mb-0.5">
                  Signed in as
                </p>
                {dropdownUserSection}
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                asChild
                className="gap-2 text-sm cursor-pointer"
              >
                <Link href="/profile">
                  <HugeiconsIcon
                    icon={UserIcon}
                    size={14}
                    color="currentColor"
                    strokeWidth={1.5}
                  />
                  Profile
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem
                asChild
                className="gap-2 text-sm cursor-pointer"
              >
                <Link href="/settings">
                  <HugeiconsIcon
                    icon={Settings01Icon}
                    size={14}
                    color="currentColor"
                    strokeWidth={1.5}
                  />
                  Settings
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="gap-2 text-sm cursor-pointer text-destructive focus:text-destructive">
                <HugeiconsIcon
                  icon={Logout03Icon}
                  size={14}
                  color="currentColor"
                  strokeWidth={1.5}
                />
                Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </aside>

      {/*
          bottom navigation bar
       */}
      <nav className="md:hidden fixed bottom-0 inset-x-0 z-50 bg-background border-t border-border">
        <div className="flex items-center justify-around px-2 h-16">
          {navItems.map((item) => {
            const isActive =
              pathname === item.href || pathname.startsWith(item.href + "/");
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "relative flex flex-col items-center justify-center gap-1 flex-1 py-2 rounded-xl transition-colors outline-none",
                  isActive ? "text-primary" : "text-muted-foreground",
                )}
              >
                {isActive && (
                  <span className="absolute top-0 left-1/2 -translate-x-1/2 w-6 h-0.5 rounded-full bg-primary" />
                )}
                <div className="relative">
                  <HugeiconsIcon
                    icon={item.icon}
                    size={22}
                    color="currentColor"
                    strokeWidth={isActive ? 2 : 1.5}
                  />
                  {item.badge > 0 && (
                    <span
                      className="absolute -top-1.5 -right-1.5 min-w-4 h-4 px-1 flex items-center justify-center rounded-full
bg-primary text-primary-foreground text-[9px] font-bold leading-none"
                    >
                      {item.badge}
                    </span>
                  )}
                </div>
                <span
                  className={cn(
                    "text-[10px] leading-none",
                    isActive ? "font-semibold" : "font-medium",
                  )}
                >
                  {item.label}
                </span>
              </Link>
            );
          })}

          {/* Profile tab */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                className="relative flex flex-col items-center justify-center gap-1 flex-1 py-2 rounded-xl text-muted-foreground
outline-none"
              >
                {loading ? (
                  <div className="size-6 rounded-full bg-muted animate-pulse" />
                ) : error ? (
                  <div className="size-6 rounded-full bg-destructive/10 flex items-center justify-center">
                    <span className="text-destructive text-[9px] font-bold">
                      !
                    </span>
                  </div>
                ) : (
                  <Avatar className="size-6 ring-1 ring-border">
                    <AvatarImage
                      src={user?.image ?? user?.image ?? undefined}
                      alt={user?.name ?? "User"}
                    />
                    <AvatarFallback className="text-[9px] font-semibold">
                      {getInitials(user?.name ?? null)}
                    </AvatarFallback>
                  </Avatar>
                )}
                <span className="text-[10px] font-medium leading-none">
                  Profile
                </span>
              </button>
            </DropdownMenuTrigger>

            <DropdownMenuContent
              side="top"
              align="end"
              sideOffset={8}
              alignOffset={-8}
              className="w-56 mb-1"
            >
              <DropdownMenuLabel className="font-normal py-2">
                {dropdownUserSection}
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                asChild
                className="gap-2 text-sm cursor-pointer"
              >
                <Link href="/profile">
                  <HugeiconsIcon
                    icon={UserIcon}
                    size={14}
                    color="currentColor"
                    strokeWidth={1.5}
                  />
                  Profile
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem
                asChild
                className="gap-2 text-sm cursor-pointer"
              >
                <Link href="/settings">
                  <HugeiconsIcon
                    icon={Settings01Icon}
                    size={14}
                    color="currentColor"
                    strokeWidth={1.5}
                  />
                  Settings
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="gap-2 text-sm cursor-pointer text-destructive focus:text-destructive">
                <HugeiconsIcon
                  icon={Logout03Icon}
                  size={14}
                  color="currentColor"
                  strokeWidth={1.5}
                />
                Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </nav>
    </>
  );
};

"use client";

import { usePathname } from "next/navigation";
import { ConversationsProvider } from "@/context/conversations";
import { ConversationList } from "@/components/conversation-list";
import { NewConversationButton } from "@/components/new-conversation-button";
import { UploadPrekeys } from "@/components/upload-prekeys";
import { cn } from "@/lib/utils";

export default function ChatsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ConversationsProvider>
      <ChatsShell>{children}</ChatsShell>
      <UploadPrekeys />
    </ConversationsProvider>
  );
}

function ChatsShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const conversationId = pathname.startsWith("/chats/")
    ? pathname.split("/chats/")[1]?.split("/")[0] || undefined
    : undefined;

  return (
    <div className="flex h-full">
      {/* Left pane — conversation list */}
      <aside
        className={cn(
          "flex flex-col border-r border-border bg-background",
          "md:w-80 md:shrink-0 md:flex",
          conversationId ? "hidden md:flex" : "flex"
        )}
      >
        <div className="p-2 flex flex-col h-full">
          <div className="flex items-center justify-between p-4">
            <NewConversationButton />
          </div>
          <div className="flex-1 overflow-y-auto">
            <ConversationList activeId={conversationId} />
          </div>
        </div>
      </aside>

      {/* Right pane — chat view or empty state */}
      <main
        className={cn(
          "flex-1 flex flex-col bg-muted/30",
          !conversationId && "hidden md:flex"
        )}
      >
        {children}
      </main>
    </div>
  );
}
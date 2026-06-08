"use client";

import { HugeiconsIcon } from "@hugeicons/react";
import { BubbleChatIcon } from "@hugeicons/core-free-icons";

export function ConversationListEmpty() {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
      <HugeiconsIcon
        icon={BubbleChatIcon}
        size={40}
        className="text-muted-foreground/40 mb-3"
        strokeWidth={1.2}
      />
      <p className="text-sm font-medium text-foreground">
        No conversations yet
      </p>
      <p className="text-xs text-muted-foreground mt-1">
        Start a new conversation to get going
      </p>
    </div>
  );
}

"use client";

import { BubbleChatIcon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";

export function ConversationNotFound() {
  return (
    <div className="flex-1 flex flex-col items-center justify-center gap-3 p-8">
      <div className="size-16 rounded-4xl bg-muted flex items-center justify-center">
        <HugeiconsIcon
          icon={BubbleChatIcon}
          size={28}
          className="text-muted-foreground"
          strokeWidth={1.2}
        />
      </div>
      <div className="text-center">
        <p className="text-sm font-medium text-foreground">
          Conversation not found
        </p>
        <p className="text-xs text-muted-foreground mt-1">
          This conversation may have been deleted
        </p>
      </div>
    </div>
  );
}

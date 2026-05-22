import { HugeiconsIcon } from "@hugeicons/react";
import { BubbleChatIcon } from "@hugeicons/core-free-icons";

export function EmptyChatState() {
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
          Select a conversation
        </p>
        <p className="text-xs text-muted-foreground mt-1">
          Choose from your existing conversations or start a new one
        </p>
      </div>
    </div>
  );
}
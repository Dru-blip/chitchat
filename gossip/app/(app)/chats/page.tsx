import { UploadPrekeys } from "@/components/upload-prekeys";
import { NewConversationButton } from "@/components/new-conversation-button";
import { ConversationsProvider } from "@/context/conversations";

export default function Page() {
  return (
    <ConversationsProvider>
      <div className="flex items-center justify-between p-4">
        <NewConversationButton />
      </div>
      <UploadPrekeys />
    </ConversationsProvider>
  );
}

import { UploadPrekeys } from "@/components/upload-prekeys";
import { NewConversationButton } from "@/components/new-conversation-button";
import { ConversationsProvider } from "@/context/conversations";
import { ConversationList } from "@/components/conversation-list";

export default function Page() {
  return (
    <ConversationsProvider>
      <div className="p-2">
        <div className="flex items-center justify-between p-4">
          <NewConversationButton />
        </div>
        <ConversationList />
      </div>
      <UploadPrekeys />
    </ConversationsProvider>
  );
}

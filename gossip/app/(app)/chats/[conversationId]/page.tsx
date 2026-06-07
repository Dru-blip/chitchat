import { ChatView } from "@/components/chat/chat-view";
import { LoadActiveConversation } from "@/components/chat/load-active-conversation";

export default async function ConversationPage() {
  return (
    <>
      <LoadActiveConversation />
      <ChatView />
    </>
  );
}

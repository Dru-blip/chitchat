import { Sidebar } from "@/components/sidebar";
import { WebsocketProvider } from "@/context/websockets";
import { ActiveConversationStoreProvider } from "@/stores/providers/active-conversation";
import { ConversationStoreProvider } from "@/stores/providers/conversation";
import { MessageStoreProvider } from "@/stores/providers/messages";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <ConversationStoreProvider>
      <MessageStoreProvider>
        <ActiveConversationStoreProvider>
          <WebsocketProvider>
            <section className="flex h-dvh">
              <Sidebar />
              <div className="flex-1 overflow-hidden">{children}</div>
            </section>
          </WebsocketProvider>
        </ActiveConversationStoreProvider>
      </MessageStoreProvider>
    </ConversationStoreProvider>
  );
}

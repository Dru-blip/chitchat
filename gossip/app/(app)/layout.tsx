import { Sidebar } from "@/components/sidebar";
import { WebsocketProvider } from "@/context/websockets";
import { ConversationStoreProvider } from "@/stores/providers/conversation";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <ConversationStoreProvider>
      <WebsocketProvider>
        <section className="flex h-dvh">
          <Sidebar />
          <div className="flex-1 overflow-hidden">{children}</div>
        </section>
      </WebsocketProvider>
    </ConversationStoreProvider>
  );
}

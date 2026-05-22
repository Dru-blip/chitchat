import { Sidebar } from "@/components/sidebar";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <section className="flex h-dvh">
      <Sidebar />
      <div className="flex-1 overflow-hidden">{children}</div>
    </section>
  );
}

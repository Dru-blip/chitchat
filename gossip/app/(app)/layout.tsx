import { Sidebar } from "@/components/sidebar";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <section className="flex">
      <Sidebar />
      <div>{children}</div>
    </section>
  );
}

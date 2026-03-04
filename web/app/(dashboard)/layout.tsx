import Footer from "@/components/Footer";
import Sidebar from "@/components/Sidebar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 md:flex">
      <Sidebar />
      <div className="flex min-h-screen flex-1 flex-col">
        <main className="flex-1 px-4 py-6 sm:px-6 lg:px-10 lg:py-8">{children}</main>
        <div className="px-4 pb-8 sm:px-6 lg:px-10">
          <Footer />
        </div>
      </div>
    </div>
  );
}

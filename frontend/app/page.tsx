import { ActiveJobBanner } from "@/components/dashboard/active-job-banner";
import { DashboardHero } from "@/components/dashboard/dashboard-hero";
import { HowItWorks } from "@/components/dashboard/how-it-works";
import { ToolCatalog } from "@/components/dashboard/tool-catalog";

export default function HomePage() {
  return (
    <main id="main-content" className="mx-auto w-full max-w-6xl px-4 pb-12 sm:px-6 lg:px-8">
      <ActiveJobBanner />
      <DashboardHero />
      <ToolCatalog />
      <HowItWorks />
    </main>
  );
}

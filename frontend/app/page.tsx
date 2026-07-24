import { ActiveJobBanner } from "@/components/dashboard/active-job-banner";
import { DashboardHero } from "@/components/dashboard/dashboard-hero";
import { HowItWorks } from "@/components/dashboard/how-it-works";
import { ToolCatalog } from "@/components/dashboard/tool-catalog";

export default function HomePage() {
  return (
    <main id="main-content" className="mx-auto w-full max-w-[92rem] px-4 pb-16 sm:px-8 lg:px-10 lg:pb-20">
      <ActiveJobBanner />
      <DashboardHero />
      <ToolCatalog />
      <HowItWorks />
    </main>
  );
}

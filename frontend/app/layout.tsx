import type { Metadata } from "next";
import "@fontsource-variable/archivo";
import "@fontsource/ibm-plex-mono/400.css";
import "@fontsource/ibm-plex-mono/500.css";
import "@fontsource/ibm-plex-mono/600.css";

import { MobileHeader } from "@/components/layout/mobile-header";
import { Sidebar } from "@/components/layout/sidebar";
import { QueryProvider } from "@/components/providers/query-provider";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { WorkspaceProvider } from "@/components/providers/workspace-provider";

import "./globals.css";

export const metadata: Metadata = {
  title: {
    default: "Recode",
    template: "%s · Recode",
  },
  description:
    "Account-free tools for converting and processing images and videos.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <ThemeProvider>
          <QueryProvider>
            <WorkspaceProvider>
              <a
                href="#main-content"
                className="fixed top-3 left-3 z-[100] -translate-y-24 border-2 border-ink bg-accent px-4 py-3 text-sm font-bold text-[#171a17] focus:translate-y-0"
              >
                Skip to content
              </a>
              <Sidebar />
              <div className="min-h-svh lg:pl-[19rem]">
                <MobileHeader />
                {children}
              </div>
            </WorkspaceProvider>
          </QueryProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}

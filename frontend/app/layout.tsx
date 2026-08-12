import type { Metadata, Viewport } from "next";
import "@fontsource-variable/manrope";
import "@fontsource/ibm-plex-mono/400.css";
import "@fontsource/ibm-plex-mono/500.css";
import "@fontsource/ibm-plex-mono/600.css";

import { MobileHeader } from "@/components/layout/mobile-header";
import { GoogleAnalytics } from "@next/third-parties/google";
import { Sidebar } from "@/components/layout/sidebar";
import { QueryProvider } from "@/components/providers/query-provider";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { WorkspaceProvider } from "@/components/providers/workspace-provider";
import { getSiteUrl, siteConfig } from "@/config/site";

import "./globals.css";

export const metadata: Metadata = {
  metadataBase: getSiteUrl(),
  title: {
    default: siteConfig.title,
    template: "%s · Recode",
  },
  description: siteConfig.description,
  applicationName: siteConfig.name,
  keywords: [...siteConfig.keywords],
  authors: [{ name: siteConfig.name }],
  creator: siteConfig.name,
  publisher: siteConfig.name,
  category: "technology",
  alternates: {
    canonical: "/",
  },
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "/",
    siteName: siteConfig.name,
    title: siteConfig.title,
    description: siteConfig.description,
    images: [
      {
        url: "/opengraph-image",
        width: 1200,
        height: 630,
        alt: "Recode — private, account-free image tools",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: siteConfig.title,
    description: siteConfig.description,
    images: ["/opengraph-image"],
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-image-preview": "large",
      "max-snippet": -1,
      "max-video-preview": 0,
    },
  },
  icons: {
    icon: "/icon.svg",
  },
  manifest: "/manifest.webmanifest",
  appleWebApp: {
    capable: true,
    title: siteConfig.name,
    statusBarStyle: "default",
  },
  formatDetection: {
    telephone: false,
    email: false,
    address: false,
  },
};

export const viewport: Viewport = {
  colorScheme: "light dark",
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#f7f6f3" },
    { media: "(prefers-color-scheme: dark)", color: "#181a19" },
  ],
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
                className="fixed top-3 left-3 z-[100] -translate-y-24 rounded-md bg-ink px-4 py-2.5 text-sm font-semibold text-canvas shadow-soft focus:translate-y-0"
              >
                Skip to content
              </a>
              <Sidebar />
              <div className="min-h-svh lg:pl-64">
                <MobileHeader />
                {children}

                <GoogleAnalytics gaId="G-M7KBWRJ0C3" />
              </div>
            </WorkspaceProvider>
          </QueryProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}

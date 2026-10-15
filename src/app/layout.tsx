import type { Metadata } from "next";
import { Archivo, Instrument_Sans } from "next/font/google";
import "./globals.css";
import "./admin.css";
import { SiteHeader } from "@/components/templates/gas-grill/site-header";
import { SiteFooter } from "@/components/templates/gas-grill/site-footer";
import { Reveal } from "@/components/reveal";
import { siteUrl } from "@/lib/seo/site-url";

const archivo = Archivo({
  subsets: ["latin"],
  weight: ["600", "700"],
  variable: "--font-archivo",
  display: "swap",
});

const instrumentSans = Instrument_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  variable: "--font-instrument-sans",
  display: "swap",
});

export const metadata: Metadata = {
  ...(siteUrl ? { metadataBase: siteUrl } : {}),
  title: { default: "Outdoor cooking, considered | EMBER / FIELD", template: "%s | EMBER / FIELD" },
  description: "Purpose-built gas grills and outdoor cooking equipment for the way you gather.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    // `suppressHydrationWarning` covers the `js` class that Reveal adds to <html>
    // before first paint. Server HTML has no way to know scripting is available.
    // `data-scroll-behavior="smooth"` tells Next.js the smooth scrolling declared in
    // globals.css is intentional, so it keeps suppressing smooth scroll on route changes.
    <html lang="en" className={`${archivo.variable} ${instrumentSans.variable}`} data-scroll-behavior="smooth" suppressHydrationWarning>
      <body>
        <SiteHeader />
        <main>{children}</main>
        <SiteFooter />
        <Reveal />
      </body>
    </html>
  );
}

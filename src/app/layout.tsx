import type { Metadata } from "next";
import "./globals.css";
import { SiteHeader } from "@/components/templates/gas-grill/site-header";
import { SiteFooter } from "@/components/templates/gas-grill/site-footer";
import { siteUrl } from "@/lib/seo/site-url";

export const metadata: Metadata = {
  ...(siteUrl ? { metadataBase: siteUrl } : {}),
  title: { default: "Outdoor cooking, considered | EMBER / FIELD", template: "%s | EMBER / FIELD" },
  description: "Purpose-built gas grills and outdoor cooking equipment for the way you gather.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <div className="announcement">Complimentary delivery on every grill · Built for the long weekend</div>
        <SiteHeader />
        <main>{children}</main>
        <SiteFooter />
      </body>
    </html>
  );
}

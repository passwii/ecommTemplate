import type { Metadata } from "next";
import "./globals.css";
import { SiteHeader } from "@/components/templates/gas-grill/site-header";
import { SiteFooter } from "@/components/templates/gas-grill/site-footer";

export const metadata: Metadata = {
  title: { default: "EMBER / FIELD | Outdoor cooking, considered", template: "%s | EMBER / FIELD" },
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

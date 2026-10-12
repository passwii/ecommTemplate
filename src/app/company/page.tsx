import type { Metadata } from "next";
import Link from "next/link";
import { CompanyStory } from "@/components/templates/gas-grill/company-story";
import { CatalogUnavailable } from "@/components/catalog-unavailable";
import { fetchCompanyContent } from "@/lib/catalog/api";

export const revalidate = 60;

export const metadata: Metadata = {
  title: "Company",
  description: "Company, manufacturing, quality, credentials, customer examples and frequently asked questions.",
};

export default async function CompanyPage() {
  const content = await fetchCompanyContent();

  if (!content.ok) {
    return (
      <div className="company-page page-wrap">
        <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>Company</span></nav>
        <CatalogUnavailable heading="Company information is temporarily unavailable." />
      </div>
    );
  }

  return <CompanyStory content={content.data} />;
}

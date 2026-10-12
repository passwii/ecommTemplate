import type { Metadata } from "next";
import { CompanyStory } from "@/components/templates/gas-grill/company-story";

export const metadata: Metadata = {
  title: "Company",
  description: "Company, manufacturing, quality, credentials, customer examples and frequently asked questions.",
};

export default function CompanyPage() {
  return <CompanyStory />;
}

import type { Metadata } from "next";
import "../../admin.css";
import { CompanyAdmin } from "@/components/company-admin";

export const metadata: Metadata = { title: "Company content | Admin" };

export default function AdminCompanyPage() {
  return <CompanyAdmin />;
}

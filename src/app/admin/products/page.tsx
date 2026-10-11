import type { Metadata } from "next";
import "../../admin.css";
import { ProductAdmin } from "@/components/product-admin";

export const metadata: Metadata = {
  title: "Products | Admin",
  robots: { index: false, follow: false },
};

export default function AdminProductsPage() {
  return <ProductAdmin />;
}

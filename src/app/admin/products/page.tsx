import type { Metadata } from "next";
import "../../admin.css";
import { ProductAdmin } from "@/components/product-admin";

export const metadata: Metadata = { title: "Products | Admin" };

export default function AdminProductsPage() {
  return <ProductAdmin />;
}

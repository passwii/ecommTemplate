import type { Metadata } from "next";
import Link from "next/link";
import { InquiryForm } from "@/components/templates/gas-grill/inquiry-form";
import { fetchProduct } from "@/lib/catalog/api";

export const metadata: Metadata = {
  title: "Request a quote",
  description: "Send your product, quantity and customization requirements and our team will reply with a quote.",
};

type Props = { searchParams: Promise<{ product?: string }> };

export default async function InquiriesPage({ searchParams }: Props) {
  const { product: productSlug } = await searchParams;

  const productResult = productSlug ? await fetchProduct(productSlug) : undefined;
  const product = productResult?.ok ? productResult.data : undefined;

  return (
    <div className="inquiry-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>Request a quote</span></nav>
      <header className="inquiry-page-heading">
        <p className="eyebrow">REQUEST A QUOTE</p>
        <h1>Tell us what<br />you need.</h1>
        <p>Share the product, quantity and any customization requirements. Our team replies with pricing, lead time and shipping options.</p>
      </header>

      <div className="inquiry-layout">
        <InquiryForm
          product={product ? { id: product.id, name: product.name } : undefined}
          sourcePath={product ? `/products/${product.slug}` : "/inquiries"}
        />
        <aside className="inquiry-aside">
          <h2>What happens next</h2>
          <ol>
            <li>We review your requirements and product selection.</li>
            <li>We confirm specifications, quantity and customization options.</li>
            <li>We send a quotation with lead time and shipping details.</li>
          </ol>
          <p className="inquiry-aside-note">Sample storefront. Confirm specifications and trade terms with the merchant. No checkout or payment is connected.</p>
        </aside>
      </div>
    </div>
  );
}

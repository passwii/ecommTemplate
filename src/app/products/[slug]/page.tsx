import type { Metadata } from "next";
import { notFound } from "next/navigation";
import Link from "next/link";
import { StorefrontProductDetail } from "@/components/storefront-product-detail";
import { CatalogUnavailable } from "@/components/catalog-unavailable";
import { fetchProduct } from "@/lib/catalog/api";
import { createPageMetadata } from "@/lib/seo/metadata";

type Props = { params: Promise<{ slug: string }> };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const product = await fetchProduct(slug);
  if (!product.ok) return {};

  const primaryImage = [...product.data.images].sort((a, b) => a.sortOrder - b.sortOrder)[0];
  return createPageMetadata({
    title: product.data.name,
    description: product.data.shortDescription,
    pathname: `/products/${slug}`,
    image: primaryImage ?? { url: "/images/grill-hero.jpg", alt: product.data.name },
  });
}

export default async function ProductPage({ params }: Props) {
  const { slug } = await params;
  const product = await fetchProduct(slug);

  if (!product.ok && product.reason === "not_found") notFound();

  if (!product.ok) {
    return (
      <div className="product-page page-wrap">
        <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><span>Product</span></div>
        <CatalogUnavailable heading="This product is temporarily unavailable." />
      </div>
    );
  }

  return <StorefrontProductDetail product={product.data} />;
}

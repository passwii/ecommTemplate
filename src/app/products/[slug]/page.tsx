import type { Metadata } from "next";
import { getPublicStorefrontProducts } from "@/lib/catalog/seo";
import { createPageMetadata } from "@/lib/seo/metadata";
import { notFound } from "next/navigation";
import { StorefrontProductDetail } from "@/components/storefront-product-detail";

type Props = { params: Promise<{ slug: string }> };

export const dynamicParams = false;

export function generateStaticParams() { return getPublicStorefrontProducts().map(({ slug }) => ({ slug })); }
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const product = getPublicStorefrontProducts().find((item) => item.slug === slug);
  if (!product) return {};

  const primaryImage = [...product.images].sort((a, b) => a.sortOrder - b.sortOrder)[0] ?? { url: product.image, alt: product.imageAlt };
  return createPageMetadata({
    title: product.name,
    description: product.shortDescription,
    pathname: `/products/${slug}`,
    image: primaryImage,
  });
}

export default async function ProductPage({ params }: Props) {
  const { slug } = await params;
  if (!getPublicStorefrontProducts().some((product) => product.slug === slug)) notFound();
  return <StorefrontProductDetail slug={slug} />;
}

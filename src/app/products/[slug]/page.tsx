import type { Metadata } from "next";
import { products } from "@/lib/catalog/mock-data";
import { StorefrontProductDetail } from "@/components/storefront-product-detail";

type Props = { params: Promise<{ slug: string }> };
export function generateStaticParams() { return products.map(({ slug }) => ({ slug })); }
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const product = products.find((item) => item.slug === slug && item.status === "published");
  return product ? { title: product.name, description: product.shortDescription } : {};
}

export default async function ProductPage({ params }: Props) {
  const { slug } = await params;
  return <StorefrontProductDetail slug={slug} />;
}

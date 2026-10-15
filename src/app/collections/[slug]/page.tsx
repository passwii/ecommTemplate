import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { getPublishedCollections, getCollection, getPublishedProducts } from "@/lib/catalog/mock-data";
import { createPageMetadata } from "@/lib/seo/metadata";
import { StorefrontCollectionCatalog } from "@/components/storefront-product-grid";

type Props = { params: Promise<{ slug: string }> };

export const dynamicParams = false;

export function generateStaticParams() { return getPublishedCollections().map(({ slug }) => ({ slug })); }

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const collection = getCollection(slug);
  if (!collection || !getPublishedCollections().some(({ slug: publishedSlug }) => publishedSlug === slug)) return {};

  const representativeProduct = getPublishedProducts().find((product) => product.collectionSlug === slug);
  return createPageMetadata({
    title: collection.name,
    description: collection.description,
    pathname: `/collections/${slug}`,
    image: {
      url: representativeProduct?.image ?? "/images/grill-hero.jpg",
      alt: representativeProduct?.imageAlt ?? "Outdoor patio with lounge seating beneath string lights",
    },
  });
}

export default async function CollectionPage({ params }: Props) {
  const { slug } = await params;
  const collection = getCollection(slug);
  if (!collection || !getPublishedCollections().some(({ slug: publishedSlug }) => publishedSlug === slug)) notFound();
  return (
    <div className="catalog-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>{collection.name}</span></nav>
      <header className="catalog-heading" data-reveal>
        <h1>{collection.name} for<br /><em>open-air living.</em></h1>
        <p>{collection.description}</p>
      </header>
      <div data-reveal style={{ "--reveal-delay": "80ms" } as React.CSSProperties}>
        <StorefrontCollectionCatalog slug={slug} />
      </div>
    </div>
  );
}

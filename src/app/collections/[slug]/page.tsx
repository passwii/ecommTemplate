import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { StorefrontCollectionCatalog } from "@/components/storefront-product-grid";
import { CatalogUnavailable } from "@/components/catalog-unavailable";
import { fetchCollection, fetchProducts } from "@/lib/catalog/api";
import { createPageMetadata } from "@/lib/seo/metadata";

type Props = { params: Promise<{ slug: string }> };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const collection = await fetchCollection(slug);
  if (!collection.ok) return {};

  const products = await fetchProducts({ collection: slug, limit: 100 });
  const representativeProduct = products.ok ? products.data[0] : undefined;

  return createPageMetadata({
    title: collection.data.name,
    description: collection.data.description,
    pathname: `/collections/${slug}`,
    image: {
      url: representativeProduct?.image || "/images/grill-hero.jpg",
      alt: representativeProduct?.imageAlt || "Outdoor patio with lounge seating beneath string lights",
    },
  });
}

export default async function CollectionPage({ params }: Props) {
  const { slug } = await params;
  const [collection, products] = await Promise.all([
    fetchCollection(slug),
    fetchProducts({ collection: slug, limit: 100 }),
  ]);

  if (!collection.ok && collection.reason === "not_found") notFound();

  if (!collection.ok) {
    return (
      <div className="catalog-page page-wrap">
        <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><span>Collection</span></div>
        <CatalogUnavailable heading="This collection is temporarily unavailable." />
      </div>
    );
  }

  return (
    <div className="catalog-page page-wrap">
      <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><span>{collection.data.name}</span></div>
      <header className="catalog-heading"><p className="eyebrow">The Ember / Field collection</p><h1>{collection.data.name} for<br /><em>open-air living.</em></h1><p>{collection.data.description}</p></header>
      {products.ok
        ? <StorefrontCollectionCatalog products={products.data} />
        : <CatalogUnavailable heading="Products in this collection are temporarily unavailable." />}
    </div>
  );
}

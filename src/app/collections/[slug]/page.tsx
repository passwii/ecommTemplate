import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { collections, getCollection } from "@/lib/catalog/mock-data";
import { StorefrontCollectionCatalog } from "@/components/storefront-product-grid";
import { CONTACT_PATH } from "@/lib/contact";

type Props = { params: Promise<{ slug: string }> };

export function generateStaticParams() { return collections.map(({ slug }) => ({ slug })); }

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const collection = getCollection(slug);
  return collection ? { title: collection.name } : {};
}

export default async function CollectionPage({ params }: Props) {
  const { slug } = await params;
  const collection = getCollection(slug);
  if (!collection) notFound();
  return (
    <div className="catalog-page page-wrap">
      <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><span>{collection.name}</span></div>
      <header className="catalog-heading"><p className="eyebrow">The Ember / Field collection</p><h1>{collection.name} for<br /><em>open-air living.</em></h1><p>{collection.description}</p><Link className="button button-dark catalog-contact" href={CONTACT_PATH}>Contact <span aria-hidden="true">↗</span></Link></header>
      <StorefrontCollectionCatalog slug={slug} />
    </div>
  );
}

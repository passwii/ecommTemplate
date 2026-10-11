import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { collections, products, getCollection } from "@/lib/catalog/mock-data";
import { ProductCard } from "@/components/templates/gas-grill/product-card";

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
  const items = products.filter((product) => product.collectionSlug === slug);

  return (
    <div className="catalog-page page-wrap">
      <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><span>{collection.name}</span></div>
      <header className="catalog-heading"><p className="eyebrow">The Ember / Field collection</p><h1>{collection.name} for<br /><em>open-air living.</em></h1><p>{collection.description}</p></header>
      <div className="catalog-toolbar"><span>{items.length} considered designs</span><span>Natural gas and propane</span></div>
      {items.length ? <div className="product-grid">{items.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div> : <div className="empty-state"><h2>More designs are on their way.</h2><Link className="text-link" href="/">Explore all grills</Link></div>}
    </div>
  );
}

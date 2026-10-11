"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { products, type Product } from "@/lib/catalog/mock-data";
import { getStoredProducts } from "@/lib/catalog/product-admin";
import { ProductCard } from "@/components/templates/gas-grill/product-card";

export function StorefrontProductGrid({ collectionSlug, featured = false, limit }: { collectionSlug?: string; featured?: boolean; limit?: number }) {
  const [catalog, setCatalog] = useState(products);

  useEffect(() => {
    setCatalog(getStoredProducts());
  }, []);

  let visible: Product[] = catalog
    .filter((product) => product.status === "published")
    .filter((product) => !collectionSlug || product.collectionSlug === collectionSlug)
    .sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name));
  if (featured) visible = visible.filter((product) => product.featured);
  if (limit) visible = visible.slice(0, limit);

  return <div className="product-grid">{visible.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div>;
}

export function StorefrontCollectionCatalog({ slug }: { slug: string }) {
  const [catalog, setCatalog] = useState(products);

  useEffect(() => {
    setCatalog(getStoredProducts());
  }, []);

  const items = catalog
    .filter((product) => product.status === "published" && product.collectionSlug === slug)
    .sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name));

  return (
    <>
      <div className="catalog-toolbar"><span>{items.length} considered designs</span><span>Natural gas and propane</span></div>
      {items.length ? <div className="product-grid">{items.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div> : <div className="empty-state"><h2>No published designs in this collection.</h2><Link className="text-link" href="/">Explore all grills</Link></div>}
    </>
  );
}

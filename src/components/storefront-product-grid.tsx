import Link from "next/link";
import type { Product } from "@/lib/catalog/types";
import { ProductCard } from "@/components/templates/gas-grill/product-card";

export function StorefrontProductGrid({ products }: { products: Product[] }) {
  return <div className="product-grid">{products.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div>;
}

export function StorefrontCollectionCatalog({ products }: { products: Product[] }) {
  return (
    <>
      <div className="catalog-toolbar"><span>{products.length} considered designs</span><span>Natural gas and propane</span></div>
      {products.length ? (
        <div className="product-grid">{products.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div>
      ) : (
        <div className="empty-state"><h2>No published designs in this collection.</h2><Link className="text-link" href="/">Explore all grills</Link></div>
      )}
    </>
  );
}

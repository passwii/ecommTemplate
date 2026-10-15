import Link from "next/link";
import type { Product } from "@/lib/catalog/types";
import { ProductCard } from "@/components/templates/gas-grill/product-card";

/**
 * The home page grid sits under the hero, so its cards never opt into image
 * priority. The hero image is that route's largest contentful paint.
 */
export function StorefrontProductGrid({ products, featured = true }: { products: Product[]; featured?: boolean }) {
  const mode = featured ? "product-grid is-featured" : "product-grid";
  return <div className={mode}>{products.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div>;
}

export function StorefrontCollectionCatalog({ products }: { products: Product[] }) {
  return (
    <>
      <div className="catalog-toolbar">
        <span>{products.length} {products.length === 1 ? "design" : "designs"}</span>
        <span>Natural gas and propane</span>
      </div>
      {products.length ? (
        <div className="product-grid">{products.map((product, index) => <ProductCard key={product.id} product={product} index={index} priority={index === 0} />)}</div>
      ) : (
        <div className="empty-state">
          <h2>No published designs in this collection.</h2>
          <p>This collection is empty right now. The full range is still available to browse.</p>
          <Link className="text-link" href="/#collection">Explore grills</Link>
        </div>
      )}
    </>
  );
}

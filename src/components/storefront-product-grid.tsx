import Link from "next/link";
import { getIndexableProducts } from "@/lib/catalog/seo";
import { ProductCard } from "@/components/templates/gas-grill/product-card";

export function StorefrontProductGrid({ collectionSlug, featured = false, limit }: { collectionSlug?: string; featured?: boolean; limit?: number }) {
  let visible = getIndexableProducts()
    .filter((product) => !collectionSlug || product.collectionSlug === collectionSlug)
    .sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name));
  if (featured) visible = visible.filter((product) => product.featured);
  if (limit) visible = visible.slice(0, limit);

  const mode = featured ? " product-grid is-featured" : " product-grid";
  return <div className={mode.trim()}>{visible.map((product, index) => <ProductCard key={product.id} product={product} index={index} />)}</div>;
}

export function StorefrontCollectionCatalog({ slug }: { slug: string }) {
  const items = getIndexableProducts()
    .filter((product) => product.collectionSlug === slug)
    .sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name));

  return (
    <>
      <div className="catalog-toolbar">
        <span>{items.length} {items.length === 1 ? "design" : "designs"}</span>
        <span>Natural gas and propane</span>
      </div>
      {items.length
        ? <div className="product-grid">{items.map((product, index) => <ProductCard key={product.id} product={product} index={index} priority={index === 0} />)}</div>
        : <div className="empty-state">
            <h2>No published designs in this collection.</h2>
            <p>This collection is empty right now. The full range is still available to browse.</p>
            <Link className="text-link" href="/#collection">Explore grills</Link>
          </div>}
    </>
  );
}

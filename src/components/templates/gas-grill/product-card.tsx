import Image from "next/image";
import Link from "next/link";
import type { Product } from "@/lib/catalog/types";
import { formatPrice } from "@/lib/catalog/mock-data";

/**
 * The comparison facts a consumer scans for. Rendered as chips under the
 * summary so values line up column-to-column across a grid row.
 */
function comparisonFacts(product: Product) {
  return [
    product.burners ? `${product.burners} burners` : null,
    product.cookingArea ?? null,
    product.fuel ?? null,
  ].filter((fact): fact is string => Boolean(fact));
}

/**
 * `priority` is set by the collection grid for its first card only, which is the
 * largest contentful paint on that route. Grids that sit under a hero leave it off.
 * A product with no published image renders the placeholder instead of a broken
 * frame, and a product with no published price renders no price line at all.
 */
export function ProductCard({ product, index = 0, priority = false }: { product: Product; index?: number; priority?: boolean }) {
  const facts = comparisonFacts(product);

  return (
    <article className={`product-card product-card-${index + 1}`}>
      <Link className="product-image" href={`/products/${product.slug}`} aria-label={`View ${product.name}`}>
        {product.image ? (
          <Image src={product.image} alt={product.imageAlt} fill priority={priority} unoptimized={!product.image.startsWith("/")} sizes="(max-width: 767px) 90vw, (max-width: 1023px) 45vw, 30vw" />
        ) : (
          <span className="product-image-placeholder" aria-hidden="true" />
        )}
        {product.badge && <span className="product-badge">{product.badge}</span>}
        <span className="image-arrow" aria-hidden="true">↗</span>
      </Link>
      <div className="product-info">
        <div>
          <p className="product-collection">{product.collection}</p>
          <h3><Link href={`/products/${product.slug}`}>{product.name}</Link></h3>
          <p className="product-summary">{product.summary}</p>
        </div>
        {product.price > 0 && <p className="product-price">{formatPrice(product.price, product.currency)}</p>}
      </div>
      {facts.length > 0 && <ul className="product-specs">{facts.map((fact) => <li key={fact}>{fact}</li>)}</ul>}
    </article>
  );
}

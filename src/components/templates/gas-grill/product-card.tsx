import Image from "next/image";
import Link from "next/link";
import type { Product } from "@/lib/catalog/mock-data";
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

export function ProductCard({ product, index = 0 }: { product: Product; index?: number }) {
  const facts = comparisonFacts(product);

  return (
    <article className={`product-card product-card-${index + 1}`}>
      <Link className="product-image" href={`/products/${product.slug}`} aria-label={`View ${product.name}`}>
        <Image src={product.image} alt={product.imageAlt} fill unoptimized={!product.image.startsWith("/")} sizes="(max-width: 767px) 90vw, (max-width: 1023px) 45vw, 30vw" />
        {product.badge && <span className="product-badge">{product.badge}</span>}
        <span className="image-arrow" aria-hidden="true">↗</span>
      </Link>
      <div className="product-info">
        <div>
          <p className="product-collection">{product.collection}</p>
          <h3><Link href={`/products/${product.slug}`}>{product.name}</Link></h3>
          <p className="product-summary">{product.summary}</p>
        </div>
        <p className="product-price">{formatPrice(product.price, product.currency)}</p>
      </div>
      {facts.length > 0 && <ul className="product-specs">{facts.map((fact) => <li key={fact}>{fact}</li>)}</ul>}
    </article>
  );
}

import Image from "next/image";
import Link from "next/link";
import type { Product } from "@/lib/catalog/mock-data";
import { formatPrice } from "@/lib/catalog/mock-data";

export function ProductCard({ product, index = 0 }: { product: Product; index?: number }) {
  return (
    <article className={`product-card product-card-${index + 1}`}>
      <Link className="product-image" href={`/products/${product.slug}`} aria-label={`View ${product.name}`}>
        <Image src={product.image} alt={product.imageAlt} fill unoptimized sizes="(max-width: 700px) 90vw, (max-width: 1100px) 45vw, 30vw" />
        {product.badge && <span className="product-badge">{product.badge}</span>}
        <span className="image-arrow" aria-hidden="true">↗</span>
      </Link>
      <div className="product-info">
        <div><p className="product-collection">{product.collection}</p><h3><Link href={`/products/${product.slug}`}>{product.name}</Link></h3><p className="product-summary">{product.summary}</p></div>
        <p className="product-price">{formatPrice(product.price)}</p>
      </div>
    </article>
  );
}

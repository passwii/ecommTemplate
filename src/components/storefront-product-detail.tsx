"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { formatPrice, products, type Product } from "@/lib/catalog/mock-data";
import { getStoredProducts } from "@/lib/catalog/product-admin";

export function StorefrontProductDetail({ slug }: { slug: string }) {
  const [product, setProduct] = useState<Product | undefined>(() => products.find((item) => item.slug === slug && item.status === "published"));
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    setProduct(getStoredProducts().find((item) => item.slug === slug && item.status === "published"));
    setLoaded(true);
  }, [slug]);

  if (!loaded && !product) {
    return <div className="product-page page-wrap" aria-busy="true"><p className="admin-loading">Loading product…</p></div>;
  }

  if (!product) {
    return <div className="product-page page-wrap"><div className="empty-state"><h1>This product is not available.</h1><Link className="text-link" href="/">Explore all grills <span aria-hidden="true">↗</span></Link></div></div>;
  }

  const primaryImage = [...product.images].sort((a, b) => a.sortOrder - b.sortOrder)[0] ?? { url: product.image, alt: product.imageAlt, sortOrder: 0 };

  return (
    <div className="product-page page-wrap">
      <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><Link href={`/collections/${product.collectionSlug}`}>{product.collection}</Link><span>/</span><span>{product.name}</span></div>
      <div className="product-detail">
        <div className="product-detail-image">
          <Image src={primaryImage.url} alt={primaryImage.alt} fill priority unoptimized sizes="(max-width: 760px) 100vw, 58vw" />
          {product.badge && <span className="product-badge">{product.badge}</span>}
        </div>
        <div className="product-detail-copy">
          <p className="eyebrow">{product.collection}</p>
          <h1>{product.name}</h1>
          <p className="detail-price">{formatPrice(product.price, product.currency)}</p>
          {product.compareAtPrice && <p className="detail-compare-price">Was {formatPrice(product.compareAtPrice, product.currency)}</p>}
          <p className="detail-summary-lead">{product.shortDescription}</p>
          <p className="detail-summary">{product.description}</p>
          {product.specifications.length > 0 && <section className="detail-specs-section" aria-labelledby="detail-specs-heading">
            <h2 id="detail-specs-heading">At a glance</h2>
            <dl className="detail-specs">{product.specifications.map((specification) => <div key={specification.name}><dt>{specification.name}</dt><dd>{specification.value}</dd></div>)}</dl>
          </section>}
          <a className="button button-dark detail-cta" href={`mailto:hello@emberandfield.example?subject=${encodeURIComponent(`Enquiry about ${product.name}`)}`}>Ask about this grill <span aria-hidden="true">↗</span></a>
          <p className="detail-note">Sample product information. No checkout or payment is connected.</p>
        </div>
      </div>
    </div>
  );
}

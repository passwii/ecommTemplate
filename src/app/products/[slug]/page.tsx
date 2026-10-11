import type { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { products, getProduct, formatPrice } from "@/lib/catalog/mock-data";

type Props = { params: Promise<{ slug: string }> };
export function generateStaticParams() { return products.map(({ slug }) => ({ slug })); }
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const product = getProduct(slug);
  return product ? { title: product.name, description: product.summary } : {};
}

export default async function ProductPage({ params }: Props) {
  const { slug } = await params;
  const product = getProduct(slug);
  if (!product) notFound();
  return (
    <div className="product-page page-wrap">
      <div className="breadcrumbs"><Link href="/">Home</Link><span>/</span><Link href={`/collections/${product.collectionSlug}`}>{product.collection}</Link><span>/</span><span>{product.name}</span></div>
      <div className="product-detail">
        <div className="product-detail-image">
          <Image src={product.image} alt={product.imageAlt} fill priority sizes="(max-width: 760px) 100vw, 58vw" />
          {product.badge && <span className="product-badge">{product.badge}</span>}
        </div>
        <div className="product-detail-copy">
          <p className="eyebrow">{product.collection}</p>
          <h1>{product.name}</h1>
          <p className="detail-price">{formatPrice(product.price)}</p>
          <p className="detail-summary-lead">{product.summary}</p>
          <p className="detail-summary">{product.description}</p>
          <section className="detail-specs-section" aria-labelledby="detail-specs-heading">
            <h2 id="detail-specs-heading">At a glance</h2>
            <dl className="detail-specs">
              <div><dt>Burners</dt><dd>{product.burners}</dd></div>
              <div><dt>Cooking area</dt><dd>{product.cookingArea}</dd></div>
              <div><dt>Fuel</dt><dd>{product.fuel}</dd></div>
              <div><dt>Finish</dt><dd>{product.finish}</dd></div>
            </dl>
          </section>
          <a className="button button-dark detail-cta" href={`mailto:hello@emberandfield.example?subject=${encodeURIComponent(`Enquiry about ${product.name}`)}`}>Ask about this grill <span aria-hidden="true">↗</span></a>
          <p className="detail-note">Sample product information. No checkout or payment is connected.</p>
        </div>
      </div>
    </div>
  );
}

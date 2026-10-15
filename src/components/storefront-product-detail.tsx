import Image from "next/image";
import Link from "next/link";
import { formatPrice } from "@/lib/catalog/mock-data";
import { isSafeProductResourceUrl } from "@/lib/catalog/product-admin";
import { getIndexableProducts } from "@/lib/catalog/seo";

const COOKING_SPECS = new Set(["Burners", "Cooking area", "Fuel"]);

export function StorefrontProductDetail({ slug }: { slug: string }) {
  const product = getIndexableProducts().find((item) => item.slug === slug);
  if (!product) return null;

  const primaryImage = [...product.images].sort((a, b) => a.sortOrder - b.sortOrder)[0] ?? { url: product.image, alt: product.imageAlt, sortOrder: 0 };
  const cookingSpecs = product.specifications.filter((specification) => COOKING_SPECS.has(specification.name));
  const buildSpecs = product.specifications.filter((specification) => !COOKING_SPECS.has(specification.name));
  const applications = product.applications?.filter((application) => application.trim()) ?? [];
  const hasSupply = applications.length > 0 || Boolean(product.minimumOrderQuantity?.trim()) || Boolean(product.leadTime?.trim());
  const resources = product.resources?.filter((resource) => resource.title.trim() && isSafeProductResourceUrl(resource.url)) ?? [];

  return (
    <div className="product-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb">
        <Link href="/">Home</Link><span>/</span>
        <Link href={`/collections/${product.collectionSlug}`}>{product.collection}</Link><span>/</span>
        <span>{product.name}</span>
      </nav>
      <div className="product-detail">
        <div className="product-detail-image" data-reveal>
          <Image src={primaryImage.url} alt={primaryImage.alt} fill priority unoptimized={!primaryImage.url.startsWith("/")} sizes="(max-width: 1023px) 100vw, 55vw" />
          {product.badge && <span className="product-badge">{product.badge}</span>}
        </div>

        <div className="product-detail-copy" data-reveal style={{ "--reveal-delay": "80ms" } as React.CSSProperties}>
          <p className="eyebrow">{product.collection}</p>
          <h1>{product.name}</h1>
          <p className="detail-price">{formatPrice(product.price, product.currency)}</p>
          {product.compareAtPrice && <p className="detail-compare-price">Was {formatPrice(product.compareAtPrice, product.currency)}</p>}
          <p className="detail-summary-lead">{product.shortDescription}</p>
          <p className="detail-summary">{product.description}</p>

          {product.specifications.length > 0 && <section className="detail-specs-section" aria-labelledby="detail-specs-heading">
            <h2 id="detail-specs-heading">Specifications</h2>
            {cookingSpecs.length > 0 && <div className="detail-spec-group">
              <h3>Cooking</h3>
              <dl className="detail-specs">{cookingSpecs.map((specification) => <div key={specification.name}><dt>{specification.name}</dt><dd>{specification.value}</dd></div>)}</dl>
            </div>}
            {buildSpecs.length > 0 && <div className="detail-spec-group">
              <h3>Build</h3>
              <dl className="detail-specs">{buildSpecs.map((specification) => <div key={specification.name}><dt>{specification.name}</dt><dd>{specification.value}</dd></div>)}</dl>
            </div>}
          </section>}

          {hasSupply && <section className="detail-trade-section" aria-labelledby="detail-trade-heading">
            <h2 id="detail-trade-heading">Supply</h2>
            {applications.length > 0 && <div className="detail-applications">
              <h3>Applications</h3>
              <ul>{applications.map((application) => <li key={application}>{application}</li>)}</ul>
            </div>}
            <dl className="detail-trade-facts">
              {product.minimumOrderQuantity?.trim() && <div><dt>Minimum order quantity</dt><dd>{product.minimumOrderQuantity}</dd></div>}
              {product.leadTime?.trim() && <div><dt>Lead time</dt><dd>{product.leadTime}</dd></div>}
            </dl>
          </section>}

          {resources.length > 0 && <section className="detail-resources-section" aria-labelledby="detail-resources-heading">
            <h2 id="detail-resources-heading">Product resources</h2>
            <ul className="detail-resources">{resources.map((resource) => <li key={`${resource.title}-${resource.url}`}>
              <a href={resource.url} target="_blank" rel="noopener noreferrer">{resource.title}<span aria-hidden="true"> ↗</span><span className="sr-only"> (opens in a new tab)</span></a>
              {(resource.fileType || resource.description) && <p>{[resource.fileType, resource.description].filter(Boolean).join(" · ")}</p>}
            </li>)}</ul>
          </section>}

          <a className="button button-primary detail-cta" href={`mailto:hello@emberandfield.example?subject=${encodeURIComponent(`Enquiry about ${product.name}`)}`}>Get in touch</a>
          <p className="detail-note">Sample product information. Confirm specifications, applications and trade terms with the merchant. No checkout or payment is connected.</p>
        </div>
      </div>
    </div>
  );
}

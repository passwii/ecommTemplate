import Link from "next/link";

/**
 * Shown when the storefront catalog API is unreachable. It deliberately does
 * not fall back to sample content, so visitors never see stale or invented data
 * presented as the live catalog.
 */
export function CatalogUnavailable({ heading = "This section is temporarily unavailable." }: { heading?: string }) {
  return (
    <section className="catalog-unavailable" aria-labelledby="catalog-unavailable-title">
      <p className="eyebrow">TEMPORARILY UNAVAILABLE</p>
      <h2 id="catalog-unavailable-title">{heading}</h2>
      <p>We could not load the product catalog just now. Please try again shortly.</p>
      <Link className="text-link" href="/">
        Return to the homepage
      </Link>
    </section>
  );
}

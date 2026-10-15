import Link from "next/link";
import Image from "next/image";
import { StorefrontProductGrid } from "@/components/storefront-product-grid";
import { CategoryLinks } from "@/components/templates/gas-grill/category-links";
import { createPageMetadata } from "@/lib/seo/metadata";

export const metadata = createPageMetadata({
  title: "Outdoor cooking, considered",
  description: "Explore thoughtfully designed gas grills for long lunches, weeknight dinners and gathering outside.",
  pathname: "/",
  image: { url: "/images/grill-hero.jpg", alt: "Outdoor patio lounge with a fire pit beneath string lights and trees" },
});

export default function HomePage() {
  return (
    <>
      <section className="hero">
        <div className="hero-copy" data-reveal>
          <h1>Outdoor cooking,<br /><em>considered.</em></h1>
          <p className="hero-description">Gas grills built for long lunches, weeknight dinners and everything in between.</p>
          <Link className="button button-primary" href="#collection">Explore grills</Link>
        </div>
        <div className="hero-image" data-reveal>
          <Image src="/images/grill-hero.jpg" alt="Outdoor patio lounge with a fire pit beneath string lights and trees" fill priority sizes="(max-width: 1023px) 100vw, 55vw" />
        </div>
      </section>

      <section className="collection-section page-wrap" id="collection">
        <div className="section-heading" data-reveal>
          <h2>A good grill<br />changes the evening.</h2>
          <Link className="text-link" href="/collections/freestanding">Explore grills</Link>
        </div>
        <div data-reveal style={{ "--reveal-delay": "80ms" } as React.CSSProperties}>
          <StorefrontProductGrid featured limit={3} />
        </div>
      </section>

      <section className="category-section page-wrap">
        <div className="section-heading" data-reveal>
          <h2>Find your kind<br />of cooking.</h2>
        </div>
        <div data-reveal style={{ "--reveal-delay": "80ms" } as React.CSSProperties}>
          <CategoryLinks />
        </div>
      </section>

      <section className="feature-section" id="our-approach">
        <div className="feature-image">
          <Image src="/images/grill-detail.jpg" alt="A cook tending a grill during an outdoor gathering" fill sizes="(max-width: 1023px) 100vw, 50vw" />
        </div>
        <div className="feature-copy" data-reveal>
          <p className="eyebrow">Made for more than one season</p>
          <h2>Good materials.<br />Better <em>memories.</em></h2>
          <p>Reliable heat, considered details and durable finishes. The practical things that let you spend less time tending equipment and more time around the table.</p>
          <Link className="text-link" href="/company">Company information</Link>
        </div>
      </section>

      <section className="closing-section page-wrap" data-reveal>
        <h2>There is always<br />room for one more.</h2>
        <Link className="button button-primary" href="/collections/freestanding">Explore grills</Link>
      </section>
    </>
  );
}

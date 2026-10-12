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
        <div className="hero-copy">
          <p className="eyebrow">The art of cooking outside</p>
          <h1>Make a little<br />more <em>of outside.</em></h1>
          <p className="hero-description">Thoughtful grills for long lunches, last-minute dinners and everything in between.</p>
          <Link className="button button-dark" href="#collection">Find your grill <span aria-hidden="true">↗</span></Link>
          <p className="hero-footnote">Designed for the way you gather.</p>
        </div>
        <div className="hero-image">
          <Image src="/images/grill-hero.jpg" alt="Outdoor patio lounge with a fire pit beneath string lights and trees" fill priority sizes="(max-width: 700px) 100vw, (max-width: 1100px) 50vw, 55vw" />
          <span className="hero-image-caption">The Ridge collection, made for open air.</span>
        </div>
        <div className="hero-index"><span>01</span><span>Outdoor cooking, considered</span></div>
      </section>

      <section className="collection-section page-wrap" id="collection">
        <div className="section-heading">
          <div><p className="eyebrow">Find your place at the table</p><h2>A good grill<br />changes the evening.</h2></div>
          <Link className="text-link" href="/collections/freestanding">Shop all grills <span aria-hidden="true">↗</span></Link>
        </div>
        <StorefrontProductGrid featured limit={3} />
      </section>

      <section className="category-section page-wrap">
        <div className="section-heading">
          <div><p className="eyebrow">A place for every kind of outside</p><h2>Find your kind<br />of cooking.</h2></div>
        </div>
        <CategoryLinks />
      </section>

      <section className="feature-section" id="our-approach">
        <div className="feature-image"><Image src="/images/grill-detail.jpg" alt="A cook tending a grill during an outdoor gathering" fill sizes="(max-width: 760px) 100vw, 50vw" /></div>
        <div className="feature-copy">
          <p className="eyebrow">Made for more than one season</p>
          <h2>Good materials.<br />Better <em>memories.</em></h2>
          <p>Reliable heat, considered details and durable finishes. The practical things that let you spend less time tending equipment and more time around the table.</p>
          <Link className="text-link" href="/company">Company information <span aria-hidden="true">↗</span></Link>
        </div>
      </section>

      <section className="closing-section page-wrap">
        <p className="eyebrow">Start with what brings you outside</p>
        <h2>There is always<br />room for one more.</h2>
        <Link className="button button-light" href="/collections/freestanding">Explore the collection <span aria-hidden="true">↗</span></Link>
      </section>
    </>
  );
}

import Link from "next/link";
import { CONTACT_PATH } from "@/lib/contact";

export function SiteFooter() {
  return (
    <footer className="site-footer">
      <Link className="wordmark footer-mark" href="/">EMBER <span>/</span> FIELD</Link>
      <p>Good food takes time. Make room for it.</p>
      <nav className="footer-links" aria-label="Footer navigation">
        <Link href="/collections/freestanding">Shop grills</Link>
        <Link href="/#our-approach">Our approach</Link>
        <Link href={CONTACT_PATH}>Contact</Link>
        <Link href="/admin/products">Product admin</Link>
      </nav>
      <small>© 2026 Ember / Field. Sample storefront content.</small>
    </footer>
  );
}

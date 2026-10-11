import Link from "next/link";

export function SiteFooter() {
  return (
    <footer className="site-footer">
      <Link className="wordmark footer-mark" href="/">EMBER <span>/</span> FIELD</Link>
      <p>Good food takes time. Make room for it.</p>
      <div className="footer-links"><Link href="/collections/freestanding">Shop grills</Link><Link href="/company">Company</Link><Link href="/admin/company">Company content</Link><Link href="/admin/products">Product admin</Link><a href="mailto:hello@emberandfield.example">Get in touch</a></div>
      <small>© 2026 Ember / Field. Sample storefront content.</small>
    </footer>
  );
}

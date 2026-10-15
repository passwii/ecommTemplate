import Link from "next/link";

export function SiteFooter() {
  return (
    <footer className="site-footer">
      <Link className="wordmark footer-mark" href="/">EMBER <span>/</span> FIELD</Link>
      <p>Good food takes time. Make room for it.</p>
      <div className="footer-links">
        <Link href="/collections/freestanding">Explore grills</Link>
        <Link href="/inquiries">Get in touch</Link>
        <Link href="/company">Company</Link>
        <Link href="/admin/company">Company content</Link>
        <Link href="/admin/products">Product admin</Link>
        <a href="mailto:hello@emberandfield.example">Email us</a>
      </div>
      <small>© 2026 Ember / Field. Sample storefront content. Product data is illustrative, nothing here is connected to checkout or payment.</small>
    </footer>
  );
}

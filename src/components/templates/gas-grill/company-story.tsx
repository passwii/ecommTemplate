import Image from "next/image";
import Link from "next/link";
import type { CompanyContent } from "@/lib/catalog/api";

export function CompanyStory({ content }: { content: CompanyContent }) {
  const about = content.about.heading.trim() && content.about.body.trim();
  const manufacturing = content.manufacturing.heading.trim() && content.manufacturing.overview.trim();
  const quality = content.manufacturing.qualityHeading.trim() && content.manufacturing.qualityDetails.trim();
  const credentials = content.credentials.filter((item) => item.name.trim() && item.issuer.trim());
  const caseStudies = content.caseStudies.filter((item) => item.customer.trim() && item.title.trim() && item.details.trim());
  const faqs = content.faqs.filter((item) => item.question.trim() && item.answer.trim());
  const hasContent = about || manufacturing || quality || credentials.length || caseStudies.length || faqs.length;
  const isSample = content.contentStatus === "sample";

  return (
    <div className="company-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>Company</span></nav>
      <header className="company-page-heading" data-reveal>
        <h1>Company information</h1>
        <p>Learn about our company, manufacturing, credentials and customer work.</p>
      </header>
      {isSample && <p className="company-sample-note">This company page currently shows sample content. Verified company, factory, certification and customer information has not been supplied yet.</p>}
      {!hasContent && <section className="company-empty" aria-labelledby="company-empty-title" data-reveal>
        <p className="eyebrow">Not published yet</p>
        <h2 id="company-empty-title">Company details are not available yet.</h2>
        <p>This storefront does not include verified company, factory, certification or customer information. Confirmed details will appear here once the store publishes them.</p>
      </section>}

      {!!about && <section className={`company-about${content.about.imageUrl.trim() && content.about.imageAlt.trim() ? " has-image" : ""}`} data-reveal>
        {content.about.imageUrl.trim() && content.about.imageAlt.trim() && <div className="company-about-image"><Image src={content.about.imageUrl} alt={content.about.imageAlt} fill unoptimized sizes="(max-width: 767px) 100vw, 48vw" /></div>}
        <div className="company-about-copy"><p className="eyebrow">About</p><h2>{content.about.heading}</h2><p>{content.about.body}</p></div>
      </section>}

      {!!manufacturing && <section className="company-detail-section" aria-labelledby="company-manufacturing-title" data-reveal>
        <p className="eyebrow">Making &amp; quality</p><h2 id="company-manufacturing-title">{content.manufacturing.heading}</h2><p>{content.manufacturing.overview}</p>
        {quality && <div className="company-quality"><h3>{content.manufacturing.qualityHeading}</h3><p>{content.manufacturing.qualityDetails}</p></div>}
      </section>}
      {!manufacturing && !!quality && <section className="company-detail-section" aria-labelledby="company-quality-title" data-reveal><p className="eyebrow">Quality</p><h2 id="company-quality-title">{content.manufacturing.qualityHeading}</h2><p>{content.manufacturing.qualityDetails}</p></section>}

      {credentials.length > 0 && <section className="company-list-section" aria-labelledby="company-credentials-title" data-reveal><p className="eyebrow">Credentials</p><h2 id="company-credentials-title">Certifications</h2><div className="company-content-list">{credentials.map((item, index) => <article className="company-content-row" key={`credential-${index}`}><h3>{item.name}</h3><p className="company-content-source">{item.issuer}</p>{item.details.trim() && <p>{item.details}</p>}</article>)}</div></section>}
      {caseStudies.length > 0 && <section className="company-list-section" aria-labelledby="company-cases-title" data-reveal><p className="eyebrow">Customer work</p><h2 id="company-cases-title">Selected examples</h2><div className="company-content-list">{caseStudies.map((item, index) => <article className="company-content-row" key={`case-${index}`}><p className="company-content-source">{item.customer}</p><h3>{item.title}</h3><p>{item.details}</p></article>)}</div></section>}
      {faqs.length > 0 && <section className="company-list-section company-faq" aria-labelledby="company-faq-title" data-reveal><p className="eyebrow">Answers</p><h2 id="company-faq-title">Frequently asked questions</h2><div className="company-content-list">{faqs.map((item, index) => <details className="company-faq-item" key={`faq-${index}`}><summary>{item.question}</summary><p>{item.answer}</p></details>)}</div></section>}
      <section className="company-contact" data-reveal><h2>Have a question?</h2><p>Contact us to discuss the information you need.</p><Link className="button button-primary" href="/inquiries">Get in touch</Link></section>
    </div>
  );
}

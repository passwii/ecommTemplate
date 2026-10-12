"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { createEmptyCompanyContent, getStoredCompanyContent, type CompanyContent } from "@/lib/catalog/company-content";

export function CompanyStory() {
  const [content, setContent] = useState<CompanyContent>(createEmptyCompanyContent);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    setContent(getStoredCompanyContent());
    setLoaded(true);
  }, []);

  const about = content.about.heading.trim() && content.about.body.trim();
  const manufacturing = content.manufacturing.heading.trim() && content.manufacturing.overview.trim();
  const quality = content.manufacturing.qualityHeading.trim() && content.manufacturing.qualityDetails.trim();
  const credentials = content.credentials.filter((item) => item.name.trim() && item.issuer.trim());
  const caseStudies = content.caseStudies.filter((item) => item.customer.trim() && item.title.trim() && item.details.trim());
  const faqs = content.faqs.filter((item) => item.question.trim() && item.answer.trim());
  const hasContent = about || manufacturing || quality || credentials.length || caseStudies.length || faqs.length;

  return (
    <div className="company-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>Company</span></nav>
      <header className="company-page-heading">
        <p className="eyebrow">EMBER / FIELD</p>
        <h1>Company information</h1>
        <p>Learn about our company, manufacturing, credentials and customer work.</p>
      </header>
      {!loaded && <p className="company-loading" aria-live="polite">Loading company information…</p>}
      {loaded && !hasContent && <section className="company-empty" aria-labelledby="company-empty-title">
        <p className="eyebrow">STOREFRONT PREVIEW</p>
        <h2 id="company-empty-title">Company details are not available yet.</h2>
        <p>This sample storefront does not include verified company, factory, certification or customer information. A store administrator can add confirmed details in the local content editor.</p>
        <Link className="text-link" href="/admin/company">Open content preview <span aria-hidden="true">↗</span></Link>
      </section>}

      {loaded && !!about && <section className={`company-about${content.about.imageUrl.trim() && content.about.imageAlt.trim() ? " has-image" : ""}`}>
        {content.about.imageUrl.trim() && content.about.imageAlt.trim() && <div className="company-about-image"><Image src={content.about.imageUrl} alt={content.about.imageAlt} fill unoptimized sizes="(max-width: 760px) 100vw, 48vw" /></div>}
        <div className="company-about-copy"><p className="eyebrow">ABOUT</p><h2>{content.about.heading}</h2><p>{content.about.body}</p></div>
      </section>}

      {loaded && !!manufacturing && <section className="company-detail-section" aria-labelledby="company-manufacturing-title">
        <p className="eyebrow">MAKING & QUALITY</p><h2 id="company-manufacturing-title">{content.manufacturing.heading}</h2><p>{content.manufacturing.overview}</p>
        {quality && <div className="company-quality"><h3>{content.manufacturing.qualityHeading}</h3><p>{content.manufacturing.qualityDetails}</p></div>}
      </section>}
      {loaded && !manufacturing && !!quality && <section className="company-detail-section" aria-labelledby="company-quality-title"><p className="eyebrow">QUALITY</p><h2 id="company-quality-title">{content.manufacturing.qualityHeading}</h2><p>{content.manufacturing.qualityDetails}</p></section>}

      {loaded && credentials.length > 0 && <section className="company-list-section" aria-labelledby="company-credentials-title"><p className="eyebrow">CREDENTIALS</p><h2 id="company-credentials-title">Certifications</h2><div className="company-content-list">{credentials.map((item, index) => <article className="company-content-row" key={`credential-${index}`}><h3>{item.name}</h3><p className="company-content-source">{item.issuer}</p>{item.details.trim() && <p>{item.details}</p>}</article>)}</div></section>}
      {loaded && caseStudies.length > 0 && <section className="company-list-section" aria-labelledby="company-cases-title"><p className="eyebrow">CUSTOMER WORK</p><h2 id="company-cases-title">Selected examples</h2><div className="company-content-list">{caseStudies.map((item, index) => <article className="company-content-row" key={`case-${index}`}><p className="company-content-source">{item.customer}</p><h3>{item.title}</h3><p>{item.details}</p></article>)}</div></section>}
      {loaded && faqs.length > 0 && <section className="company-list-section company-faq" aria-labelledby="company-faq-title"><p className="eyebrow">ANSWERS</p><h2 id="company-faq-title">Frequently asked questions</h2><div className="company-content-list">{faqs.map((item, index) => <details className="company-faq-item" key={`faq-${index}`}><summary>{item.question}</summary><p>{item.answer}</p></details>)}</div></section>}
      <section className="company-contact"><h2>Have a question?</h2><p>Contact us to discuss the information you need.</p><a className="button button-dark" href="mailto:hello@emberandfield.example">Get in touch <span aria-hidden="true">↗</span></a></section>
      <p className="company-storage-note">Sample storefront preview. Company content edits are stored in this browser only and are not published to a server.</p>
    </div>
  );
}

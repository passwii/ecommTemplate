"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import {
  createEmptyCompanyContent,
  getStoredCompanyContent,
  saveStoredCompanyContent,
  type CompanyContent,
} from "@/lib/catalog/company-content";

export function CompanyAdmin() {
  const [content, setContent] = useState<CompanyContent | null>(null);
  const [notice, setNotice] = useState("");

  useEffect(() => setContent(getStoredCompanyContent()), []);

  function save() {
    if (!content) return;
    try {
      saveStoredCompanyContent(content);
      setNotice("Company content saved in this browser.");
    } catch {
      setNotice("Could not save in this browser. Check available storage and try again.");
    }
  }

  function reset() {
    window.localStorage.removeItem("ember-field.storefront.company-content.v1");
    setContent(createEmptyCompanyContent());
    setNotice("Company content cleared from this browser.");
  }

  if (!content) return <div className="admin-loading" aria-live="polite">Loading company content…</div>;

  return (
    <div className="admin-shell">
      <header className="admin-topbar">
        <Link href="/" className="admin-brand">EMBER <span>/</span> FIELD <small>STUDIO</small></Link>
        <nav aria-label="Admin navigation"><Link href="/admin/products">Products</Link><Link href="/admin/company" aria-current="page">Company content</Link><Link href="/">View storefront ↗</Link></nav>
      </header>
      <main className="admin-main company-admin-main">
        <div className="admin-heading">
          <div><p className="admin-eyebrow">STOREFRONT CONTENT · GAS-GRILL TEMPLATE</p><h1>Company & trust</h1><p>Add only facts and assets confirmed by the merchant.</p></div>
          <button className="admin-primary" type="button" onClick={save}>Save content</button>
        </div>
        <div className="admin-callout"><span className="admin-callout-mark">i</span><p>This is a local mock editor: content is saved only in this browser’s localStorage. It does not upload images or sync to a server or other devices. Example prompts below are not published; leave unconfirmed claims blank.</p></div>
        {notice && <p className="admin-notice" role="status">{notice}</p>}

        <section className="company-editor-section" aria-labelledby="company-about-editor">
          <h2 id="company-about-editor">About</h2>
          <p>Company introduction. All fields are optional; the storefront hides this section until the heading and body are filled.</p>
          <label className="admin-field">Section heading<input value={content.about.heading} onChange={(event) => setContent({ ...content, about: { ...content.about, heading: event.target.value } })} placeholder="Example: A little about our company" /></label>
          <label className="admin-field">Introduction<textarea rows={4} value={content.about.body} onChange={(event) => setContent({ ...content, about: { ...content.about, body: event.target.value } })} placeholder="Example only: Add a merchant-approved company introduction." /></label>
          <label className="admin-field">Image URL<input value={content.about.imageUrl} onChange={(event) => setContent({ ...content, about: { ...content.about, imageUrl: event.target.value } })} placeholder="Example: /images/company-workshop.jpg" /></label>
          <label className="admin-field">Image description (alt text)<input value={content.about.imageAlt} onChange={(event) => setContent({ ...content, about: { ...content.about, imageAlt: event.target.value } })} placeholder="Example: Describe the actual image for accessibility" /></label>
          <p className="company-field-note">No image upload is available. Enter a path to an existing public image or a public image URL; add alt text before the image can appear.</p>
        </section>

        <section className="company-editor-section" aria-labelledby="company-manufacturing-editor">
          <h2 id="company-manufacturing-editor">Manufacturing & quality</h2>
          <p>Describe confirmed processes only. Leave fields blank if details are not verified.</p>
          <label className="admin-field">Manufacturing heading<input value={content.manufacturing.heading} onChange={(event) => setContent({ ...content, manufacturing: { ...content.manufacturing, heading: event.target.value } })} placeholder="Example: How our products are made" /></label>
          <label className="admin-field">Manufacturing details<textarea rows={4} value={content.manufacturing.overview} onChange={(event) => setContent({ ...content, manufacturing: { ...content.manufacturing, overview: event.target.value } })} placeholder="Example only: Enter a verified manufacturing description." /></label>
          <label className="admin-field">Quality-check heading<input value={content.manufacturing.qualityHeading} onChange={(event) => setContent({ ...content, manufacturing: { ...content.manufacturing, qualityHeading: event.target.value } })} placeholder="Example: Quality checks" /></label>
          <label className="admin-field">Quality-check details<textarea rows={4} value={content.manufacturing.qualityDetails} onChange={(event) => setContent({ ...content, manufacturing: { ...content.manufacturing, qualityDetails: event.target.value } })} placeholder="Example only: Add checks the merchant can substantiate." /></label>
        </section>

        <section className="company-editor-section" aria-labelledby="company-credentials-editor">
          <div className="company-editor-title"><div><h2 id="company-credentials-editor">Certifications</h2><p>Only add credentials the merchant has confirmed and can substantiate.</p></div><button className="admin-secondary" type="button" onClick={() => setContent({ ...content, credentials: [...content.credentials, { name: "", issuer: "", details: "" }] })}>Add certification</button></div>
          {content.credentials.map((item, index) => <div className="company-repeater" key={`credential-${index}`}>
            <div className="company-repeater-title"><strong>Certification {index + 1}</strong><button type="button" onClick={() => setContent({ ...content, credentials: content.credentials.filter((_, itemIndex) => itemIndex !== index) })}>Remove</button></div>
            <label className="admin-field">Name<input value={item.name} onChange={(event) => setContent({ ...content, credentials: content.credentials.map((entry, itemIndex) => itemIndex === index ? { ...entry, name: event.target.value } : entry) })} placeholder="Example only: Verified certification name" /></label>
            <label className="admin-field">Issuing organization<input value={item.issuer} onChange={(event) => setContent({ ...content, credentials: content.credentials.map((entry, itemIndex) => itemIndex === index ? { ...entry, issuer: event.target.value } : entry) })} placeholder="Example only: Issuing organization" /></label>
            <label className="admin-field">Details<textarea rows={2} value={item.details} onChange={(event) => setContent({ ...content, credentials: content.credentials.map((entry, itemIndex) => itemIndex === index ? { ...entry, details: event.target.value } : entry) })} placeholder="Example only: Scope, validity, or a verified reference" /></label>
          </div>)}
          {!content.credentials.length && <p className="company-editor-empty">No certifications added. Empty entries are not shown on the storefront.</p>}
        </section>

        <section className="company-editor-section" aria-labelledby="company-cases-editor">
          <div className="company-editor-title"><div><h2 id="company-cases-editor">Customer examples</h2><p>Use customer names and project details only with permission.</p></div><button className="admin-secondary" type="button" onClick={() => setContent({ ...content, caseStudies: [...content.caseStudies, { customer: "", title: "", details: "" }] })}>Add example</button></div>
          {content.caseStudies.map((item, index) => <div className="company-repeater" key={`case-${index}`}>
            <div className="company-repeater-title"><strong>Customer example {index + 1}</strong><button type="button" onClick={() => setContent({ ...content, caseStudies: content.caseStudies.filter((_, itemIndex) => itemIndex !== index) })}>Remove</button></div>
            <label className="admin-field">Customer name<input value={item.customer} onChange={(event) => setContent({ ...content, caseStudies: content.caseStudies.map((entry, itemIndex) => itemIndex === index ? { ...entry, customer: event.target.value } : entry) })} placeholder="Example only: Approved customer name" /></label>
            <label className="admin-field">Project title<input value={item.title} onChange={(event) => setContent({ ...content, caseStudies: content.caseStudies.map((entry, itemIndex) => itemIndex === index ? { ...entry, title: event.target.value } : entry) })} placeholder="Example only: Project or application" /></label>
            <label className="admin-field">Details<textarea rows={3} value={item.details} onChange={(event) => setContent({ ...content, caseStudies: content.caseStudies.map((entry, itemIndex) => itemIndex === index ? { ...entry, details: event.target.value } : entry) })} placeholder="Example only: Add a customer-approved summary." /></label>
          </div>)}
          {!content.caseStudies.length && <p className="company-editor-empty">No customer examples added. Empty entries are not shown on the storefront.</p>}
        </section>

        <section className="company-editor-section" aria-labelledby="company-faq-editor">
          <div className="company-editor-title"><div><h2 id="company-faq-editor">Frequently asked questions</h2><p>Answers should reflect the merchant’s current, confirmed policies.</p></div><button className="admin-secondary" type="button" onClick={() => setContent({ ...content, faqs: [...content.faqs, { question: "", answer: "" }] })}>Add question</button></div>
          {content.faqs.map((item, index) => <div className="company-repeater" key={`faq-${index}`}>
            <div className="company-repeater-title"><strong>Question {index + 1}</strong><button type="button" onClick={() => setContent({ ...content, faqs: content.faqs.filter((_, itemIndex) => itemIndex !== index) })}>Remove</button></div>
            <label className="admin-field">Question<input value={item.question} onChange={(event) => setContent({ ...content, faqs: content.faqs.map((entry, itemIndex) => itemIndex === index ? { ...entry, question: event.target.value } : entry) })} placeholder="Example only: A question customers often ask" /></label>
            <label className="admin-field">Answer<textarea rows={3} value={item.answer} onChange={(event) => setContent({ ...content, faqs: content.faqs.map((entry, itemIndex) => itemIndex === index ? { ...entry, answer: event.target.value } : entry) })} placeholder="Example only: Add a verified answer." /></label>
          </div>)}
          {!content.faqs.length && <p className="company-editor-empty">No questions added. Empty entries are not shown on the storefront.</p>}
        </section>

        <div className="company-editor-actions"><button className="admin-secondary" type="button" onClick={reset}>Clear company content</button><button className="admin-primary" type="button" onClick={save}>Save content</button></div>
      </main>
    </div>
  );
}

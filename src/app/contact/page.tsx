import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Contact",
  description: "Get in touch with Ember / Field about outdoor cooking and gas grills.",
};

export default function ContactPage() {
  return (
    <section className="contact-page page-wrap" aria-labelledby="contact-heading">
      <p className="eyebrow">We would be glad to hear from you</p>
      <h1 id="contact-heading">Let’s talk grills.</h1>
      <p className="contact-intro">Contact details and the enquiry form are not available yet. Please check back soon.</p>
      <p className="contact-note">This storefront sample does not yet accept or deliver enquiries.</p>
    </section>
  );
}

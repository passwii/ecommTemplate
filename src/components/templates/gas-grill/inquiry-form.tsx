"use client";

import { useRef, useState } from "react";

const PRIVACY_NOTICE_VERSION = "inquiry-privacy-v1";

type FieldErrors = Record<string, string>;
type Status = "idle" | "submitting" | "success" | "error";

type InquiryFormProps = {
  product?: { id: string; name: string; model?: string };
  sourcePath: string;
};

const fieldOrder = ["contactName", "contactEmail", "contactCompany", "contactPhone", "destinationCountry", "quantity", "customizationRequest", "privacyConsent"];

export function InquiryForm({ product, sourcePath }: InquiryFormProps) {
  const [status, setStatus] = useState<Status>("idle");
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [formError, setFormError] = useState<string>("");
  const summaryRef = useRef<HTMLDivElement>(null);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);

    const quantityRaw = String(data.get("quantity") ?? "").trim();
    const payload = {
      contact: {
        name: String(data.get("name") ?? ""),
        email: String(data.get("email") ?? ""),
        company: String(data.get("company") ?? ""),
        phone: String(data.get("phone") ?? ""),
      },
      destinationCountry: String(data.get("destinationCountry") ?? ""),
      ...(product ? { productId: product.id } : {}),
      ...(quantityRaw ? { quantity: Number(quantityRaw) } : {}),
      customizationRequest: String(data.get("customizationRequest") ?? ""),
      privacyConsent: {
        accepted: data.get("privacyConsent") === "on",
        noticeVersion: PRIVACY_NOTICE_VERSION,
      },
      sourcePath,
      referenceCode: String(data.get("referenceCode") ?? ""),
    };

    setStatus("submitting");
    setFieldErrors({});
    setFormError("");

    let response: Response;
    try {
      response = await fetch("/api/inquiries", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify(payload),
      });
    } catch {
      setStatus("error");
      setFormError("We could not reach the store just now. Please check your connection and try again.");
      return;
    }

    if (response.status === 201 || response.status === 200) {
      setStatus("success");
      form.reset();
      return;
    }

    let body: { error?: { code?: string; message?: string; fields?: FieldErrors } } = {};
    try {
      body = await response.json();
    } catch {
      body = {};
    }

    setStatus("error");
    setFieldErrors(body.error?.fields ?? {});

    if (response.status === 429) {
      setFormError("You have sent several enquiries already. Please try again a little later.");
    } else if (response.status === 404) {
      setFormError("That product is no longer available for enquiries. Please choose another product.");
    } else if (response.status >= 500 || response.status === 503) {
      setFormError("Something went wrong on our side. Please try again shortly.");
    } else {
      setFormError("Please check the highlighted fields and try again.");
    }
    summaryRef.current?.focus();
  }

  if (status === "success") {
    return (
      <div className="inquiry-success" role="status">
        <p className="eyebrow">ENQUIRY SENT</p>
        <h2>Thank you — we have your enquiry.</h2>
        <p>Our team will review your requirements and reply to the email address you provided.</p>
        <button className="button button-dark" type="button" onClick={() => setStatus("idle")}>
          Send another enquiry
        </button>
      </div>
    );
  }

  const orderedErrors = fieldOrder.filter((field) => fieldErrors[field]);

  return (
    <form className="inquiry-form" onSubmit={handleSubmit} noValidate>
      {status === "error" && (
        <div className="inquiry-error-summary" role="alert" tabIndex={-1} ref={summaryRef}>
          <p>{formError || "Please check the highlighted fields and try again."}</p>
          {orderedErrors.length > 0 && (
            <ul>
              {orderedErrors.map((field) => (
                <li key={field}>{fieldErrors[field]}</li>
              ))}
            </ul>
          )}
        </div>
      )}

      {product && (
        <p className="inquiry-product">
          Enquiry about <strong>{product.name}</strong>
          {product.model && product.model !== product.name ? ` (${product.model})` : ""}
        </p>
      )}

      <div className="inquiry-field">
        <label htmlFor="inquiry-name">Your name</label>
        <input id="inquiry-name" name="name" type="text" required autoComplete="name"
          aria-invalid={Boolean(fieldErrors.contactName)} aria-describedby={fieldErrors.contactName ? "inquiry-name-error" : undefined} />
        {fieldErrors.contactName && <p className="inquiry-field-error" id="inquiry-name-error">{fieldErrors.contactName}</p>}
      </div>

      <div className="inquiry-field">
        <label htmlFor="inquiry-email">Email</label>
        <input id="inquiry-email" name="email" type="email" required autoComplete="email"
          aria-invalid={Boolean(fieldErrors.contactEmail)} aria-describedby={fieldErrors.contactEmail ? "inquiry-email-error" : undefined} />
        {fieldErrors.contactEmail && <p className="inquiry-field-error" id="inquiry-email-error">{fieldErrors.contactEmail}</p>}
      </div>

      <div className="inquiry-field-row">
        <div className="inquiry-field">
          <label htmlFor="inquiry-company">Company <span className="inquiry-optional">optional</span></label>
          <input id="inquiry-company" name="company" type="text" autoComplete="organization"
            aria-invalid={Boolean(fieldErrors.contactCompany)} />
          {fieldErrors.contactCompany && <p className="inquiry-field-error">{fieldErrors.contactCompany}</p>}
        </div>
        <div className="inquiry-field">
          <label htmlFor="inquiry-phone">Phone / WhatsApp <span className="inquiry-optional">optional</span></label>
          <input id="inquiry-phone" name="phone" type="tel" autoComplete="tel"
            aria-invalid={Boolean(fieldErrors.contactPhone)} />
          {fieldErrors.contactPhone && <p className="inquiry-field-error">{fieldErrors.contactPhone}</p>}
        </div>
      </div>

      <div className="inquiry-field-row">
        <div className="inquiry-field">
          <label htmlFor="inquiry-country">Destination country <span className="inquiry-optional">optional</span></label>
          <input id="inquiry-country" name="destinationCountry" type="text" maxLength={2} placeholder="US" autoComplete="country"
            aria-invalid={Boolean(fieldErrors.destinationCountry)} aria-describedby="inquiry-country-hint" />
          <p className="inquiry-hint" id="inquiry-country-hint">Two-letter country code.</p>
          {fieldErrors.destinationCountry && <p className="inquiry-field-error">{fieldErrors.destinationCountry}</p>}
        </div>
        <div className="inquiry-field">
          <label htmlFor="inquiry-quantity">Estimated quantity <span className="inquiry-optional">optional</span></label>
          <input id="inquiry-quantity" name="quantity" type="number" min={1} step={1}
            aria-invalid={Boolean(fieldErrors.quantity)} />
          {fieldErrors.quantity && <p className="inquiry-field-error">{fieldErrors.quantity}</p>}
        </div>
      </div>

      <div className="inquiry-field">
        <label htmlFor="inquiry-message">Your requirements <span className="inquiry-optional">optional</span></label>
        <textarea id="inquiry-message" name="customizationRequest" rows={5}
          aria-invalid={Boolean(fieldErrors.customizationRequest)} />
        {fieldErrors.customizationRequest && <p className="inquiry-field-error">{fieldErrors.customizationRequest}</p>}
      </div>

      {/* Honeypot: hidden from visitors and assistive technology. */}
      <div className="inquiry-honeypot" aria-hidden="true">
        <label htmlFor="inquiry-reference">Reference code</label>
        <input id="inquiry-reference" name="referenceCode" type="text" tabIndex={-1} autoComplete="off" />
      </div>

      <div className="inquiry-consent">
        <input id="inquiry-consent" name="privacyConsent" type="checkbox" required
          aria-invalid={Boolean(fieldErrors.privacyConsent)} aria-describedby="inquiry-privacy" />
        <label htmlFor="inquiry-consent">
          I agree that my details may be used to respond to this enquiry.
        </label>
      </div>
      <p className="inquiry-privacy" id="inquiry-privacy">
        We use your contact details only to answer this enquiry and to discuss your requirements. We do not sell your data.
        Ask us at any time to delete it.
      </p>
      {fieldErrors.privacyConsent && <p className="inquiry-field-error">{fieldErrors.privacyConsent}</p>}

      <button className="button button-dark inquiry-submit" type="submit" disabled={status === "submitting"}>
        {status === "submitting" ? "Sending…" : "Send enquiry"}
      </button>
      <p className="inquiry-note" aria-live="polite">
        {status === "submitting" ? "Sending your enquiry…" : "We usually reply within one business day."}
      </p>
    </form>
  );
}

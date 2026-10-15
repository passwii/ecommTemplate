import type { Metadata } from "next";
import Link from "next/link";
import { notFound, redirect } from "next/navigation";
import { AdminInquiryActions } from "@/components/admin/inquiry-actions";
import { fetchAdminInquiry, formatInquiryTimestamp, statusLabel } from "@/lib/admin/api";

export const metadata: Metadata = {
  title: "Inquiry",
  robots: { index: false, follow: false },
};

type Props = { params: Promise<{ id: string }> };

export default async function AdminInquiryDetailPage({ params }: Props) {
  const { id } = await params;
  const result = await fetchAdminInquiry(id);

  if (!result.ok && result.reason === "unauthorized") redirect("/admin/login");
  if (!result.ok && result.reason === "not_found") notFound();

  if (!result.ok) {
    return (
      <div className="admin-page page-wrap">
        <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><Link href="/admin/inquiries">Inquiries</Link><span>/</span><span>Inquiry</span></nav>
        <p className="admin-empty" role="alert">This inquiry could not be loaded. Please try again.</p>
      </div>
    );
  }

  const inquiry = result.data;
  const activities = inquiry.activities ?? [];

  return (
    <div className="admin-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb">
        <Link href="/">Home</Link><span>/</span><Link href="/admin/inquiries">Inquiries</Link><span>/</span><span>{inquiry.contact.name}</span>
      </nav>

      <header className="admin-heading">
        <p className="eyebrow">INQUIRY · {statusLabel(inquiry.status)}</p>
        <h1>{inquiry.contact.name}</h1>
        <p>
          Received {formatInquiryTimestamp(inquiry.createdAt)}
          {inquiry.product.name ? ` · ${inquiry.product.name}` : ""}
        </p>
      </header>

      <div className="admin-detail-layout">
        <section className="admin-card" aria-labelledby="admin-submission-heading">
          <h2 id="admin-submission-heading">Submitted details</h2>
          <dl className="admin-facts">
            <div><dt>Email</dt><dd><a href={`mailto:${inquiry.contact.email}`}>{inquiry.contact.email}</a></dd></div>
            <div><dt>Company</dt><dd>{inquiry.contact.company || "—"}</dd></div>
            <div><dt>Phone</dt><dd>{inquiry.contact.phone || "—"}</dd></div>
            <div><dt>Destination country</dt><dd>{inquiry.destinationCountry || "—"}</dd></div>
            <div><dt>Product</dt><dd>{inquiry.product.name || "—"}{inquiry.product.model && inquiry.product.model !== inquiry.product.name ? ` (${inquiry.product.model})` : ""}</dd></div>
            <div><dt>Quantity</dt><dd>{inquiry.quantity ?? "—"}</dd></div>
            <div><dt>Source page</dt><dd>{inquiry.sourcePath || "—"}</dd></div>
            <div><dt>Privacy notice</dt><dd>{inquiry.privacyNoticeVersion} · {formatInquiryTimestamp(inquiry.privacyConsentedAt)}</dd></div>
          </dl>
          {inquiry.customizationRequest && (
            <div className="admin-message">
              <h3>Requirements</h3>
              <p>{inquiry.customizationRequest}</p>
            </div>
          )}
        </section>

        <AdminInquiryActions
          inquiryId={inquiry.id}
          status={inquiry.status}
          owner={inquiry.owner}
          nextAction={inquiry.nextAction}
          nextActionDueAt={inquiry.nextActionDueAt}
          lossReason={inquiry.lossReason}
        />
      </div>

      <section className="admin-card admin-history" aria-labelledby="admin-history-heading">
        <h2 id="admin-history-heading">History</h2>
        {activities.length === 0 ? (
          <p className="admin-empty">No activity recorded yet.</p>
        ) : (
          <ol className="admin-timeline">
            {activities.map((activity) => (
              <li key={activity.id}>
                <p className="admin-timeline-meta">{formatInquiryTimestamp(activity.createdAt)} · {activity.type} · {activity.actor}</p>
                <p>{activity.body}</p>
              </li>
            ))}
          </ol>
        )}
      </section>

      <p className="admin-note">Local MVP administration. Status changes and notes are recorded against this store&apos;s database.</p>
    </div>
  );
}

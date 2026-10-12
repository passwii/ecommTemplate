"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

type AdminInquiryActionsProps = {
  inquiryId: string;
  status: string;
  owner: string;
  nextAction: string;
  nextActionDueAt: string | null;
  lossReason: string;
};

const STATUSES = ["new", "in_progress", "quoted", "won", "lost", "on_hold"] as const;

function toDateInput(value: string | null): string {
  if (!value) return "";
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return "";
  return parsed.toISOString().slice(0, 10);
}

export function AdminInquiryActions({ inquiryId, status, owner, nextAction, nextActionDueAt, lossReason }: AdminInquiryActionsProps) {
  const router = useRouter();
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  async function send(path: string, method: string, body: unknown) {
    const response = await fetch(`/api/admin/${path}`, {
      method,
      headers: { "content-type": "application/json" },
      body: JSON.stringify(body),
    });
    return response;
  }

  async function handleDetails(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    setSaving(true);
    setMessage("");
    setError("");

    const dueDate = String(data.get("nextActionDueAt") ?? "").trim();
    try {
      const response = await send(`inquiries/${encodeURIComponent(inquiryId)}`, "PATCH", {
        status: String(data.get("status") ?? status),
        owner: String(data.get("owner") ?? ""),
        nextAction: String(data.get("nextAction") ?? ""),
        lossReason: String(data.get("lossReason") ?? ""),
        nextActionDueAt: dueDate ? new Date(`${dueDate}T00:00:00Z`).toISOString() : "",
      });

      if (response.ok) {
        setMessage("Saved.");
        router.refresh();
      } else if (response.status === 401) {
        router.replace("/admin/login");
      } else {
        setError("The changes could not be saved. Please check the values and try again.");
      }
    } catch {
      setError("The store could not be reached. Please try again.");
    } finally {
      setSaving(false);
    }
  }

  async function handleNote(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    setSaving(true);
    setMessage("");
    setError("");

    try {
      const response = await send(`inquiries/${encodeURIComponent(inquiryId)}/activities`, "POST", {
        type: String(data.get("type") ?? "note"),
        body: String(data.get("body") ?? ""),
      });

      if (response.status === 201) {
        form.reset();
        setMessage("Activity recorded.");
        router.refresh();
      } else if (response.status === 401) {
        router.replace("/admin/login");
      } else {
        setError("The activity could not be recorded. Enter a note and try again.");
      }
    } catch {
      setError("The store could not be reached. Please try again.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="admin-actions">
      {message && <p className="admin-success" role="status">{message}</p>}
      {error && <p className="admin-error" role="alert">{error}</p>}

      <form className="admin-card" onSubmit={handleDetails}>
        <h2>Follow-up</h2>
        <div className="admin-field">
          <label htmlFor="admin-status">Status</label>
          <select id="admin-status" name="status" defaultValue={status}>
            {STATUSES.map((value) => <option key={value} value={value}>{value}</option>)}
          </select>
        </div>
        <div className="admin-field">
          <label htmlFor="admin-owner">Owner</label>
          <input id="admin-owner" name="owner" type="text" defaultValue={owner} />
        </div>
        <div className="admin-field">
          <label htmlFor="admin-next-action">Next action</label>
          <input id="admin-next-action" name="nextAction" type="text" defaultValue={nextAction} placeholder="e.g. Send pricing" />
        </div>
        <div className="admin-field">
          <label htmlFor="admin-next-action-due">Next action due</label>
          <input id="admin-next-action-due" name="nextActionDueAt" type="date" defaultValue={toDateInput(nextActionDueAt)} />
        </div>
        <div className="admin-field">
          <label htmlFor="admin-loss-reason">Loss reason</label>
          <input id="admin-loss-reason" name="lossReason" type="text" defaultValue={lossReason} />
        </div>
        <button className="button button-dark" type="submit" disabled={saving}>Save changes</button>
      </form>

      <form className="admin-card" onSubmit={handleNote}>
        <h2>Add activity</h2>
        <div className="admin-field">
          <label htmlFor="admin-activity-type">Type</label>
          <select id="admin-activity-type" name="type" defaultValue="note">
            <option value="note">Note</option>
            <option value="contact">Contact</option>
          </select>
        </div>
        <div className="admin-field">
          <label htmlFor="admin-activity-body">Details</label>
          <textarea id="admin-activity-body" name="body" rows={4} required />
        </div>
        <button className="button button-dark" type="submit" disabled={saving}>Record activity</button>
      </form>

      <form action="/api/admin/sessions" method="post" onSubmit={async (event) => {
        event.preventDefault();
        await fetch("/api/admin/sessions", { method: "DELETE" });
        router.replace("/admin/login");
        router.refresh();
      }}>
        <button className="text-link admin-signout" type="submit">Sign out</button>
      </form>
    </div>
  );
}

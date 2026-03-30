import api from "./client";

export async function exportRankings() {
  const res = await api.get("/export/rankings", { responseType: "blob" });

  // Extract filename from Content-Disposition header
  const disposition = res.headers["content-disposition"] || "";
  const match = disposition.match(/filename="?([^";\s]+)"?/);
  const filename = match?.[1] || "rankings.xlsx";

  // Trigger browser download
  const url = window.URL.createObjectURL(new Blob([res.data]));
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  a.remove();
}

import type { Attendance } from "@/types";

export function StatusBadge({ status }: { status: Attendance["status"] }) {
  return (
    <span className={`status-badge ${status}`}>
      <span />
      {status === "present" ? "حاضر" : "لم يحضر"}
    </span>
  );
}

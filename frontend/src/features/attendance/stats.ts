import type { Attendance } from "@/types";

export function computeStats(people: Attendance[]) {
  const total = people.length;
  const present = people.filter((p) => p.status === "present").length;
  return {
    total,
    present,
    absent: total - present,
    rate: total ? Math.round((present / total) * 100) : 0,
  };
}
export type AttendanceStats = ReturnType<typeof computeStats>;

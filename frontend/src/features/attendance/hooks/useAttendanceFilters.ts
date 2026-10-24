import { useMemo, useState } from "react";
import type { Attendance } from "@/types";
import type { StatusFilter } from "../types";

export function useAttendanceFilters(people: Attendance[]) {
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState<StatusFilter>("all");

  const filtered = useMemo(() => {
    const term = search.trim().toLocaleLowerCase();
    return people.filter((p) => {
      const matchesStatus = status === "all" || p.status === status;
      const haystack = [p.name, p.militaryNumber, p.phoneNumber, p.location]
        .filter(Boolean)
        .join(" ")
        .toLocaleLowerCase();
      return matchesStatus && (!term || haystack.includes(term));
    });
  }, [people, search, status]);

  const reset = () => {
    setSearch("");
    setStatus("all");
  };

  return { search, setSearch, status, setStatus, filtered, reset };
}

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useBatchAttendance(batchId: number) {
  const query = useQuery({
    queryKey: ["attendances"],
    queryFn: api.attendances,
    refetchInterval: 4000,
  });
  const people = useMemo(
    () => (query.data ?? []).filter((a) => a.batchId === batchId),
    [query.data, batchId],
  );
  return { query, people };
}

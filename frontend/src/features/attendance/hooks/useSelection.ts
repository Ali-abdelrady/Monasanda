import { useState } from "react";
import type { Attendance } from "@/types";

export function useSelection(visible: Attendance[]) {
  const [selected, setSelected] = useState<number[]>([]);

  const toggle = (id: number) =>
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id],
    );

  const allVisibleSelected =
    visible.length > 0 && visible.every((p) => selected.includes(p.id));

  const toggleAllVisible = () => {
    const ids = visible.map((p) => p.id);
    setSelected((prev) =>
      ids.length && ids.every((id) => prev.includes(id))
        ? prev.filter((id) => !ids.includes(id))
        : [...new Set([...prev, ...ids])],
    );
  };

  const drop = (id: number) => setSelected((prev) => prev.filter((v) => v !== id));
  const clear = () => setSelected([]);

  return { selected, toggle, toggleAllVisible, allVisibleSelected, drop, clear };
}

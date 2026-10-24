import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { parseAttendanceFile, type ParsedRow } from "@/lib/excel";
import { errorMessage } from "@/lib/errors";
import { useToast } from "@/context/ToastContext";
import type { Attendance, VacationBatch } from "@/types";

const CONCURRENCY = 4;

export function useImport(
  batch: VacationBatch,
  existing: Attendance[],
  onFinished: () => void,
) {
  const qc = useQueryClient();
  const { notify } = useToast();
  const [rows, setRows] = useState<ParsedRow[]>([]);
  const [fileName, setFileName] = useState("");
  const [busy, setBusy] = useState(false);
  const [progress, setProgress] = useState({ done: 0, total: 0 });

  const validRows = rows.filter((r) => r.errors.length === 0);
  const invalidCount = rows.length - validRows.length;

  const reset = () => {
    setFileName("");
    setRows([]);
  };

  const readFile = async (file: File) => {
    setFileName(file.name);
    setRows([]);
    try {
      const parsed = await parseAttendanceFile(file, batch.id);
      setRows(parsed);
      if (!parsed.length) notify("لم يتم العثور على صفوف بيانات في الملف", "error");
    } catch (e) {
      notify(errorMessage(e, "تعذر قراءة ملف Excel"), "error");
    }
  };

  const run = async () => {
    if (!validRows.length) {
      notify("لا توجد صفوف سليمة للاستيراد", "error");
      return;
    }
    setBusy(true);
    setProgress({ done: 0, total: validRows.length });

    const known = new Set(existing.map((p) => p.militaryNumber).filter(Boolean));
    let success = 0,
      skipped = 0,
      failed = 0;

    const queue = validRows.filter((row) => {
      const id = row.input.militaryNumber;
      if (id && known.has(id)) {
        skipped++;
        return false;
      }
      if (id) known.add(id);
      return true;
    });

    let cursor = 0;
    const worker = async () => {
      while (cursor < queue.length) {
        const row = queue[cursor++];
        try {
          await api.createAttendance(row.input);
          success++;
        } catch {
          failed++;
        }
        setProgress((p) => ({ ...p, done: p.done + 1 }));
      }
    };
    await Promise.all(
      Array.from({ length: Math.min(CONCURRENCY, queue.length) }, () => worker()),
    );

    await qc.invalidateQueries({ queryKey: ["attendances"] });
    setBusy(false);
    reset();
    onFinished();
    notify(
      `اكتمل الاستيراد: ${success} ناجح، ${skipped} مكرر تم تخطيه، ${failed} فشل`,
      failed ? "error" : "success",
    );
  };

  return { rows, fileName, busy, progress, validRows, invalidCount, readFile, run, reset };
}

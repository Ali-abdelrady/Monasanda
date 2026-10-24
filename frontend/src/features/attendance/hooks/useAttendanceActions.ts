import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import { useToast } from "@/context/ToastContext";
import type { Attendance, AttendanceInput } from "@/types";

/** All server-side mutations for attendance rows + their user feedback. */
export function useAttendanceActions() {
  const qc = useQueryClient();
  const { notify } = useToast();
  const [bulkBusy, setBulkBusy] = useState(false);
  const refresh = () => qc.invalidateQueries({ queryKey: ["attendances"] });

  const toggle = async (person: Attendance) => {
    const input: AttendanceInput = {
      batchId: person.batchId,
      name: person.name,
      location: person.location,
      phoneNumber: person.phoneNumber,
      militaryNumber: person.militaryNumber,
      status: person.status === "present" ? "absent" : "present",
      notes: person.notes,
      startDate: person.startDate,
      endDate: person.endDate,
    };
    try {
      await api.updateAttendance(person.id, input);
      await refresh();
      notify(person.status === "present" ? "تم إلغاء التمام" : "تم تسجيل الحضور بنجاح");
    } catch (e) {
      notify(errorMessage(e, "تعذر تحديث الحضور"), "error");
    }
  };

  /** Returns true when the delete succeeded. */
  const remove = async (person: Attendance) => {
    try {
      await api.deleteAttendance(person.id);
      refresh();
      notify("تم حذف سجل الفرد");
      return true;
    } catch (e) {
      notify(errorMessage(e, "تعذر حذف السجل"), "error");
      return false;
    }
  };

  const removeMany = async (ids: number[]) => {
    setBulkBusy(true);
    let ok = 0;
    let failed = 0;
    for (const id of ids) {
      try {
        await api.deleteAttendance(id);
        ok++;
      } catch {
        failed++;
      }
    }
    await refresh();
    setBulkBusy(false);
    notify(
      failed
        ? `تم حذف ${ok} سجل، وفشل حذف ${failed} سجل`
        : `تم حذف ${ok} سجل بنجاح`,
      failed ? "error" : "success",
    );
  };

  return { toggle, remove, removeMany, bulkBusy };
}

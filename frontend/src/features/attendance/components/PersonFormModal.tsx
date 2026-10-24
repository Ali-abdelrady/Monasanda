import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { dateInput, toApiDate } from "@/lib/utils";
import { errorMessage } from "@/lib/errors";
import { useToast } from "@/context/ToastContext";
import { Modal } from "@/components/ui/Modal";
import { Button } from "@/components/ui/Button";
import type { Attendance, AttendanceInput, VacationBatch } from "@/types";

interface Props {
  batch: VacationBatch;
  /** null = add mode */
  person: Attendance | null;
  onClose: () => void;
}

function initialForm(batch: VacationBatch, person: Attendance | null) {
  return person
    ? {
        name: person.name ?? "",
        location: person.location ?? "",
        militaryNumber: person.militaryNumber ?? "",
        phoneNumber: person.phoneNumber ?? "",
        startDate: dateInput(person.startDate),
        endDate: dateInput(person.endDate),
        notes: person.notes ?? "",
      }
    : {
        name: "",
        location: "",
        militaryNumber: "",
        phoneNumber: "",
        startDate: dateInput(batch.startDate),
        endDate: dateInput(batch.endDate),
        notes: "",
      };
}

export function PersonFormModal({ batch, person, onClose }: Props) {
  const qc = useQueryClient();
  const { notify } = useToast();
  const [form, setForm] = useState(() => initialForm(batch, person));
  const set = (key: keyof typeof form) => (value: string) =>
    setForm((f) => ({ ...f, [key]: value }));

  const mutation = useMutation({
    mutationFn: async () => {
      if (!form.name.trim()) throw new Error("اسم الفرد مطلوب");
      if (form.phoneNumber.trim().length > 12)
        throw new Error("رقم التليفون يجب ألا يتجاوز 12 خانة");
      if (form.startDate && form.endDate && form.endDate < form.startDate)
        throw new Error("تاريخ نهاية الإجازة قبل تاريخ البداية");

      const input: AttendanceInput = {
        batchId: batch.id,
        name: form.name.trim(),
        location: form.location.trim() || null,
        militaryNumber: form.militaryNumber.trim() || null,
        phoneNumber: form.phoneNumber.trim() || null,
        startDate: form.startDate ? toApiDate(form.startDate) : null,
        endDate: form.endDate ? toApiDate(form.endDate, true) : null,
        notes: form.notes.trim() || null,
        status: person?.status ?? "absent",
      };
      return person ? api.updateAttendance(person.id, input) : api.createAttendance(input);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["attendances"] });
      notify(person ? "تم تحديث بيانات الفرد" : "تمت إضافة الفرد بنجاح");
      onClose();
    },
    onError: (e) => notify(errorMessage(e, "تعذر حفظ بيانات الفرد"), "error"),
  });

  return (
    <Modal title={person ? "تعديل بيانات الفرد" : "إضافة فرد جديد"} onClose={onClose} wide>
      <form
        className="form-body"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate();
        }}
      >
        <p className="form-intro">
          الحقول المعلّمة بنجمة مطلوبة. رقم التليفون اختياري ويُحفظ كنص للحفاظ على الصفر في
          البداية.
        </p>
        <div className="form-two-cols">
          <div>
            <label className="field-label">
              الاسم <span>*</span>
            </label>
            <input
              className="text-input"
              value={form.name}
              onChange={(e) => set("name")(e.target.value)}
              required
              placeholder="الاسم بالكامل"
            />
          </div>
          <div>
            <label className="field-label">الرقم العسكري</label>
            <input
              className="text-input"
              value={form.militaryNumber}
              onChange={(e) => set("militaryNumber")(e.target.value)}
              placeholder="الرقم العسكري"
              dir="ltr"
            />
          </div>
          <div>
            <label className="field-label">رقم التليفون</label>
            <input
              className="text-input"
              value={form.phoneNumber}
              onChange={(e) => set("phoneNumber")(e.target.value)}
              placeholder="01xxxxxxxxx"
              maxLength={12}
              inputMode="tel"
              dir="ltr"
            />
          </div>
          <div>
            <label className="field-label">السكن</label>
            <input
              className="text-input"
              value={form.location}
              onChange={(e) => set("location")(e.target.value)}
              placeholder="مكان السكن"
            />
          </div>
          <div>
            <label className="field-label">من</label>
            <input
              className="text-input"
              type="date"
              value={form.startDate}
              onChange={(e) => set("startDate")(e.target.value)}
            />
          </div>
          <div>
            <label className="field-label">إلى</label>
            <input
              className="text-input"
              type="date"
              value={form.endDate}
              min={form.startDate}
              onChange={(e) => set("endDate")(e.target.value)}
            />
          </div>
        </div>
        <label className="field-label">ملاحظات</label>
        <textarea
          className="text-input textarea"
          value={form.notes}
          onChange={(e) => set("notes")(e.target.value)}
          placeholder="أي ملاحظات إضافية..."
          rows={3}
        />
        <div className="modal-actions">
          <Button variant="secondary" onClick={onClose}>
            إلغاء
          </Button>
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "جاري الحفظ..." : person ? "حفظ التعديلات" : "إضافة الفرد"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

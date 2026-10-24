import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { dateInput, toApiDate } from "@/lib/utils";
import { errorMessage } from "@/lib/errors";
import { useToast } from "@/context/ToastContext";
import { Modal } from "@/components/ui/Modal";
import { Button } from "@/components/ui/Button";
import type { VacationBatch } from "@/types";

interface Props {
  /** null = create mode, otherwise edit mode */
  batch: VacationBatch | null;
  onClose: () => void;
  onSaved: (batch: VacationBatch) => void;
}

export function BatchFormModal({ batch, onClose, onSaved }: Props) {
  const qc = useQueryClient();
  const { notify } = useToast();
  // Form state lives here: the modal mounts fresh every time it opens,
  // so initial values come straight from props — no manual reset needed.
  const [name, setName] = useState(batch?.name ?? "");
  const [start, setStart] = useState(batch ? dateInput(batch.startDate) : "");
  const [end, setEnd] = useState(batch ? dateInput(batch.endDate) : "");

  const mutation = useMutation({
    mutationFn: async () => {
      if (!name.trim() || !start || !end)
        throw new Error("من فضلك أكمل اسم الدفعة والتواريخ");
      if (end < start) throw new Error("تاريخ النهاية يجب أن يكون بعد تاريخ البداية");
      const input = {
        name: name.trim(),
        startDate: toApiDate(start),
        endDate: toApiDate(end, true),
      };
      return batch ? api.updateBatch(batch.id, input) : api.createBatch(input);
    },
    onSuccess: (saved) => {
      qc.invalidateQueries({ queryKey: ["batches"] });
      notify(batch ? "تم تحديث بيانات الدفعة" : "تم إنشاء الدفعة بنجاح");
      onSaved(saved);
      onClose();
    },
    onError: (e) => notify(errorMessage(e, "تعذر حفظ الدفعة"), "error"),
  });

  return (
    <Modal title={batch ? "تعديل بيانات الدفعة" : "إنشاء دفعة جديدة"} onClose={onClose}>
      <form
        className="form-body"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate();
        }}
      >
        <p className="form-intro">
          {batch
            ? "عدّل اسم الدفعة وفترة الإجازات."
            : "أدخل بيانات الدفعة، وبعدها تقدر تستورد كشف الأفراد."}
        </p>
        <label className="field-label">
          اسم الدفعة <span>*</span>
        </label>
        <input
          className="text-input"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="مثال: دفعة إجازات أكتوبر"
          required
        />
        <div className="form-two-cols">
          <div>
            <label className="field-label">
              من <span>*</span>
            </label>
            <input
              className="text-input"
              type="date"
              value={start}
              onChange={(e) => setStart(e.target.value)}
              required
            />
          </div>
          <div>
            <label className="field-label">
              إلى <span>*</span>
            </label>
            <input
              className="text-input"
              type="date"
              value={end}
              min={start}
              onChange={(e) => setEnd(e.target.value)}
              required
            />
          </div>
        </div>
        <div className="modal-actions">
          <Button variant="secondary" onClick={onClose}>
            إلغاء
          </Button>
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending
              ? "جاري الحفظ..."
              : batch
                ? "حفظ التعديلات"
                : "إنشاء الدفعة"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

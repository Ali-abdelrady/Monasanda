import { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {
  CircleAlert,
  FileSpreadsheet,
  Plus,
  RefreshCw,
  ShieldCheck,
} from "lucide-react";
import { api } from "@/lib/api";
import { errorMessage } from "@/lib/errors";
import { useToast } from "@/context/ToastContext";
import { Button } from "@/components/ui/Button";
import { ConfirmDeleteModal } from "@/components/ui/ConfirmDeleteModal";
import type { VacationBatch } from "@/types";
import { BatchCard } from "./BatchCard";
import { BatchFormModal } from "./BatchFormModal";

type FormState = { batch: VacationBatch | null } | null;

export function BatchesPage({
  onOpen,
}: {
  onOpen: (batch: VacationBatch) => void;
}) {
  const qc = useQueryClient();
  const { notify } = useToast();
  const batchesQuery = useQuery({
    queryKey: ["batches"],
    queryFn: api.batches,
  });
  const [form, setForm] = useState<FormState>(null);
  const [deleting, setDeleting] = useState<VacationBatch | null>(null);

  const confirmDelete = async () => {
    if (!deleting) return;
    try {
      await api.deleteBatch(deleting.id);
      qc.invalidateQueries({ queryKey: ["batches"] });
      notify("تم حذف الدفعة");
    } catch (e) {
      notify(errorMessage(e, "تعذر حذف الدفعة"), "error");
    } finally {
      setDeleting(null);
    }
  };

  const batches = batchesQuery.data;

  return (
    <>
      <div className="page-heading">
        <div>
          <div className="eyebrow">مساحة العمل</div>
          <h1>دفعات الإجازات</h1>
          <p>إدارة الكشوف ومتابعة عودة الأفراد من مكان واحد.</p>
        </div>
        <Button onClick={() => setForm({ batch: null })}>
          <Plus size={18} /> إنشاء دفعة جديدة
        </Button>
      </div>

      <div className="section-line">
        <div>
          <h2>كل الدفعات</h2>
          <span className="count-pill">{batches?.length ?? 0}</span>
        </div>
        <button
          className="subtle-action"
          onClick={() => qc.invalidateQueries({ queryKey: ["batches"] })}
        >
          <RefreshCw size={15} /> تحديث
        </button>
      </div>

      {batchesQuery.isPending ? (
        <div className="batch-grid">
          {[1, 2, 3].map((n) => (
            <div className="batch-skeleton" key={n} />
          ))}
        </div>
      ) : batchesQuery.isError ? (
        <div className="empty-state">
          <CircleAlert size={28} />
          <h3>تعذر تحميل الدفعات</h3>
          <p>
            {errorMessage(
              batchesQuery.error,
              "تحقق من اتصال الخادم وحاول مرة أخرى.",
            )}
          </p>
          <Button variant="secondary" onClick={() => batchesQuery.refetch()}>
            <RefreshCw size={16} /> إعادة المحاولة
          </Button>
        </div>
      ) : batches?.length ? (
        <div className="batch-grid">
          {batches.map((batch) => (
            <BatchCard
              key={batch.id}
              batch={batch}
              onOpen={() => onOpen(batch)}
              onEdit={() => setForm({ batch })}
              onDelete={() => setDeleting(batch)}
            />
          ))}
        </div>
      ) : (
        <div className="empty-state">
          <div className="empty-illustration">
            <FileSpreadsheet size={32} />
            <Plus size={16} />
          </div>
          <h3>لسه مفيش دفعات</h3>
          <p>
            ابدأ بإنشاء دفعة إجازات جديدة، وبعدها استورد كشف الأفراد من Excel.
          </p>
          <Button onClick={() => setForm({ batch: null })}>
            <Plus size={17} /> إنشاء أول دفعة
          </Button>
        </div>
      )}

      {form && (
        <BatchFormModal
          batch={form.batch}
          onClose={() => setForm(null)}
          onSaved={onOpen}
        />
      )}
      {deleting && (
        <ConfirmDeleteModal
          title="حذف الدفعة"
          message={`سيتم حذف الدفعة «${deleting.name}». قد يؤثر ذلك على سجلات الحضور المرتبطة بها.`}
          onConfirm={confirmDelete}
          onClose={() => setDeleting(null)}
        />
      )}
    </>
  );
}

import { ArrowRight, CalendarDays, Pencil, Trash2 } from "lucide-react";
import { formatDate } from "@/lib/utils";
import type { VacationBatch } from "@/types";

interface Props {
  batch: VacationBatch;
  onOpen: () => void;
  onEdit: () => void;
  onDelete: () => void;
}

export function BatchCard({ batch, onOpen, onEdit, onDelete }: Props) {
  return (
    <article className="batch-card">
      <div className="batch-card-top">
        <div className="batch-icon">
          <CalendarDays size={21} />
        </div>
        <button className="icon-btn" title="خيارات الدفعة" onClick={onEdit}>
          <Pencil size={16} />
        </button>
      </div>
      <div className="batch-title-wrap">
        <h3>{batch.name}</h3>
        <span className="batch-status">
          <span /> دفعة مسجلة
        </span>
      </div>
      <div className="batch-dates">
        <div>
          <span>من</span>
          <strong>{formatDate(batch.startDate)}</strong>
        </div>
        <div className="date-divider" />
        <div>
          <span>إلى</span>
          <strong>{formatDate(batch.endDate)}</strong>
        </div>
      </div>
      <button className="open-batch" onClick={onOpen}>
        فتح الكشف <ArrowRight size={17} />
      </button>
      <button className="batch-delete-link" onClick={onDelete}>
        <Trash2 size={14} /> حذف الدفعة
      </button>
    </article>
  );
}

import { ArrowDownToLine, FileSpreadsheet, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/Button";

interface Props {
  count: number;
  busy: boolean;
  onClear: () => void;
  onDelete: () => void;
  onExport: () => void;
  onCertificates: () => void;
}

export function BulkBar({ count, busy, onClear, onDelete, onExport, onCertificates }: Props) {
  return (
    <div className="bulk-bar">
      <div>
        <span className="selection-count">{count}</span>
        <strong>فرد محدد</strong>
        <button className="clear-selection" onClick={onClear}>
          إلغاء التحديد
        </button>
      </div>
      <div className="bulk-actions">
        <Button variant="danger" disabled={busy} onClick={onDelete}>
          <Trash2 size={15} /> حذف المحدد
        </Button>
        <Button variant="secondary" onClick={onExport}>
          <ArrowDownToLine size={15} /> تصدير المحدد
        </Button>
        <Button variant="secondary" onClick={onCertificates}>
          <FileSpreadsheet size={15} /> شهادات الحضور والغياب
        </Button>
      </div>
    </div>
  );
}

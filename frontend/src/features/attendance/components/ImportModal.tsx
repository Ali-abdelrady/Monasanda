import { FileSpreadsheet, X } from "lucide-react";
import { Modal } from "@/components/ui/Modal";
import { Button } from "@/components/ui/Button";
import type { Attendance, VacationBatch } from "@/types";
import { useImport } from "../hooks/useImport";
import { ImportDropZone } from "./ImportDropZone";
import { ImportPreview } from "./ImportPreview";

interface Props {
  batch: VacationBatch;
  existing: Attendance[];
  onClose: () => void;
}

export function ImportModal({ batch, existing, onClose }: Props) {
  const imp = useImport(batch, existing, onClose);

  return (
    <Modal
      title="استيراد كشف الإجازات"
      onClose={() => {
        if (!imp.busy) onClose();
      }}
      wide
    >
      <div className="import-body">
        {!imp.fileName ? (
          <ImportDropZone onFile={imp.readFile} />
        ) : (
          <>
            <div className="import-file-row">
              <div className="file-icon">
                <FileSpreadsheet size={21} />
              </div>
              <div>
                <strong>{imp.fileName}</strong>
                <span>{imp.rows.length} صف بيانات</span>
              </div>
              <button className="icon-btn" disabled={imp.busy} onClick={imp.reset}>
                <X size={17} />
              </button>
            </div>

            <div className="import-summary">
              <div>
                <span>إجمالي الصفوف</span>
                <strong>{imp.rows.length}</strong>
              </div>
              <div>
                <span>جاهز للاستيراد</span>
                <strong className="text-green">{imp.validRows.length}</strong>
              </div>
              <div>
                <span>يحتاج مراجعة</span>
                <strong className={imp.invalidCount ? "text-red" : ""}>
                  {imp.invalidCount}
                </strong>
              </div>
            </div>

            <ImportPreview rows={imp.rows} />

            {imp.busy && (
              <div className="import-progress">
                <div>
                  <span>جاري الاستيراد...</span>
                  <strong>
                    {imp.progress.done} / {imp.progress.total}
                  </strong>
                </div>
                <div className="progress-track">
                  <span
                    style={{
                      width: `${imp.progress.total ? (imp.progress.done / imp.progress.total) * 100 : 0}%`,
                    }}
                  />
                </div>
              </div>
            )}

            <div className="modal-actions import-actions">
              <Button variant="secondary" disabled={imp.busy} onClick={imp.reset}>
                اختيار ملف آخر
              </Button>
              <Button disabled={imp.busy || !imp.validRows.length} onClick={imp.run}>
                {imp.busy ? "جاري الاستيراد..." : `استيراد ${imp.validRows.length} فرد`}
              </Button>
            </div>
          </>
        )}
      </div>
    </Modal>
  );
}

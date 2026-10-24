import { useState } from "react";
import { ArrowDownToLine, ArrowRight, CalendarDays, Pencil, Upload } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { ConfirmDeleteModal } from "@/components/ui/ConfirmDeleteModal";
import { exportAttendance } from "@/lib/excel";
import { formatDate } from "@/lib/utils";
import type { Attendance, VacationBatch } from "@/types";
import { BatchFormModal } from "../batches/BatchFormModal";
import { useBatchAttendance } from "./hooks/useBatchAttendance";
import { useAttendanceFilters } from "./hooks/useAttendanceFilters";
import { useSelection } from "./hooks/useSelection";
import { useAttendanceActions } from "./hooks/useAttendanceActions";
import { computeStats } from "./stats";
import type { ViewMode } from "./types";
import { StatsGrid } from "./components/StatsGrid";
import { RecordsPanel } from "./components/RecordsPanel";
import { PersonFormModal } from "./components/PersonFormModal";
import { ImportModal } from "./components/ImportModal";

interface Props {
  batch: VacationBatch;
  onBack: () => void;
  onBatchChanged: (batch: VacationBatch) => void;
}

/**
 * Orchestrator: owns which modal is open and wires hooks to presentational
 * components. Rendered with key={batch.id} in App, so all state here
 * (search, filter, selection…) resets automatically when the batch changes.
 */
export function BatchDetailPage({ batch, onBack, onBatchChanged }: Props) {
  const { query, people } = useBatchAttendance(batch.id);
  const filters = useAttendanceFilters(people);
  const selection = useSelection(filters.filtered);
  const actions = useAttendanceActions();
  const stats = computeStats(people);

  const [view, setView] = useState<ViewMode>("table");
  const [editBatch, setEditBatch] = useState(false);
  const [importOpen, setImportOpen] = useState(false);
  const [personModal, setPersonModal] = useState<{ person: Attendance | null } | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Attendance | "bulk" | null>(null);

  const confirmDelete = async () => {
    if (deleteTarget === "bulk") {
      await actions.removeMany(selection.selected);
      selection.clear();
    } else if (deleteTarget) {
      if (await actions.remove(deleteTarget)) selection.drop(deleteTarget.id);
    }
    setDeleteTarget(null);
  };

  const syncedAt = query.dataUpdatedAt ? new Date(query.dataUpdatedAt) : new Date();

  return (
    <>
      <div className="detail-heading">
        <div>
          <button className="back-link" onClick={onBack}>
            <ArrowRight size={16} /> كل الدفعات
          </button>
          <div className="detail-title-row">
            <h1>{batch.name}</h1>
            <span className="batch-status">
              <span /> نشطة
            </span>
          </div>
          <p>
            <CalendarDays size={15} /> {formatDate(batch.startDate)}{" "}
            <span className="date-dash">—</span> {formatDate(batch.endDate)}
          </p>
        </div>
        <div className="heading-actions">
          <Button variant="secondary" onClick={() => setEditBatch(true)}>
            <Pencil size={16} /> تعديل الدفعة
          </Button>
          <Button variant="secondary" onClick={() => exportAttendance(people, batch.name)}>
            <ArrowDownToLine size={16} /> تصدير Excel
          </Button>
          <Button onClick={() => setImportOpen(true)}>
            <Upload size={17} /> استيراد Excel
          </Button>
        </div>
      </div>

      <StatsGrid stats={stats} />

      <RecordsPanel
        batch={batch}
        people={people}
        stats={stats}
        query={query}
        filters={filters}
        selection={selection}
        view={view}
        onView={setView}
        bulkBusy={actions.bulkBusy}
        onAdd={() => setPersonModal({ person: null })}
        onImport={() => setImportOpen(true)}
        onEdit={(person) => setPersonModal({ person })}
        onDelete={setDeleteTarget}
        onBulkDelete={() => setDeleteTarget("bulk")}
        onToggleAttendance={actions.toggle}
      />

      <div className="sync-caption">
        <span className="live-dot" /> تتم مراجعة التحديثات تلقائيًا كل 4 ثوانٍ{" "}
        <span>
          آخر مزامنة:{" "}
          {syncedAt.toLocaleTimeString("ar-EG", {
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
          })}
        </span>
      </div>

      {editBatch && (
        <BatchFormModal batch={batch} onClose={() => setEditBatch(false)} onSaved={onBatchChanged} />
      )}
      {personModal && (
        <PersonFormModal
          batch={batch}
          person={personModal.person}
          onClose={() => setPersonModal(null)}
        />
      )}
      {importOpen && (
        <ImportModal batch={batch} existing={people} onClose={() => setImportOpen(false)} />
      )}
      {deleteTarget && (
        <ConfirmDeleteModal
          title={deleteTarget === "bulk" ? "حذف الأفراد المحددين" : "حذف سجل الفرد"}
          message={
            deleteTarget === "bulk"
              ? `سيتم حذف ${selection.selected.length} سجل. لا يمكن التراجع عن هذا الإجراء.`
              : `سيتم حذف سجل ${deleteTarget.name ?? "الفرد"} نهائيًا.`
          }
          busy={actions.bulkBusy}
          onConfirm={confirmDelete}
          onClose={() => setDeleteTarget(null)}
        />
      )}
    </>
  );
}

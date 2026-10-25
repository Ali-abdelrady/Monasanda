import { CircleAlert, FileSpreadsheet, Plus, RefreshCw, Search, Upload, Users } from "lucide-react";
import { useEffect, useState } from "react";
import type { UseQueryResult } from "@tanstack/react-query";
import { Button } from "@/components/ui/Button";
import { errorMessage } from "@/lib/errors";
import { exportAttendance, exportTemplate } from "@/lib/excel";
import { useToast } from "@/context/ToastContext";
import type { Attendance, VacationBatch } from "@/types";
import type { AttendanceStats } from "../stats";
import type { ViewMode } from "../types";
import type { useAttendanceFilters } from "../hooks/useAttendanceFilters";
import type { useSelection } from "../hooks/useSelection";
import { AttendanceToolbar } from "./AttendanceToolbar";
import { BulkBar } from "./BulkBar";
import { AttendanceTable } from "./AttendanceTable";
import { AttendanceCards } from "./AttendanceCards";

const PAGE_SIZE = 10;

interface Props {
  batch: VacationBatch;
  people: Attendance[];
  stats: AttendanceStats;
  query: UseQueryResult<Attendance[]>;
  filters: ReturnType<typeof useAttendanceFilters>;
  selection: ReturnType<typeof useSelection>;
  view: ViewMode;
  onView: (v: ViewMode) => void;
  bulkBusy: boolean;
  onAdd: () => void;
  onImport: () => void;
  onEdit: (p: Attendance) => void;
  onDelete: (p: Attendance) => void;
  onBulkDelete: () => void;
  onToggleAttendance: (p: Attendance) => void;
}

export function RecordsPanel(props: Props) {
  const { batch, people, stats, query, filters, selection, view } = props;
  const { notify } = useToast();
  const [page, setPage] = useState(1);
  const pageCount = Math.ceil(filters.filtered.length / PAGE_SIZE);
  const currentPage = Math.min(page, pageCount);
  const firstIndex = (currentPage - 1) * PAGE_SIZE;
  const pageRows = filters.filtered.slice(firstIndex, firstIndex + PAGE_SIZE);

  useEffect(() => {
    setPage(1);
  }, [filters.search, filters.status]);

  useEffect(() => {
    if (page > pageCount) setPage(Math.max(1, pageCount));
  }, [page, pageCount]);

  const renderBody = () => {
    if (query.isPending)
      return (
        <div className="table-loading">
          <div className="spinner" /> جاري تحميل كشف الأفراد...
        </div>
      );

    if (query.isError)
      return (
        <div className="empty-state compact">
          <CircleAlert size={25} />
          <h3>تعذر تحميل الكشف</h3>
          <p>{errorMessage(query.error, "حاول مرة أخرى.")}</p>
          <Button variant="secondary" onClick={() => query.refetch()}>
            <RefreshCw size={15} /> إعادة المحاولة
          </Button>
        </div>
      );

    if (!people.length)
      return (
        <div className="empty-state compact">
          <Users size={28} />
          <h3>الكشف فاضي لسه</h3>
          <p>استورد ملف Excel أو أضف أول فرد يدويًا لبدء تسجيل الحضور.</p>
          <div className="empty-actions">
            <Button onClick={props.onImport}>
              <Upload size={16} /> استيراد Excel
            </Button>
            <Button variant="secondary" onClick={props.onAdd}>
              <Plus size={16} /> إضافة فرد
            </Button>
            <button className="template-link" onClick={exportTemplate}>
              تحميل قالب Excel
            </button>
          </div>
        </div>
      );

    if (!filters.filtered.length)
      return (
        <div className="empty-state compact">
          <Search size={26} />
          <h3>مفيش نتائج مطابقة</h3>
          <p>جرّب تغيّر كلمة البحث أو فلتر الحالة.</p>
          <Button variant="secondary" onClick={filters.reset}>
            مسح الفلاتر
          </Button>
        </div>
      );

    const rowProps = {
      rows: pageRows,
      selected: selection.selected,
      onToggleSelected: selection.toggle,
      onToggleAttendance: props.onToggleAttendance,
      onEdit: props.onEdit,
      onDelete: props.onDelete,
    };

    return view === "table" ? (
      <AttendanceTable
        {...rowProps}
        firstIndex={firstIndex}
        stats={stats}
        allSelected={selection.allVisibleSelected}
        onToggleAll={selection.toggleAllVisible}
      />
    ) : (
      <AttendanceCards {...rowProps} />
    );
  };

  return (
    <section className="records-panel">
      <div className="records-header">
        <div>
          <h2>كشف الأفراد</h2>
          <p>ابحث عن الفرد وسجّل التمام في ثوانٍ.</p>
        </div>
        <Button onClick={props.onAdd}>
          <Plus size={17} /> إضافة فرد
        </Button>
      </div>

      <AttendanceToolbar
        search={filters.search}
        onSearch={filters.setSearch}
        status={filters.status}
        onStatus={filters.setStatus}
        view={view}
        onView={props.onView}
      />

      {selection.selected.length > 0 && (
        <BulkBar
          count={selection.selected.length}
          busy={props.bulkBusy}
          onClear={selection.clear}
          onDelete={props.onBulkDelete}
          onExport={() =>
            exportAttendance(
              people.filter((p) => selection.selected.includes(p.id)),
              batch.name,
              true,
            )
          }
          onCertificates={() =>
            notify("ميزة شهادات الحضور والغياب ستُضاف في مرحلة لاحقة", "error")
          }
        />
      )}

      {renderBody()}
      {!query.isPending &&
        !query.isError &&
        filters.filtered.length > 0 && (
          <div className="attendance-pagination">
            <span>
              عرض {(firstIndex + 1).toLocaleString("ar-EG")} -{" "}
              {Math.min(firstIndex + PAGE_SIZE, filters.filtered.length).toLocaleString("ar-EG")} من{" "}
              {filters.filtered.length.toLocaleString("ar-EG")} فرد
            </span>
            <div className="pagination-controls" dir="rtl">
              <Button
                variant="secondary"
                disabled={currentPage <= 1}
                onClick={() => setPage((current) => Math.max(1, current - 1))}
              >
                السابق
              </Button>
              <span>
                صفحة {currentPage.toLocaleString("ar-EG")} من{" "}
                {pageCount.toLocaleString("ar-EG")}
              </span>
              <Button
                variant="secondary"
                disabled={currentPage >= pageCount}
                onClick={() =>
                  setPage((current) => Math.min(pageCount, current + 1))
                }
              >
                التالي
              </Button>
            </div>
          </div>
        )}
    </section>
  );
}

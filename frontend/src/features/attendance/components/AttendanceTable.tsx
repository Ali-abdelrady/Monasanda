import { Check, CheckCheck, Clock3, Pencil, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { formatDate, formatTime } from "@/lib/utils";
import type { Attendance } from "@/types";
import type { AttendanceStats } from "../stats";
import { StatusBadge } from "./StatusBadge";
import { getReturnStatus } from "../../../lib/utils";

interface Props {
  rows: Attendance[];
  firstIndex: number;
  stats: AttendanceStats;
  selected: number[];
  allSelected: boolean;
  onToggleAll: () => void;
  onToggleSelected: (id: number) => void;
  onToggleAttendance: (p: Attendance) => void;
  onEdit: (p: Attendance) => void;
  onDelete: (p: Attendance) => void;
}

export function AttendanceTable(props: Props) {
  const { rows, stats, selected, allSelected, firstIndex } = props;
  return (
    <div className="table-wrap">
      <table className="attendance-table">
        <thead>
          <tr>
            <th className="check-cell">
              <input
                type="checkbox"
                aria-label="تحديد كل النتائج"
                checked={allSelected}
                onChange={props.onToggleAll}
              />
            </th>
            <th>م</th>
            <th>الاسم</th>
            <th>الرقم العسكري</th>
            <th>رقم التليفون</th>
            <th>السكن</th>
            <th>متوقع العودة</th>
            <th>الحالة</th>
            <th>وقت الحضور</th>
            <th>الإجراءات</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((p, i) => (
            <tr
              key={p.id}
              className={selected.includes(p.id) ? "row-selected" : ""}
            >
              <td className="check-cell">
                <input
                  type="checkbox"
                  aria-label={`تحديد ${p.name ?? "الفرد"}`}
                  checked={selected.includes(p.id)}
                  onChange={() => props.onToggleSelected(p.id)}
                />
              </td>
              <td className="row-number">
                {(firstIndex + i + 1).toLocaleString("ar-EG")}
              </td>
              <td>
                <div className="person-cell">
                  <div>
                    <strong>{p.name || "بدون اسم"}</strong>
                    {p.notes && <span className="person-note">{p.notes}</span>}
                  </div>
                </div>
              </td>
              <td className="mono-cell">{p.militaryNumber || "—"}</td>
              <td dir="ltr" className="phone-cell">
                {p.phoneNumber || "—"}
              </td>
              <td>{p.location || "—"}</td>
              <td>
                <div>
                  {(() => {
                    const status = getReturnStatus(p.endDate);
                    return (
                      <span className={`status-badge ${status.color}`}>
                        {status.label}
                      </span>
                    );
                  })()}
                </div>
              </td>
              <td>
                <StatusBadge status={p.status} />
              </td>
              <td>
                <div className="check-time">
                  {p.checkedInAt ? (
                    <>
                      <Clock3 size={13} />
                      {formatTime(p.checkedInAt)}
                    </>
                  ) : (
                    <span className="no-time">—</span>
                  )}
                </div>
              </td>
              <td>
                <div className="row-actions">
                  <Button
                    variant={p.status === "present" ? "soft" : "primary"}
                    className="check-btn"
                    onClick={() => props.onToggleAttendance(p)}
                  >
                    {p.status === "present" ? (
                      <>
                        <Check size={15} /> إلغاء
                      </>
                    ) : (
                      <>
                        <CheckCheck size={15} /> تمام
                      </>
                    )}
                  </Button>
                  <button
                    className="icon-btn row-edit"
                    title="تعديل البيانات"
                    onClick={() => props.onEdit(p)}
                  >
                    <Pencil size={15} />
                  </button>
                  <button
                    className="icon-btn row-delete"
                    title="حذف"
                    onClick={() => props.onDelete(p)}
                  >
                    <Trash2 size={15} />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className="table-footer">
        <span>
          <span className="mini-dot green-dot" />{" "}
          {stats.present.toLocaleString("ar-EG")} حاضر{" "}
          <span className="footer-divider">·</span>{" "}
          <span className="mini-dot amber-dot" />{" "}
          {stats.absent.toLocaleString("ar-EG")} غائب
        </span>
      </div>
    </div>
  );
}

import {
  Check,
  CheckCheck,
  CalendarDays,
  Clock3,
  MoreHorizontal,
  Smartphone,
  Trash2,
  Users,
  Notebook,
} from "lucide-react";
import { Button } from "@/components/ui/Button";
import { formatDate, formatTime } from "@/lib/utils";
import type { Attendance } from "@/types";
import { StatusBadge } from "./StatusBadge";
import { getReturnStatus } from "../../../lib/utils";

interface Props {
  rows: Attendance[];
  selected: number[];
  onToggleSelected: (id: number) => void;
  onToggleAttendance: (p: Attendance) => void;
  onEdit: (p: Attendance) => void;
  onDelete: (p: Attendance) => void;
}

export function AttendanceCards({ rows, selected, ...actions }: Props) {
  return (
    <div className="people-grid">
      {rows.map((p: Attendance) => (
        <article
          key={p.id}
          className={`person-card ${p.status === "present" ? "person-card-present" : ""}`}
        >
          <div className="person-card-top">
            <input
              type="checkbox"
              checked={selected.includes(p.id)}
              onChange={() => actions.onToggleSelected(p.id)}
              aria-label={`تحديد ${p.name ?? ""}`}
            />
            <StatusBadge status={p.status} />
            <button className="icon-btn" onClick={() => actions.onEdit(p)}>
              <MoreHorizontal size={18} />
            </button>
          </div>
          <div className="card-person-identity">
            <div
              className={`person-avatar big ${p.status === "present" ? "avatar-present" : ""}`}
            >
              {(p.name ?? "?").trim().slice(0, 1)}
            </div>
            <div>
              <h3>{p.name || "بدون اسم"}</h3>
              <span>م. {p.militaryNumber || "غير مسجل"}</span>
            </div>
          </div>
          <div className="card-info-list">
            <div>
              <Smartphone size={15} />
              <span dir="ltr">{p.phoneNumber || "لا يوجد رقم تليفون"}</span>
            </div>
            <div>
              <Users size={15} />
              <span>{p.location || "السكن غير محدد"}</span>
            </div>
            <div>
              <CalendarDays size={15} />
              {(() => {
                const status = getReturnStatus(p.endDate);
                return <span>{status.label}</span>;
              })()}
            </div>
            {p.notes && (
              <div>
                <Notebook size={15} />
                <span>{p.notes}</span>
              </div>
            )}
            {p.checkedInAt && (
              <div>
                <Clock3 size={15} />
                <span>وقت الحضور {formatTime(p.checkedInAt)}</span>
              </div>
            )}
            {p.checkedInBy && (
              <div>
                <Users size={15} />
                <span>
                  تم التمام بسواطة {p?.checkedInBy?.email.split("@")[0]}
                </span>
              </div>
            )}
          </div>
          <div className="person-card-actions">
            <Button
              variant={p.status === "present" ? "soft" : "primary"}
              onClick={() => actions.onToggleAttendance(p)}
            >
              {p.status === "present" ? (
                <>
                  <Check size={16} /> إلغاء التمام
                </>
              ) : (
                <>
                  <CheckCheck size={16} /> تسجيل التمام
                </>
              )}
            </Button>
            <button
              className="icon-btn row-delete"
              onClick={() => actions.onDelete(p)}
            >
              <Trash2 size={16} />
            </button>
          </div>
        </article>
      ))}
    </div>
  );
}

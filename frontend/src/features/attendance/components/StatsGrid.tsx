import type { ReactNode } from "react";
import { Activity, BadgeCheck, Clock3, Users } from "lucide-react";
import type { AttendanceStats } from "../stats";

const ar = (n: number) => n.toLocaleString("ar-EG");

function StatCard({
  label,
  tone,
  icon,
  value,
  footer,
}: {
  label: string;
  tone: "blue" | "green" | "amber" | "violet";
  icon: ReactNode;
  value: ReactNode;
  footer: ReactNode;
}) {
  return (
    <div className="stat-card">
      <div className="stat-top">
        <span>{label}</span>
        <div className={`stat-icon ${tone}`}>{icon}</div>
      </div>
      <div className="stat-value">{value}</div>
      {footer}
    </div>
  );
}

export function StatsGrid({ stats }: { stats: AttendanceStats }) {
  return (
    <div className="stats-grid">
      <StatCard
        label="إجمالي الأفراد"
        tone="blue"
        icon={<Users size={19} />}
        value={ar(stats.total)}
        footer={<div className="stat-foot">إجمالي المسجلين في الكشف</div>}
      />
      <StatCard
        label="حضروا"
        tone="green"
        icon={<BadgeCheck size={19} />}
        value={ar(stats.present)}
        footer={
          <div className="stat-foot">
            <span className="mini-dot green-dot" /> تم تسجيل حضورهم
          </div>
        }
      />
      <StatCard
        label="لم يحضروا"
        tone="amber"
        icon={<Clock3 size={19} />}
        value={ar(stats.absent)}
        footer={
          <div className="stat-foot">
            <span className="mini-dot amber-dot" /> في انتظار التمام
          </div>
        }
      />
      <StatCard
        label="نسبة الحضور"
        tone="violet"
        icon={<Activity size={19} />}
        value={
          <>
            {ar(stats.rate)}
            <small>%</small>
          </>
        }
        footer={
          <div className="progress-track">
            <span style={{ width: `${stats.rate}%` }} />
          </div>
        }
      />
    </div>
  );
}

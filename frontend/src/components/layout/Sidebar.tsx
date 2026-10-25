import { Activity, LayoutGrid, LogOut, ShieldCheck } from "lucide-react";
import type { User } from "@/types";

interface Props {
  user: User;
  isHome: boolean;
  onHome: () => void;
  onLogout: () => void;
}

export function Sidebar({ user, isHome, onHome, onLogout }: Props) {
  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <div className="brand-mark">
          <ShieldCheck size={23} />
        </div>
        <div>
          <strong>تمامات</strong>
          <span>إدارة عودة الإجازات</span>
        </div>
      </div>
      <div className="side-label">القائمة الرئيسية</div>
      <button className={`nav-item ${isHome ? "active" : ""}`} onClick={onHome}>
        <LayoutGrid size={18} /> دفعات الإجازات
      </button>
      {/* <div className="side-note">
        <div className="side-note-icon">
          <Activity size={17} />
        </div>
        <strong>كل شيء في مكان واحد</strong>
        <p>تابع الحضور، وراجع الأعداد، وصدّر الكشوف بسهولة.</p>
      </div> */}
      <div className="sidebar-bottom">
        <div className="user-avatar">
          {user.email.slice(0, 1).toUpperCase()}
        </div>
        <div className="user-meta">
          <strong>{user.email}</strong>
          <span>مستخدم مسجل</span>
        </div>
        <button
          className="icon-btn logout-icon"
          title="تسجيل الخروج"
          onClick={onLogout}
        >
          <LogOut size={17} />
        </button>
      </div>
    </aside>
  );
}

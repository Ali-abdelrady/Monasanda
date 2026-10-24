import { ShieldCheck } from "lucide-react";
import type { User } from "@/types";

export function Topbar({ user }: { user: User }) {
  return (
    <header className="topbar">
      <div className="topbar-left">
        <div className="mobile-brand">
          <div className="brand-mark">
            <ShieldCheck size={21} />
          </div>
          <strong>تمامات</strong>
        </div>
        <span className="topbar-caption">نظام إدارة الحضور والعودة من الإجازات</span>
      </div>
      <div className="topbar-right">
        <span className="live-dot" />
        <span className="topbar-status">النظام جاهز</span>
        <div className="topbar-avatar">{user.email.slice(0, 1).toUpperCase()}</div>
      </div>
    </header>
  );
}

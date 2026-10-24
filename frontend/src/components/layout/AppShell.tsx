import type { ReactNode } from "react";
import type { User } from "@/types";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";

interface Props {
  user: User;
  isHome: boolean;
  onHome: () => void;
  onLogout: () => void;
  children: ReactNode;
}

export function AppShell({ user, isHome, onHome, onLogout, children }: Props) {
  return (
    <div className="app-shell">
      <Sidebar user={user} isHome={isHome} onHome={onHome} onLogout={onLogout} />
      <main className="main-area">
        <Topbar user={user} />
        <div className="page-content">{children}</div>
        <footer className="app-footer">
          <span>
            تمامات <span className="footer-divider">·</span> نظام إدارة العودة من الإجازات
          </span>
          <span>واجهة إدارية داخلية</span>
        </footer>
      </main>
    </div>
  );
}

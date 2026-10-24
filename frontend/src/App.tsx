import { useState } from "react";
import { useAuth } from "./hooks/useAuth";
import { useToast } from "./context/ToastContext";
import { VacationBatch } from "./types";
import { FullLoader } from "./components/ui/FullLoader";
import { LoginPage } from "./features/auth/LoginPage";
import { AppShell } from "./components/layout/AppShell";
import { BatchDetailPage } from "./features/attendance/BatchDetailPage";
import { BatchesPage } from "./features/batches/BatchesPage";


export default function App() {
  const auth = useAuth();
  const { notify } = useToast();
  const [activeBatch, setActiveBatch] = useState<VacationBatch | null>(null);

  if (auth.isChecking) return <FullLoader />;
  if (!auth.user) return <LoginPage onLogin={auth.login} />;

  return (
    <AppShell
      user={auth.user}
      isHome={!activeBatch}
      onHome={() => setActiveBatch(null)}
      onLogout={() => {
        setActiveBatch(null);
        auth.logout();
        notify("تم تسجيل الخروج من الواجهة");
      }}
    >
      {activeBatch ? (
        <BatchDetailPage
          key={activeBatch.id}
          batch={activeBatch}
          onBack={() => setActiveBatch(null)}
          onBatchChanged={setActiveBatch}
        />
      ) : (
        <BatchesPage onOpen={setActiveBatch} />
      )}
    </AppShell>
  );
}

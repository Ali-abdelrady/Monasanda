import { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from "react";
import { Toast, type ToastState } from "@/components/ui/Toast";

type Notify = (message: string, type?: ToastState["type"]) => void;

const ToastContext = createContext<{ notify: Notify } | null>(null);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toast, setToast] = useState<ToastState | null>(null);
  const notify = useCallback<Notify>(
    (message, type = "success") => setToast({ message, type }),
    [],
  );
  const close = useCallback(() => setToast(null), []);
  const value = useMemo(() => ({ notify }), [notify]);

  return (
    <ToastContext.Provider value={value}>
      {children}
      <Toast toast={toast} onClose={close} />
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used inside <ToastProvider>");
  return ctx;
}

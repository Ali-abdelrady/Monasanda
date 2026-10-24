import { useEffect } from "react";
import { Check, CircleAlert, X } from "lucide-react";

export interface ToastState {
  type: "success" | "error";
  message: string;
}

export function Toast({
  toast,
  onClose,
}: {
  toast: ToastState | null;
  onClose: () => void;
}) {
  useEffect(() => {
    if (!toast) return;
    const timer = window.setTimeout(onClose, 4200);
    return () => window.clearTimeout(timer);
  }, [toast, onClose]);

  if (!toast) return null;
  return (
    <div className={`toast toast-${toast.type}`}>
      <span>
        {toast.type === "success" ? <Check size={18} /> : <CircleAlert size={18} />}
      </span>
      {toast.message}
      <button onClick={onClose}>
        <X size={16} />
      </button>
    </div>
  );
}

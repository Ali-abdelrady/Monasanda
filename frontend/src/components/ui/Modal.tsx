import { useEffect, type ReactNode } from "react";
import { X } from "lucide-react";

interface Props {
  title: string;
  children: ReactNode;
  onClose: () => void;
  wide?: boolean;
}

export function Modal({ title, children, onClose, wide = false }: Props) {
  useEffect(() => {
    const listener = (e: KeyboardEvent) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", listener);
    return () => window.removeEventListener("keydown", listener);
  }, [onClose]);

  return (
    <div
      className="modal-backdrop"
      onMouseDown={(e) => e.target === e.currentTarget && onClose()}
    >
      <section
        className={`modal ${wide ? "modal-wide" : ""}`}
        role="dialog"
        aria-modal="true"
      >
        <div className="modal-head">
          <h2>{title}</h2>
          <button className="icon-btn" onClick={onClose} aria-label="إغلاق">
            <X size={19} />
          </button>
        </div>
        {children}
      </section>
    </div>
  );
}

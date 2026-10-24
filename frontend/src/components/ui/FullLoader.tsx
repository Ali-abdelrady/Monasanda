import { ShieldCheck } from "lucide-react";

export function FullLoader() {
  return (
    <div className="full-loader">
      <div className="brand-mark">
        <ShieldCheck size={27} />
      </div>
      <div className="spinner" />
      <p>جاري التحقق من جلسة الدخول...</p>
    </div>
  );
}

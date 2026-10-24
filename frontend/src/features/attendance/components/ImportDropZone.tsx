import { useRef } from "react";
import { FileSpreadsheet, Upload } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { exportTemplate } from "@/lib/excel";
import { useToast } from "@/context/ToastContext";

export function ImportDropZone({ onFile }: { onFile: (file: File) => void }) {
  const { notify } = useToast();
  const fileRef = useRef<HTMLInputElement>(null);

  return (
    <div
      className="drop-zone"
      onDragOver={(e) => e.preventDefault()}
      onDrop={(e) => {
        e.preventDefault();
        const file = e.dataTransfer.files[0];
        if (file && /\.xlsx?$/.test(file.name)) onFile(file);
        else notify("اختر ملف Excel بامتداد xlsx أو xls", "error");
      }}
    >
      <div className="upload-icon">
        <FileSpreadsheet size={27} />
      </div>
      <h3>ارفع كشف الإجازات من Excel</h3>
      <p>اسحب الملف هنا أو اختره من جهازك. ندعم ملفات XLSX وXLS.</p>
      <Button onClick={() => fileRef.current?.click()}>
        <Upload size={16} /> اختيار ملف Excel
      </Button>
      <input
        ref={fileRef}
        type="file"
        accept=".xlsx,.xls"
        hidden
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) onFile(file);
          e.currentTarget.value = "";
        }}
      />
      <button className="template-link" onClick={exportTemplate}>
        تحميل قالب Excel المعتمد
      </button>
      <div className="import-columns">
        الأعمدة: م، الاسم، السكن، الرقم العسكري، رقم التليفون، من، إلى، ملاحظات
      </div>
    </div>
  );
}

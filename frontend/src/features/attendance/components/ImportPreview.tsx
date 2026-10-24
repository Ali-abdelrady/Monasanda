import { Check } from "lucide-react";
import { formatDate } from "@/lib/utils";
import type { ParsedRow } from "@/lib/excel";

const PREVIEW_LIMIT = 100;

export function ImportPreview({ rows }: { rows: ParsedRow[] }) {
  return (
    <div className="import-preview">
      <table>
        <thead>
          <tr>
            <th>صف</th>
            <th>الاسم</th>
            <th>الرقم العسكري</th>
            <th>رقم التليفون</th>
            <th>من</th>
            <th>إلى</th>
            <th>التحقق</th>
          </tr>
        </thead>
        <tbody>
          {rows.slice(0, PREVIEW_LIMIT).map((row) => (
            <tr key={row.rowNumber}>
              <td>{row.rowNumber}</td>
              <td>{row.input.name || "—"}</td>
              <td>{row.input.militaryNumber || "—"}</td>
              <td dir="ltr">{row.input.phoneNumber || "—"}</td>
              <td>{formatDate(row.input.startDate)}</td>
              <td>{formatDate(row.input.endDate)}</td>
              <td>
                {row.errors.length ? (
                  <span className="error-chip" title={row.errors.join("، ")}>
                    {row.errors.join("، ")}
                  </span>
                ) : (
                  <span className="valid-chip">
                    <Check size={12} /> سليم
                  </span>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {rows.length > PREVIEW_LIMIT && (
        <div className="preview-limit">
          يتم عرض أول 100 صف في المعاينة، وسيتم فحص جميع الصفوف عند الاستيراد.
        </div>
      )}
    </div>
  );
}

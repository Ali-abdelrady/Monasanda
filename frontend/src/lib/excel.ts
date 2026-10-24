import * as XLSX from "xlsx";
import type { Attendance, AttendanceInput } from "../types";
import { dateInput } from "./utils";

const HEADERS = [
  "م",
  "الاسم",
  "السكن",
  "الرقم العسكري",
  "رقم التليفون",
  "من",
  "إلى",
  "ملاحظات",
];
const clean = (value: unknown) => (value == null ? "" : String(value).trim());
function excelDate(value: unknown): string | null {
  if (value == null || value === "") return null;
  if (value instanceof Date && !Number.isNaN(value.getTime()))
    return value.toISOString();
  if (typeof value === "number") {
    const parsed = XLSX.SSF.parse_date_code(value);
    if (parsed)
      return new Date(Date.UTC(parsed.y, parsed.m - 1, parsed.d)).toISOString();
  }
  const text = String(value).trim();
  if (!text) return null;
  const date = new Date(text);
  if (!Number.isNaN(date.getTime())) return date.toISOString();
  const match = text.match(/^(\d{1,2})[/.\-](\d{1,2})[/.\-](\d{4})$/);
  if (match)
    return new Date(
      Date.UTC(Number(match[3]), Number(match[2]) - 1, Number(match[1])),
    ).toISOString();
  return null;
}
export type ParsedRow = {
  rowNumber: number;
  input: AttendanceInput;
  errors: string[];
};
export function parseAttendanceFile(
  file: File,
  batchId: number,
): Promise<ParsedRow[]> {
  return file.arrayBuffer().then((buffer) => {
    const workbook = XLSX.read(buffer, { type: "array", cellDates: true });
    const sheet = workbook.Sheets[workbook.SheetNames[0]];
    if (!sheet) throw new Error("ملف Excel لا يحتوي على ورقة بيانات");
    const rows = XLSX.utils.sheet_to_json<unknown[]>(sheet, {
      header: 1,
      defval: "",
      raw: true,
    });
    if (!rows.length) return [];
    const first = (rows[0] as unknown[]).map(clean);
    const hasHeader = first.some((v) =>
      ["الاسم", "السكن", "الرقم العسكري", "رقم التليفون", "ملاحظات"].includes(
        v,
      ),
    );
    const start = hasHeader ? 1 : 0;
    return rows
      .slice(start)
      .map((raw, index) => {
        const row = raw as unknown[];
        const get = (i: number) => row[i] ?? "";
        const name = clean(get(1));
        const militaryNumber = clean(get(3)) || null;
        const phoneNumber = clean(get(4)) || null;
        const startDate = excelDate(get(5));
        const endDate = excelDate(get(6));
        const errors: string[] = [];
        if (!name) errors.push("الاسم مطلوب");
        if (phoneNumber && phoneNumber.length > 12)
          errors.push("رقم التليفون أطول من 12 خانة");
        if (get(5) !== "" && !startDate)
          errors.push("تاريخ بداية الإجازة غير صحيح");
        if (get(6) !== "" && !endDate)
          errors.push("تاريخ نهاية الإجازة غير صحيح");
        if (startDate && endDate && new Date(endDate) < new Date(startDate))
          errors.push("تاريخ النهاية قبل البداية");
        return {
          rowNumber: start + index + 1,
          errors,
          input: {
            batchId,
            name: name || null,
            location: clean(get(2)) || null,
            militaryNumber,
            phoneNumber,
            status: "absent",
            notes: clean(get(7)) || null,
            startDate,
            endDate,
          },
        };
      })
      .filter(
        (item) =>
          item.input.name ||
          item.input.militaryNumber ||
          item.input.phoneNumber ||
          item.input.location,
      );
  });
}
export function exportAttendance(
  records: Attendance[],
  batchName: string,
  selectedOnly = false,
) {
  const rows = records.map((item, index) => ({
    م: index + 1,
    الاسم: item.name ?? "",
    السكن: item.location ?? "",
    "الرقم العسكري": item.militaryNumber ?? "",
    "رقم التليفون": item.phoneNumber ?? "",
    من: item.startDate ? dateInput(item.startDate) : "",
    إلى: item.endDate ? dateInput(item.endDate) : "",
    ملاحظات: item.notes ?? "",
    "حالة الحضور": item.status === "present" ? "حاضر" : "غائب",
    "وقت الحضور": item.checkedInAt
      ? new Date(item.checkedInAt).toLocaleString("ar-EG", { hour12: false })
      : "",
    "تم التمام بواسطة (ID)": item.checkedInBy ?? "",
  }));
  const sheet = XLSX.utils.json_to_sheet(rows, {
    header: [...HEADERS, "حالة الحضور", "وقت الحضور", "تم التمام بواسطة (ID)"],
  });
  sheet["!cols"] = [
    { wch: 6 },
    { wch: 26 },
    { wch: 18 },
    { wch: 18 },
    { wch: 17 },
    { wch: 14 },
    { wch: 14 },
    { wch: 28 },
    { wch: 16 },
    { wch: 22 },
    { wch: 22 },
  ];
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, sheet, "كشف الحضور");
  const date = new Date().toISOString().slice(0, 10);
  XLSX.writeFile(
    workbook,
    `${batchName}-${selectedOnly ? "المحدد-" : ""}${date}.xlsx`,
  );
}
export function exportTemplate() {
  const sheet = XLSX.utils.aoa_to_sheet([HEADERS]);
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, sheet, "كشف الإجازات");
  XLSX.writeFile(workbook, "قالب-كشف-الإجازات.xlsx");
}

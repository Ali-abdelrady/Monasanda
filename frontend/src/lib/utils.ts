export function cn(...values: Array<string | false | null | undefined>) {
  return values.filter(Boolean).join(" ");
}
export function formatDate(value?: string | null) {
  if (!value) return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return new Intl.DateTimeFormat("ar-EG", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(date);
}
export function formatTime(value?: string | null) {
  if (!value) return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return new Intl.DateTimeFormat("ar-EG", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  }).format(date);
}
export function toApiDate(value: string, end = false) {
  if (!value) return "";
  return new Date(`${value}T${end ? "23:59:59" : "00:00:00"}`).toISOString();
}
export function dateInput(value?: string | null) {
  return value ? new Date(value).toISOString().slice(0, 10) : "";
}

export function getReturnStatus(endDate: string) {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const end = new Date(endDate);
  end.setHours(0, 0, 0, 0);

  const days = Math.round((end.getTime() - today.getTime()) / 86400000);

  if (days > 1)
    return {
      label: `بعد ${days} يوم`,
      color: "bg-blue-100 text-blue-700",
    };

  if (days === 1)
    return {
      label: "بكرة",
      color: "bg-yellow-100 text-yellow-700",
    };

  if (days === 0)
    return {
      label: "العودة اليوم",
      color: "bg-green-100 text-green-700",
    };

  const overdue = Math.abs(days);

  return {
    label: `متأخر ${overdue} يوم`,
    color: "bg-red-100 text-red-700",
  };
}

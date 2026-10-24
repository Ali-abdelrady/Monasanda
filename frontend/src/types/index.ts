export type User = { id: number; email: string };
export type VacationBatch = {
  id: number;
  name: string;
  startDate: string;
  endDate: string;
  createdBy?: number;
  createdAt?: string;
  updatedAt?: string;
};
export type AttendanceStatus = "present" | "absent";
export type Attendance = {
  id: number;
  batchId: number;
  name: string | null;
  location: string | null;
  phoneNumber: string | null;
  militaryNumber: string | null;
  status: AttendanceStatus;
  notes: string | null;
  startDate: string | null;
  endDate: string | null;
  checkedInAt: string | null;
  checkedInBy: User | null;
  createdBy?: number;
  createdAt?: string;
  updatedAt?: string;
};
export type AttendanceInput = Omit<
  Attendance,
  "id" | "checkedInAt" | "checkedInBy" | "createdBy" | "createdAt" | "updatedAt"
>;
export type BatchInput = Pick<VacationBatch, "name" | "startDate" | "endDate">;

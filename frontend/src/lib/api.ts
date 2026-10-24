import type {
  Attendance,
  AttendanceInput,
  BatchInput,
  User,
  VacationBatch,
} from "../types";

const API_BASE = (import.meta.env.VITE_API_BASE_URL ?? "").replace(/\/$/, "");

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers: { "Content-Type": "application/json", ...(init.headers ?? {}) },
  });
  let body: any = null;
  const text = await response.text();
  if (text) {
    try {
      body = JSON.parse(text);
    } catch {
      body = text;
    }
  }
  if (!response.ok || body?.hasError === true) {
    const message =
      typeof body === "object" && body
        ? body.ErrorMessage || body.message || body.error
        : null;
    throw new Error(message || `حصل خطأ في الطلب (${response.status})`);
  }
  return body as T;
}

export const api = {
  async login(email: string, password: string) {
    return request<{ hasError: false }>("/api/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },

  async me() {
    const result = await request<{ hasError: false; user: User }>("/api/me");
    return result.user;
  },

  async batches() {
    const result = await request<{
      hasError: false;
      vacationBatches: VacationBatch[];
    }>("/api/vacation-batches");
    return result.vacationBatches ?? [];
  },

  async batch(id: number) {
    const result = await request<{
      hasError: false;
      vacationBatch: VacationBatch;
    }>(`/api/vacation-batches/${id}`);
    return result.vacationBatch;
  },

  async createBatch(input: BatchInput) {
    const result = await request<{
      hasError: false;
      vacationBatch: VacationBatch;
    }>("/api/vacation-batches", {
      method: "POST",
      body: JSON.stringify(input),
    });
    return result.vacationBatch;
  },

  async updateBatch(id: number, input: BatchInput) {
    const result = await request<{
      hasError: false;
      vacationBatch: VacationBatch;
    }>(`/api/vacation-batches/${id}`, {
      method: "PUT",
      body: JSON.stringify(input),
    });
    return result.vacationBatch;
  },

  deleteBatch: (id: number) =>
    request<{ hasError: false }>(`/api/vacation-batches/${id}`, {
      method: "DELETE",
    }),

  async attendances() {
    const result = await request<{
      hasError: false;
      attendances: Attendance[];
    }>("/api/attendances");
    return result.attendances ?? [];
  },

  async createAttendance(input: AttendanceInput) {
    const result = await request<{ hasError: false; attendance: Attendance }>(
      "/api/attendances",
      { method: "POST", body: JSON.stringify(input) },
    );
    return result.attendance;
  },

  async updateAttendance(id: number, input: AttendanceInput) {
    const result = await request<{ hasError: false; attendance: Attendance }>(
      `/api/attendances/${id}`,
      { method: "PUT", body: JSON.stringify(input) },
    );
    return result.attendance;
  },
  
  deleteAttendance: (id: number) =>
    request<{ hasError: false }>(`/api/attendances/${id}`, {
      method: "DELETE",
    }),
};

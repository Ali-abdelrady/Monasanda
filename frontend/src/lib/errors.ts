export const errorMessage = (e: unknown, fallback: string) =>
  e instanceof Error ? e.message : fallback;

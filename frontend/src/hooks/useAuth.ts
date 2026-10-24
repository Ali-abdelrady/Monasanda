import { useEffect, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { User } from "@/types";

export function useAuth() {
  const qc = useQueryClient();
  const [user, setUser] = useState<User | null>(null);
  const [authChecked, setAuthChecked] = useState(false);

  const meQuery = useQuery({
    queryKey: ["me"],
    queryFn: api.me,
    enabled: !authChecked,
    retry: false,
  });

  useEffect(() => {
    if (meQuery.data) {
      setUser(meQuery.data);
      setAuthChecked(true);
    } else if (meQuery.isError) {
      setAuthChecked(true);
    }
  }, [meQuery.data, meQuery.isError]);

  const login = async (email: string, password: string) => {
    await api.login(email.trim(), password);
    setUser(await api.me());
    setAuthChecked(true);
  };

  const logout = () => {
    setUser(null);
    qc.clear();
  };

  return { user, isChecking: !authChecked && meQuery.isPending, login, logout };
}

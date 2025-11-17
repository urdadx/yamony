import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";

export interface User {
  id: number;
  username: string;
  email: string;
  emailVerified: boolean;
  image: string;
}

interface MeResponse {
  user: User;
}

export function useMe() {
  return useQuery<User>({
    queryKey: ["me"],
    queryFn: async () => {
      const response = await api.get<MeResponse>("/me");
      return response.data.user;
    },
    staleTime: 5 * 60 * 1000, 
    retry: false,
  });
}

import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";

export interface Vault {
  id: number;
  user_id: number;
  name: string;
  description?: string;
  icon?: string;
  theme?: string;
  is_favorite: boolean;
  item_count: number;
  created_at: string;
  updated_at: string;
}

export function useVaults() {
  return useQuery<Vault[]>({
    queryKey: ["vaults"],
    queryFn: async () => {
      const response = await api.get<Vault[]>("/vaults");
      return response.data;
    },
    staleTime: 5 * 60 * 1000,
  });
}

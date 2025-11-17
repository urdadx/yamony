import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

export interface Session {
  id: number;
  user_id: number;
  session_token: string;
  expires_at: {
    Time: string;
    Valid: boolean;
  };
  created_at: {
    Time: string;
    Valid: boolean;
  };
  updated_at: {
    Time: string;
    Valid: boolean;
  };
  active_page_id: {
    Int32: number;
    Valid: boolean;
  };
}

interface SessionsResponse {
  sessions: Session[];
}

const isSessionExpired = (session: Session): boolean => {
  if (!session.expires_at || !session.expires_at.Valid) return true;
  const expiresAt = new Date(session.expires_at.Time).getTime();
  return expiresAt <= Date.now();
};

export function useSession() {
  const query = useQuery<Session[]>({
    queryKey: ["sessions"],
    queryFn: async () => {
      const response = await api.get<SessionsResponse>("/sessions");
      return response.data.sessions || [];
    },
    staleTime: 7 * 60 * 1000, 
    retry: false, 
  });

  const currentSession = useMemo(() => {
    if (!query.data || query.data.length === 0) return null;
    
    return query.data.reduce((latest, session) => {
      if (!latest) return session;
      
      const latestTime = latest.updated_at && latest.updated_at.Valid 
        ? new Date(latest.updated_at.Time).getTime() 
        : 0;
      const sessionTime = session.updated_at && session.updated_at.Valid 
        ? new Date(session.updated_at.Time).getTime() 
        : 0;
      
      return sessionTime > latestTime ? session : latest;
    }, null as Session | null);
  }, [query.data]);

  const activePageId = useMemo(() => {
    if (!currentSession?.active_page_id || !currentSession.active_page_id.Valid) return null;
    return currentSession.active_page_id.Int32;
  }, [currentSession]);

  return {
    ...query,
    sessions: query.data || [],
    currentSession,
    activePageId,
    isSessionExpired,
  };
}

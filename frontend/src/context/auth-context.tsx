import { useSession } from "@/hooks/use-session";
import { api } from "@/lib/api";
import type React from "react";
import { createContext, useContext, useEffect, useState } from "react";

export interface AuthState {
  isAuthenticated: boolean;
}

export interface User {
  id: number;
  username: string;
  email: string;
  emailVerified: boolean;
  image: string;
}

export interface AuthResponse {
  message: string;
  user: User;
}

export interface LogoutResponse {
  message: string;
}

interface AuthContextType {
  authState: AuthState;
  setAuthState: React.Dispatch<React.SetStateAction<AuthState>>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { currentSession, isLoading } = useSession();

  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
  });

  useEffect(() => {
    if (!isLoading) {
      setAuthState({
        isAuthenticated: !!currentSession,
      });
    }
  }, [currentSession, isLoading]);

  const logout = async () => {
    try {
      await api.post<LogoutResponse>("/logout");
      setAuthState({ isAuthenticated: false });
    }
    catch (error) {
      console.error("Logout error:", error);
      setAuthState({ isAuthenticated: false });
    }
  };

  return (
    <AuthContext.Provider
      value={{
        authState,
        setAuthState,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
const LAST_LOGIN_KEY = 'yamony_last_login_method';

export type LoginMethod = 'google' | 'email';

export interface LastLoginInfo {
  method: LoginMethod;
  timestamp: number;
  email?: string;
}

export function getLastLoginMethod(): LastLoginInfo | null {
  if (typeof window === 'undefined') return null;
  
  try {
    const stored = localStorage.getItem(LAST_LOGIN_KEY);
    if (!stored) return null;
    
    const parsed = JSON.parse(stored) as LastLoginInfo;
    
    // Check if the stored data is older than 30 days
    const thirtyDaysAgo = Date.now() - (30 * 24 * 60 * 60 * 1000);
    if (parsed.timestamp < thirtyDaysAgo) {
      localStorage.removeItem(LAST_LOGIN_KEY);
      return null;
    }
    
    return parsed;
  } catch {
    localStorage.removeItem(LAST_LOGIN_KEY);
    return null;
  }
}

export function setLastLoginMethod(method: LoginMethod, email?: string): void {
  if (typeof window === 'undefined') return;
  
  const info: LastLoginInfo = {
    method,
    timestamp: Date.now(),
    email,
  };
  
  localStorage.setItem(LAST_LOGIN_KEY, JSON.stringify(info));
}

export function clearLastLoginMethod(): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(LAST_LOGIN_KEY);
}
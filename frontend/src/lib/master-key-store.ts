import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface MasterKeyStore {
  salt: string | null;
  hasSetup: boolean;
  setSalt: (salt: string) => void;
  setHasSetup: (value: boolean) => void;
  clear: () => void;
}

export const useMasterKeyStore = create<MasterKeyStore>()(
  persist(
    (set) => ({
      salt: null,
      hasSetup: false,
      setSalt: (salt: string) => set({ salt, hasSetup: true }),
      setHasSetup: (value: boolean) => set({ hasSetup: value }),
      clear: () => set({ salt: null, hasSetup: false }),
    }),
    {
      name: 'master-key-storage',
    }
  )
);

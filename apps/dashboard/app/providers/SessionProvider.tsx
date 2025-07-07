import React, { createContext, useCallback, useEffect, useState } from "react";
import { getAuthClient } from "~/lib/api-clients";

interface SessionProviderProps {
  children?: React.ReactNode;
  storageKey?: string;
}

export interface SessionProviderState {
  authenticated: boolean;
  setAuthenticated: (authenticated: boolean) => void;
}

const initialState = {
  authenticated: false,
  setAuthenticated: (authenticated: boolean) => {
    throw new Error("Called setAuthenticated outside of SessionProvider");
  },
};

const SessionProviderContext = createContext<SessionProviderState>(initialState);

let externalSetAuthenticated: ((authenticated: boolean) => void) | null = null;

export function setSessionAuthenticated(authenticated: boolean) {
  if (externalSetAuthenticated) {
    externalSetAuthenticated(authenticated);
  } else {
    throw new Error("SessionProvider is not mounted");
  }
}

export default function SessionProvider({ children, storageKey = "authenticated" }: SessionProviderProps) {
  let isAuthenticated = false;

  if (typeof window !== "undefined") {
    isAuthenticated = sessionStorage.getItem(storageKey) === "true";
  }
  const [authenticated, setAuthenticated] = useState<boolean>(isAuthenticated);

  useEffect(() => {
    externalSetAuthenticated = setAuthenticated;
    return () => {
      externalSetAuthenticated = null;
    };
  }, [setAuthenticated]);

  useEffect(() => {
    if (authenticated) {
      // TODO:
      // try if we jwt and session id cookies are set by sending a request to the backend
      // if not, follow through the authentication flow
    }

    // TODO:
    // call the backend to check if refresh token is set and valid, if so, set jwt and session id cookies and update the authenticated state
    // else, set authenticated to false and redirect to the login page.
  }, []);

  return (
    <SessionProviderContext.Provider
      value={{
        authenticated: authenticated,
        setAuthenticated: (auth: boolean) => {
          setAuthenticated(auth);
          if (typeof window !== "undefined") {
            sessionStorage.setItem(storageKey, auth.toString());
          }
        },
      }}
    >
      {children}
    </SessionProviderContext.Provider>
  );
}

export function useSession() {
  const context = React.useContext(SessionProviderContext);
  if (!context) {
    throw new Error("useSession must be used within a SessionProvider");
  }
  return context;
}

export function useLogout() {
  const { setAuthenticated } = useSession();

  return useCallback(async () => {
    setAuthenticated(false);
    const res = await (await getAuthClient()).logout({});

    if (!res.success) {
      console.error("Logout failed (?????)");
    }

    return;
  }, [setAuthenticated]);
}

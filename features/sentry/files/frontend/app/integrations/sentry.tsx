import * as Sentry from "@sentry/react";
import type { ReactNode } from "react";

const dsn = import.meta.env.VITE_SENTRY_DSN;

if (dsn) {
  Sentry.init({
    dsn,
    integrations: [
      Sentry.browserTracingIntegration(),
      Sentry.replayIntegration(),
    ],
    tracesSampleRate: 1.0,
    replaysSessionSampleRate: 0.1,
    replaysOnErrorSampleRate: 1.0,
  });
}

export function SentryProvider({ children }: { children: ReactNode }) {
  if (!dsn) return <>{children}</>;
  return (
    <Sentry.ErrorBoundary fallback={<p>An error has occurred</p>}>
      {children}
    </Sentry.ErrorBoundary>
  );
}

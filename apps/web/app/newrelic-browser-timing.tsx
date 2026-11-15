import newrelic from 'newrelic';
import Script from 'next/script';

interface NewRelicAgent {
  collector: { isConnected: () => boolean };
  on: (event: 'connected', listener: () => void) => void;
}

export async function NewRelicBrowserTiming() {
  const { agent } = newrelic as typeof newrelic & { agent?: NewRelicAgent };
  if (!agent) {
    return null;
  }

  if (agent.collector.isConnected() === false) {
    await new Promise<void>((resolve) => {
      agent.on('connected', resolve);
    });
  }

  const browserTimingHeader = newrelic.getBrowserTimingHeader({
    hasToRemoveScriptWrapper: true,
    allowTransactionlessInjection: true,
  });

  return (
    <Script
      id='nr-browser-agent'
      dangerouslySetInnerHTML={{ __html: browserTimingHeader }}
    />
  );
}

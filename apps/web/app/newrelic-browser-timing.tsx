import newrelic from 'newrelic';
import Script from 'next/script';

export async function NewRelicBrowserTiming() {
  // @ts-ignore
  const agent = newrelic.agent;
  if (!agent) {
    return null;
  }

  if (agent.collector.isConnected() === false) {
    await new Promise((resolve) => {
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
      strategy='beforeInteractive'
      dangerouslySetInnerHTML={{ __html: browserTimingHeader }}
    />
  );
}

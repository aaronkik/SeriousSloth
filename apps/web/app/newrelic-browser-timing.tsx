import newrelic from 'newrelic';
import Script from 'next/script';

export async function NewRelicBrowserTiming() {
  // @ts-ignore
  if (newrelic.agent.collector.isConnected() === false) {
    await new Promise((resolve) => {
      // @ts-ignore
      newrelic.agent.on('connected', resolve);
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

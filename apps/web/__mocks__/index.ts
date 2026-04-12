const initMocks = async () => {
  if (typeof window === 'undefined') {
    const { server } = await import('./msw-server');
    server.listen();
  } else {
    const { worker } = await import('./msw-browser');
    worker.start();
  }
};

initMocks();

export {};

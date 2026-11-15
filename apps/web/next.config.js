// @ts-check

/** @type {import('next').NextConfig} */
const nextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  poweredByHeader: false,
  serverExternalPackages: ['newrelic'],
  // The agent requires newrelic.js via a runtime path, which file tracing cannot see.
  outputFileTracingIncludes: {
    '/*': ['./newrelic.js'],
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'static-cdn.jtvnw.net',
      },
    ],
  },
  reactStrictMode: true,
  transpilePackages: ['radix-ui', 'sonner', 'lucide-react'],
};

module.exports = nextConfig;

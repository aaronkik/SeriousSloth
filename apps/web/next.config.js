// @ts-check

/** @type {import('next').NextConfig} */
const nextConfig = {
  cacheComponents: true,
  poweredByHeader: false,
  serverExternalPackages: ['newrelic'],
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

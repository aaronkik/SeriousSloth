/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['static-cdn.jtvnw.net'],
  },
  reactStrictMode: true,
  swcMinify: true,
};

module.exports = nextConfig;

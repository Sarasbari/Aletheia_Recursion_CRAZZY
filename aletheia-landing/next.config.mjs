/** @type {import('next').NextConfig} */
const nextConfig = {
  allowedDevOrigins: ["http://192.168.0.113:3000"],
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://192.168.0.117:8080/api/:path*",
      },
    ];
  },
};

export default nextConfig;

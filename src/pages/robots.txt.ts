import type { APIRoute } from 'astro';

export const GET: APIRoute = ({ site }) => {
  if (!site) {
    throw new Error('robots.txt를 생성하려면 astro.config.mjs의 site 설정이 필요합니다.');
  }

  const sitemapUrl = new URL('/sitemap.xml', site);
  const robots = `User-agent: *\nAllow: /\n\nSitemap: ${sitemapUrl.href}\n`;

  return new Response(robots, { headers: { 'Content-Type': 'text/plain; charset=utf-8' } });
};

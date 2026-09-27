import { defineCollection, reference } from 'astro:content';
import { file, glob } from 'astro/loaders';
import { z } from 'astro/zod';

const slug = z.string().regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/);
const localizedName = z.object({
  ko: z.string().min(1),
  en: z.string().min(1),
  ja: z.string().min(1),
});

const topics = defineCollection({
  loader: file('./src/data/topics.json'),
  schema: z.object({
    slug,
    name: z.string().min(1),
  }),
});

const series = defineCollection({
  loader: file('./src/data/series.json'),
  schema: z.object({
    slug,
    name: localizedName,
  }),
});

const posts = defineCollection({
  loader: glob({ pattern: '**/*.mdx', base: './src/content/posts' }),
  schema: z
    .object({
      title: z.string().min(1),
      description: z.string().min(1),
      author: z.string().min(1).optional(),
      thumbnail: z.url().optional(),
      publishedAt: z.coerce.date(),
      updatedAt: z.coerce.date().optional(),
      topics: z.array(reference('topics')).min(1),
      series: reference('series').optional(),
      seriesOrder: z.number().int().positive().optional(),
      draft: z.boolean().default(false),
    })
    .superRefine((post, context) => {
      if (Boolean(post.series) !== Boolean(post.seriesOrder)) {
        context.addIssue({
          code: 'custom',
          message: 'series와 seriesOrder는 함께 작성해야 합니다.',
        });
      }
    }),
});

export const collections = { posts, topics, series };

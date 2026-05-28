import { z } from 'zod';

const contentFormatSchema = z.enum([
	'markdown',
	'html',
	'raw_html',
	'json',
	'links',
	'metadata',
]);

export const partialConfigSchema = z
	.object({
		request: z
			.object({
				headers: z.record(z.string(), z.string()).optional(),
				timeout: z.string().optional(),
				retry_count: z.number().int().min(0).max(10).optional(),
				retry_interval: z.string().optional(),
			})
			.optional(),
		content: z
			.object({
				formats: z.array(contentFormatSchema).min(1).optional(),
				only_main_content: z.boolean().optional(),
				selector: z.string().optional(),
				extract_links: z.boolean().optional(),
				extract_metadata: z.boolean().optional(),
			})
			.optional(),
		crawl: z
			.object({
				enabled: z.boolean().optional(),
				max_depth: z.number().int().min(0).max(10).optional(),
				max_pages: z.number().int().min(1).max(100000).optional(),
				respect_robots_txt: z.boolean().optional(),
			})
			.optional(),
	})
	.passthrough();

export function parsePartialConfig(json: unknown) {
	return partialConfigSchema.safeParse(json);
}

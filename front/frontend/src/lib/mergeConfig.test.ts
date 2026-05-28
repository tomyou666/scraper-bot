import { describe, expect, it } from 'vitest';
import { DEFAULT_APP_CONFIG } from './defaults';
import { configForMode2, mergeConfig } from './mergeConfig';

describe('mergeConfig', () => {
	it('node overrides workspace', () => {
		const merged = mergeConfig({}, { crawl: { max_depth: 5 } }, undefined, {
			crawl: { max_depth: 1 },
		});
		expect(merged.crawl?.max_depth).toBe(1);
	});

	it('mode2 uses app defaults only', () => {
		const m2 = configForMode2({ crawl: { max_depth: 9 } });
		expect(m2.crawl?.max_depth).toBe(9);
		expect(m2.crawl?.max_pages).toBe(DEFAULT_APP_CONFIG.crawl?.max_pages);
	});
});

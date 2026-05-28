import { describe, expect, it } from 'vitest';
import { normalizeUrl } from './normalizeUrl';

describe('normalizeUrl', () => {
	it('lowercases host and removes fragment', () => {
		expect(normalizeUrl('HTTPS://Example.COM/path#frag')).toBe(
			'https://example.com/path',
		);
	});

	it('sorts query keys', () => {
		expect(normalizeUrl('https://ex.com/?b=2&a=1')).toBe(
			'https://ex.com/?a=1&b=2',
		);
	});

	it('removes default https port', () => {
		expect(normalizeUrl('https://ex.com:443/')).toBe('https://ex.com/');
	});
});

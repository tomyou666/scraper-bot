const DEFAULT_PORTS: Record<string, string> = {
	'http:': '80',
	'https:': '443',
};

/**
 * ノード一意判定用 URL 正規化。
 */
export function normalizeUrl(raw: string): string {
	const trimmed = raw.trim();
	if (!trimmed) {
		throw new Error('URL is empty');
	}
	const withScheme = /^[a-zA-Z][a-zA-Z\d+\-.]*:/.test(trimmed)
		? trimmed
		: `https://${trimmed}`;
	const parsed = new URL(withScheme);
	parsed.protocol = parsed.protocol.toLowerCase();
	parsed.hostname = parsed.hostname.toLowerCase();
	parsed.hash = '';

	const defaultPort = DEFAULT_PORTS[parsed.protocol];
	if (defaultPort && parsed.port === defaultPort) {
		parsed.port = '';
	}

	if (
		(parsed.protocol === 'http:' || parsed.protocol === 'https:') &&
		parsed.pathname === ''
	) {
		parsed.pathname = '/';
	} else if (parsed.pathname.length > 1 && parsed.pathname.endsWith('/')) {
		parsed.pathname = parsed.pathname.replace(/\/+$/, '') || '/';
	}

	if (parsed.search) {
		const params = new URLSearchParams(parsed.search);
		const sorted = new URLSearchParams();
		[...params.keys()].sort().forEach((key) => {
			const values = params.getAll(key).sort();
			values.forEach((v) => {
				sorted.append(key, v);
			});
		});
		parsed.search = sorted.toString() ? `?${sorted.toString()}` : '';
	}

	return parsed.toString();
}

export function hostFromUrl(url: string): string {
	return new URL(url).hostname;
}

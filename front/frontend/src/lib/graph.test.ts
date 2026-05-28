import { describe, expect, it } from 'vitest';
import type { GraphEdge, GraphNode } from '@/types/graph';
import { getDescendantNodeIds, getForwardReachableExisting } from './graph';

const nodes: GraphNode[] = [
	{
		id: 'a',
		urlNormalized: 'https://x/a',
		label: 'a',
		position: { x: 0, y: 0 },
		nodeSettings: {},
		crawlExclude: false,
		status: 'idle',
	},
	{
		id: 'b',
		urlNormalized: 'https://x/b',
		label: 'b',
		position: { x: 0, y: 0 },
		nodeSettings: {},
		crawlExclude: false,
		status: 'idle',
	},
	{
		id: 'c',
		urlNormalized: 'https://x/c',
		label: 'c',
		position: { x: 0, y: 0 },
		nodeSettings: {},
		crawlExclude: false,
		status: 'idle',
	},
];
const edges: GraphEdge[] = [
	{ id: 'e1', source: 'a', target: 'b' },
	{ id: 'e2', source: 'b', target: 'c' },
];

describe('graph', () => {
	it('getForwardReachableExisting returns BFS order without start', () => {
		expect(getForwardReachableExisting('a', nodes, edges)).toEqual(['b', 'c']);
	});

	it('getDescendantNodeIds', () => {
		const d = getDescendantNodeIds('a', edges);
		expect([...d]).toEqual(['b', 'c']);
	});
});

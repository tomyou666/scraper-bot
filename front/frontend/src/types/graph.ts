import type { PartialConfig } from './config';
import type { CrawlResultPreview } from './crawl';

export type NodeStatus = 'idle' | 'running' | 'success' | 'error' | 'skipped';

export interface GraphNode {
	id: string;
	urlNormalized: string;
	label: string;
	position: { x: number; y: number };
	nodeSettings: PartialConfig;
	crawlExclude: boolean;
	status: NodeStatus;
	lastResult?: CrawlResultPreview;
	lastError?: string;
}

export interface GraphEdge {
	id: string;
	source: string;
	target: string;
}

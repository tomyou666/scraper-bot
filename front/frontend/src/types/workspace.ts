import type { PartialConfig } from './config';
import type { GraphEdge, GraphNode } from './graph';

export interface Workspace {
	id: string;
	name: string;
	seedUrl: string;
	settings: PartialConfig;
	exclude_urls: string[];
	nodes: GraphNode[];
	edges: GraphEdge[];
	domainSettings: Record<string, PartialConfig>;
}

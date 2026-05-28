import { useMemo, useState } from 'react';
import { Alert } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { messages } from '@/i18n/messages';
import { useAppStore } from '@/stores/appStore';
import type { ContentFormat } from '@/types/config';
import { getActiveFormats } from '@/types/crawl';
import type { GraphNode } from '@/types/graph';

export function RightSidebar() {
	const ws = useAppStore((s) => s.getActiveWorkspace());
	const node = useAppStore((s) => s.getSelectedNode());
	const selectedDomain = useAppStore((s) => s.selectedDomain);
	const runHistory = useAppStore((s) => s.runHistory);
	const crawlError = useAppStore((s) => s.crawlError);
	const clearCrawlError = useAppStore((s) => s.clearCrawlError);
	const appDefaults = useAppStore((s) => s.appDefaults);
	const setNodeCrawlExclude = useAppStore((s) => s.setNodeCrawlExclude);
	const updateDomainSettings = useAppStore((s) => s.updateDomainSettings);

	const formats = useMemo(
		() =>
			getActiveFormats(
				ws?.settings.content?.formats ?? appDefaults.content?.formats,
			),
		[ws, appDefaults],
	);

	const [tab, setTab] = useState<ContentFormat>(formats[0] ?? 'markdown');

	if (selectedDomain && ws) {
		const domainCfg = ws.domainSettings[selectedDomain] ?? {};
		return (
			<aside className='flex w-72 shrink-0 flex-col border-l border-border bg-card'>
				<div className='border-b border-border px-3 py-2 text-xs font-semibold'>
					{messages.right.domainSettings}: {selectedDomain}
				</div>
				<ScrollArea className='flex-1 p-3'>
					<Label className='flex items-center gap-2'>
						<Checkbox
							checked={domainCfg.crawl?.respect_robots_txt ?? true}
							onCheckedChange={(checked) =>
								updateDomainSettings(selectedDomain, {
									crawl: { respect_robots_txt: checked },
								})
							}
						/>
						robots.txt に従う
					</Label>
					<p className='mt-4 text-xs text-muted-foreground'>
						max_depth: {domainCfg.crawl?.max_depth ?? '—'}
					</p>
					<textarea
						className='mt-2 min-h-32 w-full rounded-lg border border-input bg-background p-2 font-mono text-xs'
						placeholder='{"crawl":{"max_depth":2}}'
						defaultValue={JSON.stringify(domainCfg, null, 2)}
						onBlur={(e) => {
							try {
								const parsed = JSON.parse(e.target.value);
								updateDomainSettings(selectedDomain, parsed);
							} catch {
								/* ignore invalid json on blur */
							}
						}}
					/>
				</ScrollArea>
			</aside>
		);
	}

	if (node) {
		return (
			<aside className='flex w-72 shrink-0 flex-col border-l border-border bg-card'>
				<div className='border-b border-border px-3 py-2'>
					<p className='text-xs font-semibold'>{messages.right.nodeResult}</p>
					<p className='truncate text-xs text-muted-foreground'>
						{node.urlNormalized}
					</p>
					{node.status === 'error' && node.lastError && (
						<Alert variant='destructive' className='mt-2 text-xs'>
							{messages.error.nodeFailed}: {node.lastError}
						</Alert>
					)}
				</div>
				<div className='border-b border-border px-3 py-2'>
					<Label className='flex items-center gap-2 text-xs'>
						<Checkbox
							checked={node.crawlExclude}
							onCheckedChange={(c) => setNodeCrawlExclude(node.id, c)}
						/>
						{messages.right.crawlExclude}
					</Label>
				</div>
				<Tabs
					value={tab}
					onValueChange={(v) => setTab(v as ContentFormat)}
					className='flex flex-1 flex-col px-3'
				>
					<TabsList>
						{formats.map((f) => (
							<TabsTrigger key={f} value={f}>
								{f}
							</TabsTrigger>
						))}
					</TabsList>
					<ScrollArea className='flex-1 pb-3'>
						{formats.map((f) => (
							<TabsContent key={f} value={f}>
								<NodeFormatContent format={f} node={node} />
							</TabsContent>
						))}
					</ScrollArea>
				</Tabs>
			</aside>
		);
	}

	return (
		<aside className='flex w-72 shrink-0 flex-col border-l border-border bg-card'>
			<div className='border-b border-border px-3 py-2 text-xs font-semibold'>
				{messages.right.runSummary}
			</div>
			{crawlError && (
				<Alert variant='destructive' className='m-2 text-xs'>
					<div className='flex justify-between gap-2'>
						<span>
							{messages.error.crawlFailed}: {crawlError.message}
						</span>
						<button
							type='button'
							onClick={clearCrawlError}
							className='shrink-0'
						>
							×
						</button>
					</div>
				</Alert>
			)}
			<ScrollArea className='flex-1 p-3'>
				{runHistory.length === 0 ? (
					<p className='text-xs text-muted-foreground'>
						{messages.right.noSelection}
					</p>
				) : (
					<div className='space-y-2'>
						<p className='text-xs font-medium text-muted-foreground'>
							{messages.right.history}
						</p>
						{runHistory.map((run) => (
							<div
								key={run.id}
								className='rounded-lg border border-border p-2 text-xs'
							>
								<div className='flex items-center justify-between'>
									<Badge variant='secondary'>モード {run.mode}</Badge>
									<span className='text-muted-foreground'>
										{run.stoppedReason ?? '—'}
									</span>
								</div>
								<p className='mt-1 text-muted-foreground'>
									{new Date(run.startedAt).toLocaleString()}
								</p>
								<p className='mt-1'>
									成功 {run.succeeded} / 失敗 {run.failed} / スキップ{' '}
									{run.skipped}
								</p>
							</div>
						))}
					</div>
				)}
			</ScrollArea>
		</aside>
	);
}

function NodeFormatContent({
	format,
	node,
}: {
	format: string;
	node: GraphNode;
}) {
	const r = node.lastResult;
	if (!r) {
		return <p className='text-xs text-muted-foreground'>結果がありません</p>;
	}
	if (format === 'markdown') {
		return (
			<pre className='whitespace-pre-wrap text-xs'>{r.markdown ?? '—'}</pre>
		);
	}
	if (format === 'links') {
		return (
			<ul className='list-inside list-disc text-xs'>
				{(r.links ?? []).map((l) => (
					<li key={l} className='truncate'>
						{l}
					</li>
				))}
			</ul>
		);
	}
	if (format === 'metadata') {
		return (
			<dl className='space-y-1 text-xs'>
				{Object.entries(r.metadata ?? {}).map(([k, v]) => (
					<div key={k}>
						<dt className='text-muted-foreground'>{k}</dt>
						<dd>{v}</dd>
					</div>
				))}
			</dl>
		);
	}
	return (
		<p className='text-xs text-muted-foreground'>（モック未対応: {format}）</p>
	);
}

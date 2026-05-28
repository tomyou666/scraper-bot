import { Plus, Trash2 } from 'lucide-react';
import { useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { messages } from '@/i18n/messages';
import { hostFromUrl } from '@/lib/normalizeUrl';
import { cn } from '@/lib/utils';
import { useAppStore } from '@/stores/appStore';

export function LeftSidebar() {
	const workspaces = useAppStore((s) => s.workspaces);
	const activeWorkspaceId = useAppStore((s) => s.activeWorkspaceId);
	const setActiveWorkspace = useAppStore((s) => s.setActiveWorkspace);
	const deleteWorkspace = useAppStore((s) => s.deleteWorkspace);
	const openNewWorkspaceDialog = useAppStore((s) => s.openNewWorkspaceDialog);
	const openAddNodeDialog = useAppStore((s) => s.openAddNodeDialog);
	const openDeleteNodeDialog = useAppStore((s) => s.openDeleteNodeDialog);
	const selectedNodeId = useAppStore((s) => s.selectedNodeId);
	const selectedDomain = useAppStore((s) => s.selectedDomain);
	const selectDomain = useAppStore((s) => s.selectDomain);
	const activeWorkspace = useAppStore((s) =>
		s.workspaces.find((w) => w.id === s.activeWorkspaceId),
	);
	const crawlStatus = useAppStore((s) => s.crawlStatus);
	const domains = useMemo(() => {
		if (!activeWorkspace) return [];
		const hosts = new Set(
			activeWorkspace.nodes.map((n) => hostFromUrl(n.urlNormalized)),
		);
		return [...hosts].sort();
	}, [activeWorkspace]);

	return (
		<aside className='flex w-56 shrink-0 flex-col border-r border-border bg-sidebar'>
			<div className='flex items-center justify-between border-b border-sidebar-border px-2 py-2'>
				<span className='text-xs font-semibold'>
					{messages.sidebar.workspaces}
				</span>
				<Button variant='ghost' size='icon-xs' onClick={openNewWorkspaceDialog}>
					<Plus className='size-3.5' />
				</Button>
			</div>
			<ScrollArea className='max-h-40 flex-none px-1 py-1'>
				{workspaces.length === 0 ? (
					<p className='px-2 py-2 text-xs text-muted-foreground'>
						{messages.sidebar.emptyWorkspaces}
					</p>
				) : (
					workspaces.map((ws) => (
						<button
							key={ws.id}
							type='button'
							className={cn(
								'flex w-full items-center justify-between rounded-md px-2 py-1.5 text-left text-xs hover:bg-sidebar-accent',
								activeWorkspaceId === ws.id &&
									'bg-sidebar-accent font-medium text-sidebar-accent-foreground',
							)}
							onClick={() => setActiveWorkspace(ws.id)}
						>
							<span className='truncate'>{ws.name}</span>
							<Button
								variant='ghost'
								size='icon-xs'
								className='shrink-0'
								onClick={(e) => {
									e.stopPropagation();
									deleteWorkspace(ws.id);
								}}
							>
								<Trash2 className='size-3' />
							</Button>
						</button>
					))
				)}
			</ScrollArea>

			<div className='flex flex-1 flex-col border-t border-sidebar-border'>
				<div className='flex items-center justify-between px-2 py-2'>
					<span className='text-xs font-semibold'>
						{messages.sidebar.domains}
					</span>
					<div className='flex gap-0.5'>
						<Button
							variant='ghost'
							size='icon-xs'
							onClick={() => openAddNodeDialog()}
							title={messages.sidebar.newNode}
						>
							<Plus className='size-3.5' />
						</Button>
						<Button
							variant='ghost'
							size='icon-xs'
							disabled={!selectedNodeId || crawlStatus !== 'idle'}
							onClick={openDeleteNodeDialog}
							title={messages.sidebar.deleteNode}
						>
							<Trash2 className='size-3.5' />
						</Button>
					</div>
				</div>
				<ScrollArea className='flex-1 px-1 pb-2'>
					{domains.length === 0 ? (
						<p className='px-2 py-2 text-xs text-muted-foreground'>
							{messages.sidebar.emptyDomains}
						</p>
					) : (
						domains.map((host) => (
							<button
								key={host}
								type='button'
								className={cn(
									'block w-full truncate rounded-md px-2 py-1.5 text-left text-xs hover:bg-sidebar-accent',
									selectedDomain === host && 'bg-sidebar-accent font-medium',
								)}
								onClick={() => selectDomain(host)}
							>
								{host}
							</button>
						))
					)}
				</ScrollArea>
			</div>
		</aside>
	);
}

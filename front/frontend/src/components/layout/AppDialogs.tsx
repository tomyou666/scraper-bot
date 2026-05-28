import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { messages } from '@/i18n/messages';
import { useAppStore } from '@/stores/appStore';

export function AppDialogs() {
	const showNewWorkspaceDialog = useAppStore((s) => s.showNewWorkspaceDialog);
	const closeNewWorkspaceDialog = useAppStore((s) => s.closeNewWorkspaceDialog);
	const createWorkspace = useAppStore((s) => s.createWorkspace);
	const workspaces = useAppStore((s) => s.workspaces);

	const showAddNodeDialog = useAppStore((s) => s.showAddNodeDialog);
	const closeAddNodeDialog = useAppStore((s) => s.closeAddNodeDialog);
	const addNode = useAppStore((s) => s.addNode);

	const showDeleteNodeDialog = useAppStore((s) => s.showDeleteNodeDialog);
	const closeDeleteNodeDialog = useAppStore((s) => s.closeDeleteNodeDialog);
	const deleteSelectedSubtree = useAppStore((s) => s.deleteSelectedSubtree);
	const selectedNode = useAppStore((s) => s.getSelectedNode());
	const ws = useAppStore((s) => s.getActiveWorkspace());

	const [wsName, setWsName] = useState('My Workspace');
	const [wsUrl, setWsUrl] = useState('https://example.com/');
	const [nodeUrl, setNodeUrl] = useState('https://');

	const mustShowNewWs = showNewWorkspaceDialog || workspaces.length === 0;

	return (
		<>
			<Dialog
				open={mustShowNewWs}
				onOpenChange={(open) => {
					if (!open && workspaces.length > 0) closeNewWorkspaceDialog();
				}}
			>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>{messages.dialog.newWorkspaceTitle}</DialogTitle>
					</DialogHeader>
					<div className='space-y-3'>
						<div>
							<Label>{messages.dialog.newWorkspaceName}</Label>
							<Input
								className='mt-1'
								value={wsName}
								onChange={(e) => setWsName(e.target.value)}
							/>
						</div>
						<div>
							<Label>{messages.dialog.newWorkspaceUrl}</Label>
							<Input
								className='mt-1'
								value={wsUrl}
								onChange={(e) => setWsUrl(e.target.value)}
							/>
						</div>
					</div>
					<DialogFooter>
						{workspaces.length > 0 && (
							<Button
								variant='outline'
								size='sm'
								onClick={closeNewWorkspaceDialog}
							>
								{messages.dialog.cancel}
							</Button>
						)}
						<Button size='sm' onClick={() => createWorkspace(wsName, wsUrl)}>
							{messages.dialog.create}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>

			<Dialog open={showAddNodeDialog} onOpenChange={closeAddNodeDialog}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>{messages.dialog.addNodeTitle}</DialogTitle>
					</DialogHeader>
					<Label>{messages.dialog.addNodeUrl}</Label>
					<Input
						className='mt-1'
						value={nodeUrl}
						onChange={(e) => setNodeUrl(e.target.value)}
					/>
					<DialogFooter>
						<Button variant='outline' size='sm' onClick={closeAddNodeDialog}>
							{messages.dialog.cancel}
						</Button>
						<Button size='sm' onClick={() => addNode(nodeUrl)} disabled={!ws}>
							{messages.dialog.add}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>

			<Dialog open={showDeleteNodeDialog} onOpenChange={closeDeleteNodeDialog}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>{messages.dialog.deleteNodeTitle}</DialogTitle>
					</DialogHeader>
					<p className='text-sm'>{messages.dialog.deleteNodeConfirm}</p>
					{selectedNode && (
						<p className='mt-2 truncate text-xs text-muted-foreground'>
							{selectedNode.urlNormalized}
						</p>
					)}
					<DialogFooter>
						<Button variant='outline' size='sm' onClick={closeDeleteNodeDialog}>
							{messages.dialog.cancel}
						</Button>
						<Button
							variant='destructive'
							size='sm'
							onClick={deleteSelectedSubtree}
						>
							{messages.dialog.delete}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</>
	);
}

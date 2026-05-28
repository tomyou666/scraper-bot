import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { ScrollArea } from '@/components/ui/scroll-area';
import { messages } from '@/i18n/messages';
import { parsePartialConfig } from '@/schemas/config';
import { useAppStore } from '@/stores/appStore';

export function MenuBar() {
	const setAppDefaults = useAppStore((s) => s.setAppDefaults);
	const appDefaults = useAppStore((s) => s.appDefaults);
	const [settingsOpen, setSettingsOpen] = useState(false);
	const [jsonText, setJsonText] = useState('');
	const [jsonError, setJsonError] = useState<string | null>(null);

	const openSettings = () => {
		setJsonText(JSON.stringify(appDefaults, null, 2));
		setJsonError(null);
		setSettingsOpen(true);
	};

	const saveSettings = () => {
		try {
			const parsed = JSON.parse(jsonText);
			const result = parsePartialConfig(parsed);
			if (!result.success) {
				setJsonError(result.error.message);
				return;
			}
			setAppDefaults(result.data);
			setSettingsOpen(false);
		} catch {
			setJsonError('JSON の解析に失敗しました');
		}
	};

	return (
		<>
			<div className='flex h-8 items-center gap-1 border-b border-border bg-card px-2 text-xs'>
				<span className='px-2 font-medium'>{messages.menu.file}</span>
				<Button variant='ghost' size='xs' disabled title={messages.scrbPhase2}>
					{messages.menu.openScrb}
				</Button>
				<Button variant='ghost' size='xs' disabled title={messages.scrbPhase2}>
					{messages.menu.saveScrb}
				</Button>
				<span className='mx-1 text-muted-foreground'>|</span>
				<span className='px-2 font-medium'>{messages.menu.settings}</span>
				<Button variant='ghost' size='xs' onClick={openSettings}>
					{messages.menu.appDefaults}
				</Button>
			</div>

			<Dialog open={settingsOpen} onOpenChange={setSettingsOpen}>
				<DialogContent className='max-h-[80vh]'>
					<DialogHeader>
						<DialogTitle>{messages.menu.appDefaults}</DialogTitle>
					</DialogHeader>
					<Label>JSON</Label>
					<ScrollArea className='max-h-96'>
						<textarea
							className='mt-1 min-h-64 w-full rounded-lg border border-input bg-background p-2 font-mono text-xs'
							value={jsonText}
							onChange={(e) => setJsonText(e.target.value)}
						/>
					</ScrollArea>
					{jsonError && (
						<p className='mt-2 text-xs text-destructive'>{jsonError}</p>
					)}
					<DialogFooter>
						<Button
							variant='outline'
							size='sm'
							onClick={() => setSettingsOpen(false)}
						>
							{messages.dialog.cancel}
						</Button>
						<Button size='sm' onClick={saveSettings}>
							{messages.dialog.confirm}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</>
	);
}

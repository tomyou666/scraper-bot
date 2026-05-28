import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { messages } from '@/i18n/messages';
import { useAppStore } from '@/stores/appStore';

export function GlobalErrorBanner() {
	const globalError = useAppStore((s) => s.globalError);
	const clearGlobalError = useAppStore((s) => s.clearGlobalError);

	if (!globalError) return null;

	return (
		<Alert
			variant='destructive'
			className='mx-2 mt-2 flex items-center justify-between gap-2'
		>
			<div>
				<strong className='font-medium'>{messages.error.globalBanner}: </strong>
				{globalError.message}
			</div>
			<Button variant='ghost' size='xs' onClick={clearGlobalError}>
				×
			</Button>
		</Alert>
	);
}

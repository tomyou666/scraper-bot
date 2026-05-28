import { cn } from '@/lib/utils';

function Alert({
	className,
	variant = 'default',
	children,
}: {
	className?: string;
	variant?: 'default' | 'destructive';
	children: React.ReactNode;
}) {
	return (
		<div
			role='alert'
			className={cn(
				'rounded-lg border px-3 py-2 text-sm',
				variant === 'destructive'
					? 'border-destructive/50 bg-destructive/10 text-destructive'
					: 'border-border bg-muted/50 text-foreground',
				className,
			)}
		>
			{children}
		</div>
	);
}

export { Alert };

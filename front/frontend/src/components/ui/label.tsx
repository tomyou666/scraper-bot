import type * as React from 'react';
import { cn } from '@/lib/utils';

function Label({ className, ...props }: React.ComponentProps<'label'>) {
	return (
		// biome-ignore lint/a11y/noLabelWithoutControl: This reusable component forwards htmlFor/children from props.
		<label
			className={cn('text-xs font-medium text-muted-foreground', className)}
			{...props}
		/>
	);
}

export { Label };

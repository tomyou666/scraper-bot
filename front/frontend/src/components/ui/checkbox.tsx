import type * as React from 'react';
import { cn } from '@/lib/utils';

function Checkbox({
	className,
	checked,
	onCheckedChange,
	...props
}: Omit<React.ComponentProps<'input'>, 'type' | 'onChange'> & {
	onCheckedChange?: (checked: boolean) => void;
}) {
	return (
		<input
			type='checkbox'
			checked={checked}
			onChange={(e) => onCheckedChange?.(e.target.checked)}
			className={cn('size-4 accent-primary', className)}
			{...props}
		/>
	);
}

export { Checkbox };

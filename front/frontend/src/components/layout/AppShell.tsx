import { CrawlGraph } from '@/components/graph/CrawlGraph';
import { AppBootstrap } from './AppBootstrap';
import { AppDialogs } from './AppDialogs';
import { ControlBar } from './ControlBar';
import { GlobalErrorBanner } from './GlobalErrorBanner';
import { LeftSidebar } from './LeftSidebar';
import { MenuBar } from './MenuBar';
import { RightSidebar } from './RightSidebar';

export function AppShell() {
	return (
		<AppBootstrap>
			<div className='flex h-screen flex-col overflow-hidden'>
				<MenuBar />
				<GlobalErrorBanner />
				<ControlBar />
				<div className='flex min-h-0 flex-1'>
					<LeftSidebar />
					<main className='relative flex min-w-0 flex-1 flex-col'>
						<CrawlGraph />
					</main>
					<RightSidebar />
				</div>
				<AppDialogs />
			</div>
		</AppBootstrap>
	);
}

import { useStores } from '@/hooks/useStores';
import { observer } from 'mobx-react-lite';
import { Skeleton } from './ui/skeleton';
import SessionHistoryItem from './SessionHistoryItem';

const SessionsHistory = observer(() => {
    const { rootStore } = useStores();

    return (
        <div className='max-w-md w-full mx-auto'>
            <div className='space-y-4'>
                {rootStore.sessionsLoading
                    ? Array.from({ length: 5 }).map((_, i) => (
                          <div
                              key={i}
                              className='bg-gray-200 rounded-lg p-4 transition-colors duration-300 relative'
                          >
                              <div className='flex items-center justify-between'>
                                  <div className='space-y-2'>
                                      <Skeleton className='h-4 w-32' />
                                      <Skeleton className='h-3 w-24' />
                                  </div>
                                  <Skeleton className='h-5 w-5' />
                              </div>
                          </div>
                      ))
                    : rootStore.sessions
                          .slice()
                          .sort(
                              (a, b) =>
                                  new Date(b.created_at).getTime() -
                                  new Date(a.created_at).getTime()
                          )
                          .map((session) => (
                              <SessionHistoryItem session={session} key={session.id} />
                          ))}
            </div>
        </div>
    );
});

export default SessionsHistory;

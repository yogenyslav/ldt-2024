import { useStores } from '@/hooks/useStores';
import { TrashIcon } from 'lucide-react';
import { observer } from 'mobx-react-lite';
import { Skeleton } from './ui/skeleton';

const SessionsHistory = observer(() => {
    const { rootStore } = useStores();

    return (
        <div className=' p-6 max-w-md mx-auto'>
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
                              <div
                                  key={session.id}
                                  className='session-history-item bg-gray-100 hover:bg-gray-200 rounded-lg p-4 transition-colors duration-300 hover:cursor-pointer'
                              >
                                  <h3 className='text-sm font-medium mb-1'>
                                      {session.title || session.id.slice(0, 8) + '...'}
                                  </h3>
                                  <div className='flex items-center justify-between'>
                                      <p className='text-gray-500 text-sm'>
                                          {new Date(session.created_at).toLocaleDateString()}
                                      </p>
                                      <button className='session-history-item__delete-button text-gray-400 hover:text-gray-600 focus:outline-none group hidden'>
                                          <TrashIcon className='h-4 w-4 group-hover:block' />
                                          <span className='sr-only'>Удалить из истории</span>
                                      </button>
                                  </div>
                              </div>
                          ))}
            </div>
        </div>
    );
});

export default SessionsHistory;

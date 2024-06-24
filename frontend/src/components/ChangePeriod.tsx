import { useStores } from '@/hooks/useStores';
import { Button } from './ui/button';
import { ChatCommand } from '@/api/models';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from './ui/dialog';
import { useState } from 'react';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { cn } from '@/lib/utils';
import { CalendarIcon } from 'lucide-react';
import { Calendar } from './ui/calendar';
import { format } from 'date-fns';
import { Label } from './ui/label';

type Props = {
    setShowInvalidButton: (show: boolean) => void;
};

const ChangePeriod = ({ setShowInvalidButton }: Props) => {
    const { rootStore } = useStores();

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [startDate, setStartDate] = useState<Date>();
    const [endDate, setEndDate] = useState<Date>();

    const handlePeriodChange = () => {
        rootStore.sendMessage({
            command: ChatCommand.Invalid,
        });

        const startDatePrompt = startDate ? format(startDate, 'dd-MM-yyyy') : '';
        const endDatePrompt = endDate ? format(endDate, 'dd-MM-yyyy') : '';

        rootStore.sendMessage({
            prompt: `Период с ${startDatePrompt} по ${endDatePrompt}`,
        });

        setShowInvalidButton(false);
    };

    return (
        <>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogContent className='sm:max-w-[425px]'>
                    <DialogHeader>
                        <DialogTitle>Изменить период прогноза</DialogTitle>
                        <DialogDescription>Задайте новый период для прогноза</DialogDescription>
                    </DialogHeader>

                    <Label>Начальная дата</Label>
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button
                                variant={'outline'}
                                className={cn(
                                    'w-[240px] justify-start text-left font-normal',
                                    !startDate && 'text-muted-foreground'
                                )}
                            >
                                <CalendarIcon className='mr-2 h-4 w-4' />
                                {startDate ? (
                                    format(startDate, 'PPP')
                                ) : (
                                    <span>Выбор начальной даты</span>
                                )}
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent className='w-auto p-0' align='start'>
                            <Calendar
                                mode='single'
                                selected={startDate}
                                onSelect={setStartDate}
                                initialFocus
                            />
                        </PopoverContent>
                    </Popover>

                    <Label>Конечная дата</Label>
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button
                                variant={'outline'}
                                className={cn(
                                    'w-[240px] justify-start text-left font-normal',
                                    !endDate && 'text-muted-foreground'
                                )}
                            >
                                <CalendarIcon className='mr-2 h-4 w-4' />
                                {endDate ? (
                                    format(endDate, 'PPP')
                                ) : (
                                    <span>Выбор конечной даты</span>
                                )}
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent className='w-auto p-0' align='start'>
                            <Calendar
                                mode='single'
                                selected={endDate}
                                onSelect={setEndDate}
                                initialFocus
                            />
                        </PopoverContent>
                    </Popover>

                    <DialogFooter>
                        <Button disabled={!(startDate && endDate)} onClick={handlePeriodChange}>
                            Изменить
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Button
                onClick={() => {
                    setIsDialogOpen(true);
                }}
                variant='outline'
                className='flex-1'
            >
                Уточнить даты
            </Button>
        </>
    );
};

export default ChangePeriod;

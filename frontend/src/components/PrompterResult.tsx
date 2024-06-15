import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { CalendarDays, FileBox, PackageSearch } from 'lucide-react';
import { mapIncomingMessageTypeToText } from '@/utils/query-type-mapper';
import { IncomingMessageType } from '@/api/models';

type Props = {
    product?: string;
    period?: string;
    type?: IncomingMessageType;
};

const PrompterResult = ({ period, product, type }: Props) => {
    return (
        <>
            <div className='grid gap-4 md:grid-cols-2 lg:grid-cols-3'>
                {type && (
                    <Card>
                        <CardHeader className='flex flex-row items-center justify-between pb-2'>
                            <CardTitle className='text-sm font-medium'>Тип запроса</CardTitle>
                            <PackageSearch className='w-4 h-4 text-gray-500 dark:text-gray-400' />
                        </CardHeader>
                        <CardContent>
                            <div className='text-lg font-bold'>
                                {mapIncomingMessageTypeToText(type)}
                            </div>
                        </CardContent>
                    </Card>
                )}

                {product && (
                    <Card>
                        <CardHeader className='flex flex-row items-center justify-between pb-2'>
                            <CardTitle className='text-sm font-medium'>Продукт</CardTitle>
                            <FileBox className='w-4 h-4 text-gray-500 dark:text-gray-400' />
                        </CardHeader>
                        <CardContent>
                            <div className='text-lg font-bold'>{product}</div>
                        </CardContent>
                    </Card>
                )}

                {period && (
                    <Card>
                        <CardHeader className='flex flex-row items-center justify-between pb-2'>
                            <CardTitle className='text-sm font-medium'>Период</CardTitle>
                            <CalendarDays className='w-4 h-4 text-gray-500 dark:text-gray-400' />
                        </CardHeader>
                        <CardContent>
                            <div className='text-2xl font-bold'>{period}</div>
                        </CardContent>
                    </Card>
                )}
            </div>
        </>
    );
};

export default PrompterResult;

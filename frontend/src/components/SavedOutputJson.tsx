import { DeliveryRow, OutputJson } from '@/api/models';
import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Label } from './ui/label';
import { Input } from './ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/table';
import FavoritesApiService from '@/api/FavoritesApiService';
import { LoaderButton } from './ui/loader-button';
import { toast } from './ui/use-toast';

type Props = {
    outputJson: OutputJson;
    id: number;
};

const SavedOutputJson = ({ outputJson, id }: Props) => {
    const [deliverySchedules, setDeliverySchedules] = useState<DeliveryRow[]>(outputJson.rows);
    const [editMode, setEditMode] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    const handleSave = () => {
        setIsSaving(true);

        FavoritesApiService.updateFavorite({
            id,
            response: { ...outputJson, rows: deliverySchedules },
        })
            .then(() => {
                setEditMode(false);

                toast({
                    title: 'Успех',
                    description: 'Изменения сохранены',
                });
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось сохранить изменения',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsSaving(false);
            });
    };

    const handleDelete = () => {
        setIsDeleting(true);

        FavoritesApiService.deleteFavorite(id)
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Прогноз успешно удален',
                });
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось удалить прогноз',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsDeleting(false);
            });
    };

    const downloadFile = () => {
        if (outputJson) {
            const data = JSON.stringify(outputJson);

            const blob = new Blob([data], { type: 'application/json' });
            const url = URL.createObjectURL(blob);

            const a = document.createElement('a');
            a.href = url;
            a.download = 'result.json';
            document.body.appendChild(a);
            a.click();
            a.remove();
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Сохраненные планы закупок</CardTitle>
                <div className='ml-auto gap-2 flex'>
                    <LoaderButton variant='outline' isLoading={isDeleting} onClick={handleDelete}>
                        Удалить
                    </LoaderButton>

                    {editMode ? (
                        <LoaderButton isLoading={isSaving} onClick={handleSave}>
                            Сохранить
                        </LoaderButton>
                    ) : (
                        <Button onClick={() => setEditMode(true)}>Редактировать</Button>
                    )}
                </div>
            </CardHeader>
            <CardContent className='grid gap-6'>
                <div className='grid gap-2'>
                    <div className='flex items-center justify-between'>
                        <Label>Customer ID</Label>
                        <span>{outputJson.CustomerId}</span>
                    </div>
                </div>
                <div className='grid gap-4'>
                    <div className='flex items-center justify-between'>
                        <Label>Графики доставки</Label>
                    </div>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Дата начала</TableHead>
                                <TableHead>Дата окончания</TableHead>
                                <TableHead>Год доставки</TableHead>
                                <TableHead>Количество</TableHead>
                                <TableHead>Товар</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {deliverySchedules?.map((schedule, scheduleIndex) => (
                                <TableRow key={scheduleIndex}>
                                    <TableCell>
                                        {editMode ? (
                                            <Input
                                                type='date'
                                                value={schedule.DeliverySchedule.start_date}
                                                onChange={(e) => {
                                                    const updatedSchedules = [...deliverySchedules];
                                                    updatedSchedules[
                                                        scheduleIndex
                                                    ].DeliverySchedule.start_date = e.target.value;
                                                    setDeliverySchedules(updatedSchedules);
                                                }}
                                            />
                                        ) : (
                                            <span>{schedule.DeliverySchedule.start_date}</span>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {editMode ? (
                                            <Input
                                                type='date'
                                                value={schedule.DeliverySchedule.end_date}
                                                onChange={(e) => {
                                                    const updatedSchedules = [...deliverySchedules];
                                                    updatedSchedules[
                                                        scheduleIndex
                                                    ].DeliverySchedule.end_date = e.target.value;
                                                    setDeliverySchedules(updatedSchedules);
                                                }}
                                            />
                                        ) : (
                                            <span>{schedule.DeliverySchedule.end_date}</span>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {editMode ? (
                                            <Input
                                                type='number'
                                                value={schedule.DeliverySchedule.year}
                                                onChange={(e) => {
                                                    const updatedSchedules = [...deliverySchedules];
                                                    updatedSchedules[
                                                        scheduleIndex
                                                    ].DeliverySchedule.year = parseInt(
                                                        e.target.value
                                                    );
                                                    setDeliverySchedules(updatedSchedules);
                                                }}
                                            />
                                        ) : (
                                            <span>{schedule.DeliverySchedule.year}</span>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {editMode ? (
                                            <Input
                                                type='number'
                                                value={
                                                    schedule.DeliverySchedule.deliveryAmount || 0
                                                }
                                                onChange={(e) => {
                                                    const updatedSchedules = [...deliverySchedules];
                                                    updatedSchedules[
                                                        scheduleIndex
                                                    ].DeliverySchedule.deliveryAmount = parseFloat(
                                                        e.target.value
                                                    );
                                                    setDeliverySchedules(updatedSchedules);
                                                }}
                                            />
                                        ) : (
                                            <span>{schedule.DeliverySchedule.deliveryAmount}</span>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {editMode ? (
                                            <div className='grid gap-2'>
                                                <Input
                                                    type='number'
                                                    value={schedule.purchaseAmount || 0}
                                                    onChange={(e) => {
                                                        const updatedSchedules = [
                                                            ...deliverySchedules,
                                                        ];
                                                        updatedSchedules[
                                                            scheduleIndex
                                                        ].purchaseAmount = parseFloat(
                                                            e.target.value
                                                        );
                                                        setDeliverySchedules(updatedSchedules);
                                                    }}
                                                />
                                                <span>
                                                    СПГЗ ID: {schedule.spgzCharacteristics.spgzId}
                                                </span>
                                                <span>
                                                    Название СПГЗ:{' '}
                                                    {schedule.spgzCharacteristics.spgzName}
                                                </span>
                                                <span>
                                                    Код КПГЗ:{' '}
                                                    {schedule.spgzCharacteristics.kpgzCode}
                                                </span>
                                                <span>
                                                    Название KPGZ:{' '}
                                                    {schedule.spgzCharacteristics.kpgzName}
                                                </span>
                                            </div>
                                        ) : (
                                            <div className='grid gap-2'>
                                                <span>
                                                    {(schedule.purchaseAmount || 0).toFixed(2)}
                                                </span>
                                                <span>
                                                    СПГЗ ID: {schedule.spgzCharacteristics.spgzId}
                                                </span>
                                                <span>
                                                    Название СПГЗ:{' '}
                                                    {schedule.spgzCharacteristics.spgzName}
                                                </span>
                                                <span>
                                                    Код КПГЗ:{' '}
                                                    {schedule.spgzCharacteristics.kpgzCode}
                                                </span>
                                                <span>
                                                    Название КПГЗ:{' '}
                                                    {schedule.spgzCharacteristics.kpgzName}
                                                </span>
                                            </div>
                                        )}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>

                <Button variant='outline' onClick={downloadFile}>
                    Загрузить план закупок (.json)
                </Button>
            </CardContent>
        </Card>
    );
};

export default SavedOutputJson;

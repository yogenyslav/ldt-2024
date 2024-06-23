import { Stock } from '@/api/models';
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from './ui/table';
import { mapQuarterToLabel } from '@/utils/map-quarter';
import { CSVLink } from 'react-csv';
import { jsonToCSV } from 'react-papaparse';
import { Button } from './ui/button';
import BasicBarChart from './BasicBarChart';
import { hslStringToHex } from '@/utils/hsl-to-hex';

type Props = {
    stocks: Stock[];
};

const Stocks = ({ stocks }: Props) => {
    const csvData = jsonToCSV(stocks);

    const rootStyles = getComputedStyle(document.documentElement);

    const dataAmount = stocks.map((stock) => ({
        x: `${stock.year} ${mapQuarterToLabel(stock.quarter)}`,
        y: stock.amount,
        fill: hslStringToHex(rootStyles.getPropertyValue('--primary').trim()),
    }));

    const dataSum = stocks.map((stock) => ({
        x: `${stock.year} ${mapQuarterToLabel(stock.quarter)}`,
        y: stock.sum,
        fill: '#e11d48',
    }));

    return (
        <div className='w-full'>
            <div className='rounded-lg border p-2 mt-4'>
                <div>
                    <BasicBarChart
                        data={dataAmount}
                        title='Остатки товара (количество)'
                        xLabel='Кварталы'
                        tooltipItemName='Количество товара'
                        tooltipPostfix='шт.'
                    />
                </div>
            </div>

            <div className='rounded-lg border p-2 mt-4'>
                <div>
                    <BasicBarChart
                        data={dataSum}
                        title='Остатки товара (сумма руб.)'
                        xLabel='Кварталы'
                        tooltipItemName='Сумма'
                        tooltipPostfix='₽'
                    />
                </div>
            </div>

            <Table className='w-full mt-4'>
                <TableCaption>Остатки товара</TableCaption>
                <TableHeader>
                    <TableRow>
                        <TableHead>Наименование</TableHead>
                        <TableHead>Стоимость</TableHead>
                        <TableHead>Количество</TableHead>
                        <TableHead className='text-right'>Год</TableHead>
                        <TableHead className='text-right'>Квартал</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {stocks.map((stock) => (
                        <TableRow key={stock.quarter}>
                            <TableCell className='font-medium'>{stock.name}</TableCell>
                            <TableCell className='text-right'>{stock.price}</TableCell>
                            <TableCell className='text-right'>{stock.amount}</TableCell>
                            <TableCell className='text-right'>{stock.year}</TableCell>
                            <TableCell className='text-right'>
                                {mapQuarterToLabel(stock.quarter)}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
            <CSVLink data={csvData} filename={'data.csv'}>
                <Button variant={'outline'}>Скачать CSV</Button>
            </CSVLink>
        </div>
    );
};

export default Stocks;

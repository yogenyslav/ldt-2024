import { Stock, StockResponse } from '@/api/models';
import Stocks from './Stocks';
import {
    SelectItem,
    Select,
    SelectContent,
    SelectGroup,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from './ui/select';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { useState } from 'react';

type Props = {
    stocksGroup: StockResponse['data'];
};

const StocksGroup = ({ stocksGroup }: Props) => {
    const [selectedStock, setSelectedStock] = useState<string | null>(null);

    return (
        <>
            <Card className='mt-2'>
                <CardHeader>
                    <CardTitle>Нашлось несколько товаров</CardTitle>
                </CardHeader>
                <CardContent>
                    <p className='text-sm'>Выберите товар, чтобы увидеть его остатки на складе</p>

                    <div className='grid gap-6 mt-2'>
                        <div className='grid gap-3'>
                            <Select
                                onValueChange={(stock) => {
                                    setSelectedStock(stock);
                                }}
                            >
                                <SelectTrigger id='status' aria-label='Select status'>
                                    <SelectValue placeholder='Выберите товар' />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectGroup>
                                        <SelectLabel>Товары</SelectLabel>
                                        {stocksGroup.map((stocks, index) => (
                                            <SelectItem key={index} value={stocks.name}>
                                                {stocks.name}
                                            </SelectItem>
                                        ))}
                                    </SelectGroup>
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <div className='gap-4'>
                {selectedStock && stocksGroup.find((stocks) => stocks.name === selectedStock) && (
                    <Stocks
                        stocks={
                            stocksGroup
                                .find((stocks) => stocks.name === selectedStock)
                                ?.history.sort((a, b) => a.quarter - b.quarter) as Stock[]
                        }
                    />
                )}
            </div>
        </>
    );
};

export default StocksGroup;

import { PurchasePlan } from '@/api/models';
import { hslStringToHex } from '@/utils/hsl-to-hex';
import { useCallback, useEffect, useState } from 'react';
import {
    BarChart,
    Bar,
    Rectangle,
    XAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from 'recharts';
import { useCurrentPng } from 'recharts-to-png';
import FileSaver from 'file-saver';
import { Button } from './ui/button';

type Props = {
    history: PurchasePlan[];
    forecast: PurchasePlan[];
};

const Prediction = ({ history, forecast }: Props) => {
    const [getPng, { ref }] = useCurrentPng();

    const handleDownload = useCallback(async () => {
        const png = await getPng();

        if (png) {
            FileSaver.saveAs(png, 'Прогноз закупок.png');
        }
    }, [getPng]);

    const [screenWidth, setScreenWidth] = useState(window.innerWidth);

    useEffect(() => {
        const handleResize = () => setScreenWidth(window.innerWidth);
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    const rootStyles = getComputedStyle(document.documentElement);

    const mapData = (data: PurchasePlan[], isHistory: boolean) => {
        return data.map(({ date, value, volume }) => ({
            date,
            value,
            volume,
            fill: isHistory
                ? hslStringToHex(rootStyles.getPropertyValue('--primary').trim())
                : '#e11d48',
        }));
    };

    const data = mapData(history, true).concat(mapData(forecast, false));

    const getInterval = () => {
        if (screenWidth < 600) {
            return Math.ceil(data.length / 4); // Show fewer labels on small screens
        } else if (screenWidth < 1024) {
            return Math.ceil(data.length / 8); // Medium screens
        } else {
            return 0; // Show all labels on large screens
        }
    };

    return (
        <>
            <div>
                <div className='flex items-center justify-between flex-wrap'>
                    <h2 className='text-xl font-bold'>Прогноз закупок</h2>

                    <Button onClick={handleDownload} variant='outline'>
                        Скачать диаграмму
                    </Button>
                </div>

                <div className='mt-2 rounded-lg border w-full h-[300px]'>
                    <ResponsiveContainer width='100%' height={300}>
                        <BarChart
                            ref={ref}
                            height={300}
                            data={data}
                            margin={{
                                top: 5,
                                right: 30,
                                left: 20,
                                bottom: 5,
                            }}
                        >
                            <CartesianGrid strokeDasharray='3 3' />
                            <XAxis
                                dataKey='date'
                                interval={getInterval()}
                                axisLine={false}
                                angle={45}
                                dx={window.innerWidth < 600 ? 0 : 20}
                                dy={20}
                            />
                            {/* <YAxis /> */}
                            <Tooltip />
                            <Legend />
                            <Bar
                                dataKey='value'
                                fill={hslStringToHex(
                                    rootStyles.getPropertyValue('--primary').trim()
                                )}
                                activeBar={<Rectangle />}
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>
        </>
    );
};

export default Prediction;

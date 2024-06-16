import { hslStringToHex } from '@/utils/hsl-to-hex';
import { useCallback, useEffect, useState } from 'react';
import {
    BarChart,
    Bar,
    Rectangle,
    XAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    TooltipProps,
} from 'recharts';
import { useCurrentPng } from 'recharts-to-png';
import FileSaver from 'file-saver';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { NameType, ValueType } from 'recharts/types/component/DefaultTooltipContent';

type Props = {
    data: {
        x: string | number;
        y: number;
        fill: string;
    }[];
    title: string;
    xLabel: string;
    tooltipItemName?: string;
    tooltipPostfix?: string;
};

const BasicBarChart = ({ data, title, tooltipItemName, tooltipPostfix }: Props) => {
    const [getPng, { ref }] = useCurrentPng();

    const handleDownload = useCallback(async () => {
        const png = await getPng();

        if (png) {
            FileSaver.saveAs(png, `${title}.png`);
        }
    }, [getPng, title]);

    const [screenWidth, setScreenWidth] = useState(window.innerWidth);

    useEffect(() => {
        const handleResize = () => setScreenWidth(window.innerWidth);
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    const rootStyles = getComputedStyle(document.documentElement);

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
                    <h2 className='text-xl font-bold'>{title}</h2>

                    <Button onClick={handleDownload} variant='outline'>
                        Скачать диаграмму
                    </Button>
                </div>

                <div className='mt-2 rounded-lg border w-full h-[320px]'>
                    <ResponsiveContainer width='100%' height={320}>
                        <BarChart
                            ref={ref}
                            height={320}
                            data={data}
                            margin={{
                                top: 5,
                                right: 30,
                                left: 20,
                                bottom: 45,
                            }}
                        >
                            <CartesianGrid strokeDasharray='3 3' />
                            <XAxis
                                dataKey='x'
                                interval={getInterval()}
                                axisLine={false}
                                angle={45}
                                dx={window.innerWidth < 600 ? 0 : 20}
                                dy={20}
                            ></XAxis>

                            <Tooltip
                                content={
                                    <CustomTooltip
                                        itemName={tooltipItemName}
                                        postfix={tooltipPostfix}
                                    />
                                }
                            />

                            <Bar
                                dataKey='y'
                                fill={hslStringToHex(
                                    rootStyles.getPropertyValue('--primary').trim()
                                )}
                                activeBar={<Rectangle />}
                                radius={[10, 10, 0, 0]}
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>
        </>
    );
};

export default BasicBarChart;

interface CustomTooltipProps extends TooltipProps<ValueType, NameType> {
    itemName?: string;
    postfix?: string;
}

const CustomTooltip = ({ active, payload, label, postfix, itemName }: CustomTooltipProps) => {
    if (active && payload && payload.length) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className='text-sm'>{label}</CardTitle>
                </CardHeader>
                <CardContent>
                    {payload.map((item, index) => (
                        <div key={index}>
                            <span>{itemName ? itemName : item.name}</span>
                            <span>
                                : {item.value} {postfix}
                            </span>
                        </div>
                    ))}
                </CardContent>
            </Card>
        );
    }
};

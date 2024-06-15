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
    Label,
} from 'recharts';
import { useCurrentPng } from 'recharts-to-png';
import FileSaver from 'file-saver';
import { Button } from './ui/button';

type Props = {
    data: {
        x: string | number;
        y: number;
        fill: string;
    }[];
    title: string;
    xLabel: string;
};

const BasicBarChart = ({ data, title, xLabel }: Props) => {
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
                                dataKey='x'
                                interval={getInterval()}
                                axisLine={false}
                                angle={45}
                                dx={window.innerWidth < 600 ? 0 : 20}
                                dy={20}
                            >
                                <Label value={xLabel} offset={0} position='insideBottom' />
                            </XAxis>
                            {/* <YAxis /> */}
                            <Tooltip />
                            <Legend />
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

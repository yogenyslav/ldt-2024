import { PurchasePlan } from '@/api/models';
import { hslStringToHex } from '@/utils/hsl-to-hex';
import BasicBarChart from './BasicBarChart';

type Props = {
    history: PurchasePlan[];
    forecast: PurchasePlan[];
};

const Prediction = ({ history, forecast }: Props) => {
    const rootStyles = getComputedStyle(document.documentElement);

    const mapData = (data: PurchasePlan[], isHistory: boolean) => {
        return data.map(({ date, value, volume }) => ({
            date,
            value: Math.round(value),
            volume: volume ? Math.round(volume) : null,
            fill: isHistory
                ? hslStringToHex(rootStyles.getPropertyValue('--primary').trim())
                : '#e11d48',
        }));
    };

    const data = mapData(history, true).concat(mapData(forecast, false));

    return (
        <>
            <div className='rounded-lg border p-2'>
                {' '}
                <BasicBarChart
                    data={data.map(({ date, value, fill }) => ({
                        x: date,
                        y: value,
                        fill,
                    }))}
                    title='Прогноз закупок'
                    xLabel='Сумма закупки'
                />
            </div>
        </>
    );
};

export default Prediction;

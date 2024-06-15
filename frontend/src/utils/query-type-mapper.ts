import { IncomingMessageType } from '@/api/models';

export const mapIncomingMessageTypeToText = (type: IncomingMessageType) => {
    switch (type) {
        case IncomingMessageType.Prediction:
            return 'Прогноз';
        case IncomingMessageType.Stock:
            return 'Запасы';
        case IncomingMessageType.Undefined:
            return 'Не определено';
        default:
            return 'Не определено';
    }
};

export const mapQuarterToLabel = (quarter: number) => {
    switch (quarter) {
        case 1:
            return 'Q1';
        case 2:
            return 'Q2';
        case 3:
            return 'Q3';
        case 4:
            return 'Q4';
        default:
            return '';
    }
};

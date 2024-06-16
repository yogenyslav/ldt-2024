export const hslStringToHex = (hslString: string): string => {
    // Split the string by spaces
    const [h, s, l] = hslString.split(/\s+/).map((value, index) => {
        // Remove the '%' sign from saturation and lightness and convert to number
        if (index === 0) return parseFloat(value); // Hue remains as is
        return parseFloat(value.replace('%', '')); // Remove '%' from saturation and lightness
    });

    // Convert to hex using the hslToHex function
    return hslToHex(h, s, l);
};

// Helper function to convert HSL to Hex
const hslToHex = (h: number, s: number, l: number): string => {
    l /= 100;
    const a = (s * Math.min(l, 1 - l)) / 100;
    const f = (n: number) => {
        const k = (n + h / 30) % 12;
        const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
        return Math.round(255 * color)
            .toString(16)
            .padStart(2, '0'); // Convert to hex and pad with zeros
    };
    return `#${f(0)}${f(8)}${f(4)}`;
};

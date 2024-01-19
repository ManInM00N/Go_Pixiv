
export function padStart(value, length, padChar) {
    value = String(value);
    while (value.length < length) {
        value = padChar + value;
    }
    return value;
}
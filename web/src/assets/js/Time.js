import { ref } from "vue";
import { Events } from "@wailsio/runtime";
export function padStart(value, length, padChar) {
    value = String(value);
    while (value.length < length) {
        value = padChar + value;
    }
    return value;
}

export let timeElement = ref("");
Events.On('time', (time) => {
    timeElement.value = time.data[0];
});
export const sleep = ms => new Promise(r => setTimeout(r, ms));

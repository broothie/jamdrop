import * as api from "./api";
import {notify} from "./notify";

export const startPing = () => {
    api.ping();

    setInterval(() => {
        api.ping().then((data) => {
            if (data.song_queued_events !== null) {
                data.song_queued_events.forEach((event) => {
                    notify(`${event.user_name} dropped "${event.song_name}" to your queue`);
                });
            }
        });
    }, 10 * 1000);
};


export const notify = (message) => {
    if (!("Notification" in window)) {
        console.log('Notifications not available');
    } else if (Notification.permission === "granted") {
        new Notification(message);
    } else if (Notification.permission !== "denied") {
        Notification.requestPermission((permission) => {
            if (permission === "granted") {
                new Notification(message);
            }
        });
    }
};

import m from "mithril";

export default {
    message: null,
    setEl(el, time = 5) {
        this.message = el;
        m.redraw();

        setTimeout(() => {
            this.message = null;
            m.redraw();
        }, time * 1000);
    },
    setMessage(message, time = 5) {
        this.setEl(m('p.toast.message', message), time);
    },
    setError(error, time = 5) {
        this.setEl(m('p.toast.error', error), time);
    }
};

import m from 'mithril';
import {Sharers} from "./sharers";
import {Shares} from "./shares";
import jam from '../assets/jam.svg';

export const Main = (vnode) => {
    const { userData } = vnode.attrs;

    const messenger = {
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

    return {
        view: () => m('main',
            m('.welcome',
                m('.logo', m('img', {src: jam}), m('p', 'JamDrop')),
                messenger.message,
                m('div', m('p', `Welcome, ${userData.user.name} 👋`))
            ),
            m(Sharers, { sharers: userData.sharers, messenger }),
            m(Shares, { shares: userData.shares, messenger }),
        )
    };
};

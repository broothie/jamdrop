import m from 'mithril';
import {Sharers} from "./sharers";
import {Shares} from "./shares";
import jam from '../assets/jam.svg';

export const Main = () => ({
    view(vnode) {
        const { userData } = vnode.attrs;

        return m('main',
            m('.welcome',
                m('.logo', m('img', { src: jam }), m('p', 'JamDrop')),
                m('div', m('p', `Welcome, ${userData.user.name} 👋`))
            ),
            userData.sharers.length > 0 && m(Sharers, { sharers: userData.sharers }),
            m(Shares, { shares: userData.shares }),
        );
    }
});

import m from 'mithril';
import {Add} from "./add";
import {Sharer} from "./sharer";
import {Share} from "./share";

export const Main = () => ({
    view(vnode) {
        const { userData } = vnode.attrs;

        return m('main', [
            m('.welcome', m('p', `👋 ${userData.user.name}`)),
            m('p.sharers-title', 'queues you can drop to'),
            m('.sharers', [...userData.sharers.map((sharer) => m(Sharer, { sharer })), m(Add)]),
            m('p.shares-title', "users you've shared your queue with"),
            m('.shares', userData.shares.map((share) => m(Share, { share })))
        ]);
    }
});

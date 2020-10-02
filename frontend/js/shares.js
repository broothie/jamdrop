import m  from 'mithril'
import {AddShare} from "./add_share";
import * as api from "./api";

export const Shares = (vnode) => {
    let { shares } = vnode.attrs;
    let message = null;

    const reload = () => {
        api.getShares().then((data) => shares = data);
    };

    return {
        view: () => m('.shares-container',
            shares.length > 0 && m('.shares-header',
                m('p.shares-title', "↓ users you've shared your queue with"),
                m('p.shares-message', message),
            ),
            m('.shares', m(AddShare, { reload }), ...shares.map((share) => m(Share, { share })))
        )
    };
};

export const Share = () => ({
    view(vnode) {
        const { share } = vnode.attrs;

        return m('.share',
            m('img.image', { src: share.image_url }),
            m('.name', m('p', share.name))
        );
    }
});

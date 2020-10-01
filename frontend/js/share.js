import m  from 'mithril'

export const Share = () => ({
    view(vnode) {
        const { share } = vnode.attrs;

        return m('.share', [
            m('.name', m('p', share.name))
        ]);
    }
});

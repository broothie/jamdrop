import m from 'mithril';

export const Sharer = () => ({
    view(vnode) {
        const { sharer } = vnode.attrs;

        return m('.sharer.card', [
            m('img.image', {src: sharer.image_url}),
            m('.name', m('p', sharer.name)),
        ]);
    }
});

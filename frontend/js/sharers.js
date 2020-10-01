import m from 'mithril';

export const Sharers = (vnode) => {
    let { sharers } = vnode.attrs;
    let message = null;

    const setMessage = (msg) => {
        message = msg;
        setTimeout(() => message = null, 3000);
    };

    return {
        view: () => m('div',
            m('.sharers-header',
                m('p.sharers-title', '↓ queues you can drop to'),
                m('p.sharers-message', message),
            ),
            m('.sharers', sharers.map((sharer) => m(Sharer, { sharer, setMessage }))),
        )
    };
};

export const Sharer = () => ({
    view(vnode) {
        const { sharer, setMessage } = vnode.attrs;

        const ondrop = (event) => {
            event.preventDefault();
            const song_identifier = event.dataTransfer.getData('text/plain');

            m.request({ method: 'post', url: '/api/users/:id/queue', params: { id: sharer.id, song_identifier } })
                .then((res) => setMessage(res.message));
        };

        const ondragover = (event) => {
            event.preventDefault();
            event.dataTransfer.dropEffect = 'link';
        };

        return m('.sharer.card', { ondrop, ondragover },
            m('img.image', {src: sharer.image_url}),
            m('.name', m('p', sharer.name)),
        );
    }
});

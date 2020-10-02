import m from 'mithril';
import * as api from "./api";

export const Sharers = (vnode) => {
    let { sharers } = vnode.attrs;
    let message = null;

    const setMessage = (msg) => message = msg;

    return {
        view: () => m('.sharers-container',
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
            const songIdentifier = event.dataTransfer.getData('text/plain');

            api.queueSong(sharer.id, songIdentifier).then((res) => {
                setMessage(res.message);

                setTimeout(() => setMessage(null), 3000);
            });
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

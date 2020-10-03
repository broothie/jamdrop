import m from 'mithril';
import * as api from "./api";

export const Sharers = (vnode) => {
    let { sharers, messenger } = vnode.attrs;

    return {
        view: () => {
            sharers = sharers.filter((sharer) => sharer.enabled && sharer.is_playing && sharer.is_active);

            return m('.sharers-container',
                m('.sharers-header', m('p.sharers-title', '↓ queues you can drop to')),
                m('.sharers', sharers.map((sharer) => m(Sharer, {key: sharer.id, sharer, messenger}))),
            );
        }
    };
};

export const Sharer = (vnode) => {
    const { sharer, messenger } = vnode.attrs;

    const ondragstart = (event) => {
        event.dataTransfer.setData('text/plain', sharer.id);
    };

    const ondrop = (event) => {
        event.preventDefault();
        const songIdentifier = event.dataTransfer.getData('text/plain');

        api.queueSong(sharer.id, songIdentifier)
            .then((res) => messenger.setMessage(res.message))
            .catch((e) => messenger.setError(e.response.error));
    };

    const ondragover = (event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'link';
    };

    return {
        view: () => m('.sharer.card', { draggable: true, ondragstart, ondrop, ondragover },
            m('img.image', { src: sharer.image_url, draggable: false }),
            m('.name', m('p', sharer.name)),
        )
    };
};

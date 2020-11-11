import m from 'mithril';
import * as api from "./api";
import toaster from "./toaster";

export const Sharers = (vnode) => {
    let { sharers } = vnode.attrs;

    setInterval(() => {
        api.getSharers()
            .then((data) => sharers = data)
            .then(m.redraw);
    }, 10 * 1000);

    return {
        view: () => {
            sharers = sharers.sort((a, b) => a.id < b.id);
            const enabledSharers = sharers.filter((sharer) => sharer.droppable);
            const disabledSharers = sharers.filter((sharer) => !sharer.droppable);

            return sharers.length > 0 && m('.sharers-container',
                m('.sharers-header', m('p.title', '↓ drop a jam')),
                m('.sharers', [
                    ...enabledSharers.map((sharer) => m(Sharer, {key: sharer.id, sharer})),
                    ...disabledSharers.map((sharer) => m(Sharer, {key: sharer.id, sharer}))
                ]),
            );
        }
    };
};

export const Sharer = (vnode) => {
    const { sharer } = vnode.attrs;
    const ondragstart = (event) => {
        event.dataTransfer.setData('text/plain', sharer.id);
    };

    const queueSong = (songIdentifier) => {
        api.queueSong(sharer.id, songIdentifier)
            .then((res) => toaster.setMessage(res.message))
            .catch((e) => toaster.setError(e.response.error));
    };

    const ondrop = (event) => {
        event.preventDefault();

        const songIdentifier = event.dataTransfer.getData('text/plain');

        queueSong(songIdentifier);
    };

    const onclick = () => {
        const songIdentifier = window.prompt(`Paste a Spotify song ID or link here to drop it to ${sharer.name}'s queue`);
        if (songIdentifier) queueSong(songIdentifier);
    };

    return {
        view: (vnode) => {
            const { sharer } = vnode.attrs;
            const light = sharer.droppable ? m('span.light.active', '●') : m('span.light', '○');

            const ondragover = (event) => {
                event.preventDefault();
                if (!sharer.droppable) return;

                event.dataTransfer.dropEffect = 'link';
            };

            return m(
                '.sharer.card',
                {
                    class: sharer.droppable ? '' : 'disabled',
                    draggable: true,
                    ondragstart,
                    ondrop,
                    ondragover,
                    onclick
                },
                [
                    m('img.image', {src: sharer.image_url, draggable: false}),
                    m('.name', light, m('p', sharer.name)),
                ]
            );
        }
    };
};

import m from 'mithril';
import * as api from "./api";

export const AddShare = (vnode) => {
    const { reload, messenger } = vnode.attrs;

    const ondrop = (event) => {
        event.preventDefault();
        const userIdentifier = event.dataTransfer.getData('text/plain');
        console.log(userIdentifier);

        api.addShare(userIdentifier)
            .then(reload)
            .catch((e) => messenger.setError(e.response.error));
    };

    const ondragover = (event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'link';
    };

    return {
        view() {
            return m('.add.share', { ondrop, ondragover }, m('p', 'drag and drop a Spotify user here to add them'));
        }
    };
};

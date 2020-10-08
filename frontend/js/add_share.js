import m from 'mithril';
import * as api from "./api";
import toaster from "./toaster";

export const AddShare = (vnode) => {
    const { reload } = vnode.attrs;

    const ondrop = (event) => {
        event.preventDefault();
        const userIdentifier = event.dataTransfer.getData('text/plain');

        api.addShare(userIdentifier)
            .then(reload)
            .catch((e) => toaster.setError(e.response.error));
    };

    const ondragover = (event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'link';
    };

    return {
        view() {
            return m('.add.share', { ondrop, ondragover }, m('p', 'drag a Spotify user here to share your queue'));
        }
    };
};

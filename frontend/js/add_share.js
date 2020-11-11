import m from 'mithril';
import * as api from "./api";
import toaster from "./toaster";

export const AddShare = (vnode) => {
    const { reload } = vnode.attrs;

    const addShare = (userIdentifier) => {
        api.addShare(userIdentifier)
            .then(reload)
            .catch((e) => toaster.setError(e.response.error));
    };

    const ondrop = (event) => {
        event.preventDefault();
        const userIdentifier = event.dataTransfer.getData('text/plain');

        addShare(userIdentifier);
    };

    const onclick = () => {
        const userIdentifier = window.prompt('Paste a Spotify user ID or link here to share your queue');
        if (userIdentifier) addShare(userIdentifier);
    };

    const ondragover = (event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'link';
    };

    return {
        view() {
            return m('.add.share', { ondrop, ondragover, onclick }, m('p', 'drag a Spotify user here to share your queue'));
        }
    };
};

import m from 'mithril';

export const AddShare = (vnode) => {
    const { reload } = vnode.attrs;
    const startMessage = 'drag and drop a Spotify user here to add them';
    let text = startMessage;

    const handleError = (message) => {
        text = message;

        setTimeout(() => text = startMessage, 3 * 1000)
    };

    const ondrop = (event) => {
        event.preventDefault();
        const user_identifier = event.dataTransfer.getData('text/plain');

        m.request({ method: 'post', url: '/api/share', params: { user_identifier } })
            .then(reload)
            .catch((error) => handleError(error));
    };

    const ondragover = (event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'link';
    };

    return {
        view() {
            return m('.add.share', { ondrop, ondragover }, m('p', text));
        }
    };
};

import m from 'mithril';

export const Add = () => {
    const ondrop = (event) => {
        event.preventDefault();
        const user_identifier = event.dataTransfer.getData('text/plain');

        m.request({ method: 'post', url: '/api/share', params: { user_identifier } })
    };

    const ondragover = (event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'link';
    };

    return {
        view() {
            return m('.add.card', { ondrop, ondragover }, m('p', 'drag and drop a Spotify user here to add them'));
        }
    };
};

import m from 'mithril';
import jam from "../assets/jam.svg";

export const AuthorizeSpotify = () => {
    const onclick = () => location.href = '/spotify/authorize';

    return {
        view() {
            return m('main',
                m('.welcome', m('.logo', m('img', { src: jam }), m('p', 'JamDrop'))),
                m('.spotify-authorize', m('button', { onclick }, 'Log In with Spotify')),
            );
        }
    };
};

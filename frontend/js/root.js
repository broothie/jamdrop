import m from 'mithril';
import {Main} from "./main";
import {AuthorizeSpotify} from "./authorize_spotify";
import * as api from './api';

export const Root = () => {
    let userData = null;
    let error = null;

    return {
        oninit() {
            api.getMe()
                .then((data) => userData = data)
                .catch((e) => error = e);
        },
        view() {
            if (error !== null) return m(AuthorizeSpotify);

            return userData && m(Main, { userData });
        }
    };
};

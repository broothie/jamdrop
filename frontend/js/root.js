import m from 'mithril';
import {Main} from "./main";

export const Root = () => {
    let userData = null;

    return {
        oninit() {
            m.request('/api/users/me')
                .then((res) => userData = res)
                .catch(() => location.href = '/spotify/authorize');
        },
        view() {
            return userData && m(Main, { userData });
        }
    };
};

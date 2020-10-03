import m from 'mithril';
import {Root} from "./root";
import {startPing} from "./ping";

document.addEventListener('DOMContentLoaded', () => {
    startPing();

    const root = document.getElementById('root');
    m.mount(root, Root);
});

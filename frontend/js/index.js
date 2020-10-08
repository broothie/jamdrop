import m from 'mithril';
import {Root} from "./root";
import {startPing} from "./ping";

document.addEventListener('DOMContentLoaded', () => {
    startPing();

    m.mount(document.body, Root);
});

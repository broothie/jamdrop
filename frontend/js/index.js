import m from 'mithril';
import {Root} from "./root";

document.addEventListener('DOMContentLoaded', () => {
    const root = document.getElementById('root');
    m.mount(root, Root);
});

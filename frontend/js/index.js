import m from 'mithril';
import {Root} from "./root";
import * as api from './api';

document.addEventListener('DOMContentLoaded', () => {
    activeTracker();

    const root = document.getElementById('root');
    m.mount(root, Root);
});

const activeTracker = () => {
    api.ping();
    setInterval(api.ping, 60 * 1000);
};

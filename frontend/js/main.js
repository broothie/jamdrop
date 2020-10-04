import m from 'mithril';
import {Sharers} from "./sharers";
import {Shares} from "./shares";
import jam from '../assets/jam.svg';
import * as api from "./api";

export const Main = (vnode) => {
    const { userData } = vnode.attrs;
    const user = userData.user;

    const messenger = {
        message: null,
        setEl(el, time = 5) {
            this.message = el;
            m.redraw();

            setTimeout(() => {
                this.message = null;
                m.redraw();
            }, time * 1000);
        },
        setMessage(message, time = 5) {
            this.setEl(m('p.toast.message', message), time);
        },
        setError(error, time = 5) {
            this.setEl(m('p.toast.error', error), time);
        }
    };

    return {
        view: () => m('main',
            m('.welcome',
                m('.logo', m('img', {src: jam}), m('p', 'JamDrop')),
                messenger.message,
                m('.user', m('p', `Hi, ${user.name} 👋`), m(Settings, {user}))
            ),
            m(Sharers, {sharers: userData.sharers, messenger}),
            m(Shares, {shares: userData.shares, messenger}),
        )
    };
};

export const Settings = (vnode) => {
    const {user} = vnode.attrs;
    let stayActive = user.stay_active;
    let phoneNumber = user.phone_number;

    return {
        view: () => {
            let stayActiveDisabled = false;
            let phoneNumberDisabled = false;

            const setStayActive = () => {
                stayActiveDisabled = true;
                api.setStayActive(!stayActive)
                    .then(() => stayActive = !stayActive)
                    .then(m.redraw);
            };

            const setPhoneNumber = (event) => {
                const newPhoneNumber = event.target.value;

                phoneNumberDisabled = true;
                api.setPhoneNumber(newPhoneNumber)
                    .then(() => phoneNumber = newPhoneNumber)
                    .then(m.redraw);
            };

            return m('.settings',
                m('.setting',
                    m('input', {
                        type: 'text',
                        placeholder: 'phone number',
                        value: phoneNumber,
                        onblur: setPhoneNumber,
                    })
                ),
                m('.setting',
                    m('input#stay-active', {
                        type: 'checkbox',
                        disabled: stayActiveDisabled,
                        checked: stayActive,
                        onchange: setStayActive
                    }),
                    m('label', {for: 'stay-active'}, 'Stay active')
                ),
            );
        }
    };
};

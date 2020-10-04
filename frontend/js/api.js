import m from 'mithril';

export const getMe = () => m.request('/api/users/me');
export const getShares = () => m.request('/api/users/me/shares');
export const getSharers = () => m.request('/api/users/me/sharers');
export const ping = () => m.request('/api/users/me/ping');

export const addShare = (user_identifier) => m.request({
    method: 'post',
    url: '/api/share',
    params: { user_identifier }
});

export const queueSong = (user_id, song_identifier) => m.request({
    method: 'post',
    url: '/api/users/:user_id/queue',
    params: { user_id, song_identifier }
});

export const setEnabled = (user_id, enabled) => m.request({
    method: 'patch',
    url: '/api/users/:user_id/enabled',
    params: { user_id, enabled }
});

export const setStayActive = (stay_active) => m.request({
    method: 'patch',
    url: '/api/users/me/stay_active',
    params: { stay_active }
});

export const setPhoneNumber = (phone_number) => m.request({
    method: 'patch',
    url: '/api/users/me/phone_number',
    params: { phone_number }
});

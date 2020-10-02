import m from 'mithril';

export const getMe = () => m.request('/api/users/me');
export const getShares = () => m.request('/api/users/me/shares');

export const addShare = (user_identifier) => m.request({
    method: 'post',
    url: '/api/share',
    params: { user_identifier }
});

export const queueSong = (userID, song_identifier) => m.request({
    method: 'post',
    url: '/api/users/:userID/queue',
    params: { userID, song_identifier }
});

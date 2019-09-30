import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {id as pluginId} from './manifest';

// TODO: Move this into mattermost-redux or mattermost-webapp.
export const getPluginServerRoute = (state) => {
    const config = getConfig(state);

    let basePath = '/';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + pluginId;
};

export const checkCanaryCookie = () => async (dispatch, getState) => {
    const canaryCookie = getCookieValue('canary');
    if (canaryCookie) {
        return;
    }

    const resp = await fetch(getPluginServerRoute(getState()) + '/api/v1/check', {
        method: 'GET',
        credentials: 'same-origin',
        headers: {
            'X-Requested-With': 'XMLHttpRequest',
        },
    });

    const data = await resp.json();

    if (data && data.cookieValue === 'always') {
        window.location.reload();
    }
};

// From https://stackoverflow.com/a/25490531
function getCookieValue(a) {
    var b = document.cookie.match('(^|[^;]+)\\s*' + a + '\\s*=\\s*([^;]+)');
    return b ? b.pop() : '';
}
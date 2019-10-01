import {id as pluginId} from './manifest';
import Root from './components/root';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry) {
        registry.registerRootComponent(Root);
    }
}

window.registerPlugin(pluginId, new Plugin());

import {App, Flyout} from './app';



const TITLE = 'AI';
const ID = 'AI';

((window: any) => {
    window?.extensionsAPI?.registerStatusPanelExtension(App, TITLE, ID, Flyout);
})(window);
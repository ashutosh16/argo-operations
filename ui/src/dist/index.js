"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.component = exports.Flyout = exports.Extension = void 0;
const react_1 = require("react");
require('./styles.scss');
const TITLE = "Argo AI";
const ID = "ARGOAI";
const sectionLabel = (info) => (react_1.default.createElement("label", { style: { fontSize: '12px', fontWeight: 600, color: "#6D7F8B" } }, info.title));
const Extension = (props) => {
    var _a, _b, _c;
    return (((_c = (_b = (_a = props.application) === null || _a === void 0 ? void 0 : _a.status) === null || _b === void 0 ? void 0 : _b.health) === null || _c === void 0 ? void 0 : _c.status) != 'Healthy' && react_1.default.createElement("div", { className: "application-status-panel__item" },
        react_1.default.createElement("div", { style: { lineHeight: '19.5px', marginBottom: '0.3em' } }, sectionLabel({ title: 'Intuit GenAI Failure Analysis' })),
        react_1.default.createElement("img", { style: { fill: "red" }, className: "genai-image", src: "assets/images/genai.svg", onClick: props.openFlyout }),
        react_1.default.createElement("div", { className: 'application-status-panel__item-name' },
            react_1.default.createElement("div", null, "Click to start analysising the failed resources"))));
};
exports.Extension = Extension;
const Flyout = (props) => {
    return (react_1.default.createElement(react_1.default.Fragment, null,
        react_1.default.createElement("div", { className: 'application-status-panel__item', style: { position: 'relative' } },
            sectionLabel({
                title: TITLE,
            }),
            react_1.default.createElement("div", { className: 'application-status-panel__item-value', style: { margin: 'auto 0' } },
                react_1.default.createElement("a", { className: 'neutral' },
                    react_1.default.createElement("i", { className: `fa fa-pause-circle` }),
                    " Progressive Sync")))));
};
exports.Flyout = Flyout;
exports.component = exports.Extension;
// Register the component extension in ArgoCD
((window) => {
    var _a;
    (_a = window === null || window === void 0 ? void 0 : window.extensionsAPI) === null || _a === void 0 ? void 0 : _a.registerStatusPanelExtension(exports.component, TITLE, ID, exports.Flyout);
})(window);
//# sourceMappingURL=index.js.map
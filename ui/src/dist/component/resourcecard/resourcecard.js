"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResourceCard = void 0;
const React = require("react");
require("./resourcecard.scss");
const ResourceCard = (props) => {
    const { resNode, application, tree } = Object.assign({}, props);
    return (React.createElement("div", { style: { width: '100%', height: '100%' } }, (React.createElement(DataLoader, { load: () => __awaiter(void 0, void 0, void 0, function* () {
            const res = yield services.applications.managedResources(application.metadata.name, application.metadata.namespace, {
                id: {
                    name: resNode.name,
                    namespace: resNode.namespace,
                    kind: resNode.kind,
                    group: resNode.group
                }
            });
            const node = managedResources.find(item => AppUtils.isSameNode(resNode, item));
            const summary = application.status.resources.find(item => AppUtils.isSameNode(resNode, item));
            const nodeState = (node && summary && { summary, state: node }) || null;
            const resQuery = Object.assign({}, resNode);
            if (node && node.targetState) {
                resQuery.version = AppUtils.parseApiVersion(node.targetState.apiVersion).version;
            }
            const liveState = yield services.applications.getResource(application.metadata.name, application.metadata.namespace, resQuery).catch(() => null);
            const events = (liveState &&
                (yield services.applications.resourceEvents(application.metadata.name, application.metadata.namespace, {
                    name: liveState.metadata.name,
                    namespace: liveState.metadata.namespace,
                    uid: liveState.metadata.uid
                }))) ||
                [];
            return { nodeState, liveState, events };
        }) }, data => (React.createElement(React.Fragment, null,
        React.createElement("div", { key: resNode.uid, className: `argo-table-list__row applications-list__entry applications-list__entry--health-${app.status.health.status}}` },
            React.createElement("div", { className: 'row resourcecard-tiles__wrapper' },
                React.createElement("div", { className: `columns small-12 resourcecard-tiles__info qe-applications-list-${AppUtils.appInstanceName(app)} resourcecard-tiles__item` },
                    React.createElement("div", { className: 'row' },
                        React.createElement("div", { className: 'columns small-3', title: 'Project:' }, "Project:"),
                        React.createElement("div", { className: 'columns small-9' }, app.spec.project)))))))))));
};
exports.ResourceCard = ResourceCard;
//# sourceMappingURL=resourcecard.js.map
import {NodeServer} from '@bitdev/node.node-server';

export default NodeServer.from({
    name: 'dashboard-service',
    mainPath: './dashboard-service.app-root.js',
});

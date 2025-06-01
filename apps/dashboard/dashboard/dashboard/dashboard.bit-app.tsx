import {Platform} from '@bitdev/platforms.platform';

const DashboardWeb = import.meta.resolve('@panelium/dashboard.dashboard-web');
const DashboardService = import.meta.resolve('@panelium/dashboard.dashboard-service');

export const Dashboard = Platform.from({
    name: 'dashboard',

    frontends: {
        main: DashboardWeb,
        mainPortRange: [3000, 3100]
    },

    backends: {
        main: DashboardService,
    },
});

export default Dashboard;

import {DashboardService} from './dashboard-service.js';

describe('corporate service', () => {
    it('should say hello', async () => {
        const dashboardService = DashboardService.from();
        const announcements = await dashboardService.listAnnouncements();
        expect(announcements.length).toEqual(2);
    })
});

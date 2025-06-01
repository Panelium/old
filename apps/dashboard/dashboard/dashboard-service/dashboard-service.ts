import {mockAnnouncements} from "@panelium/dashboard.entities.announcement";

/**
 * corporate service
 */
export class DashboardService {
    /**
     * create a new instance of a corporate service.
     */
    static from() {
        return new DashboardService();
    }

    /**
     * say hello.
     */
    async listAnnouncements() {
        return mockAnnouncements();
    }
}    

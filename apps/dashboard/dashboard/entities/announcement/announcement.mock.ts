import {Announcement} from "./announcement.js";

export function mockAnnouncements() {
    return [
        Announcement.from({
            title: 'Announcing Acme Cloud. A new era for cloud compting',
            date: new Date().toLocaleDateString()
        }),
        Announcement.from({
            title: 'Investing $1B into building Data Center in Chicago',
            date: new Date().toLocaleDateString()
        })
    ];
}

export type PlainAnnouncement = {
    /**
     * title of the announcement
     */
    title: string;

    /**
     * date of the announcement
     */
    date: string;
}

export class Announcement {
    constructor(
        /**
         * name of the announcement.
         */
        readonly title: string,
        /**
         * date of the announcement.
         */
        readonly date: Date,
    ) {
    }

    /**
     * create a Announcement object from a
     * plain object.
     */
    static from(plainAnnouncement: PlainAnnouncement) {
        const date = new Date(plainAnnouncement.date);

        return new Announcement(
            plainAnnouncement.title,
            date
        );
    }

    /**
     * serialize a Announcement into
     * a serializable object.
     */
    toObject() {
        return {
            title: this.title,
            date: this.date.toString()
        };
    }
}

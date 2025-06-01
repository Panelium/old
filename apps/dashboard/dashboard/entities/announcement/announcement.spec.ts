import {Announcement} from './announcement.js';

it('should create an announcement instance', () => {
    const title = 'Introducing new Bit Cloud';
    const announcement = Announcement.from({
        title,
        date: new Date().toString()
    });

    expect(announcement.title).toBe(title);
});

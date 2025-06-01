import {MockProvider} from '@acme/acme.testing.mock-provider';
import {Announcements} from './announcements.js';

export const BasicAnnouncements = () => {
    return (
        // create an env to apply providers on all of your compositions
        // learn more: https://bit.dev/docs/react-env/component-previews
        <MockProvider>
            <Announcements/>
        </MockProvider>
    );
}

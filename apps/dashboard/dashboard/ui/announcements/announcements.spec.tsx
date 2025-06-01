import {render} from '@testing-library/react';
import {BasicAnnouncements} from './announcements.composition.js';

it('should render with text from the mock', () => {
    const {getByText} = render(<BasicAnnouncements/>);
    const rendered = getByText('Announcing Acme Cloud. A new era for cloud compting');
    expect(rendered).toBeTruthy();
});

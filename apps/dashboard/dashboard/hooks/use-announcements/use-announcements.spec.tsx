import {MockProvider} from '@acme/acme.testing.mock-provider';
import {renderHook} from '@testing-library/react';
import {useAnnouncements} from './use-announcements.js';

it('should return two announcements', () => {
    const {result} = renderHook(() => useAnnouncements(), {wrapper: MockProvider});
    expect(result.current?.length).toBe(2);
});
